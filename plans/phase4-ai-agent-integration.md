# Phase 4: AI Agent Integration

## Context

Phases 1-3 are complete: proto definitions, Go services, desktop+mobile UI, Sprint service, status workflow transitions, Kanban boards, AI triage pipeline. Phase 4 adds the tools for AI coding agents to work with l8bugs: an MCP server for programmatic access, GitHub webhook-based git integration with auto-transitions, and enhanced AI root cause analysis for bugs with stack traces.

---

## Step 1: Proto Changes

### 1.1 Add fields to `proto/bugs.proto`

**BugsProject** — add git repo association (fields 10-11, between `created_date=9` and `labels=20`):
```protobuf
string repository_url = 10;
string webhook_secret = 11;
```

**Bug** — add commit linking (field 46, after `related_bug_ids=45`):
```protobuf
string linked_commit_sha = 46;
```

**Feature** — add commit linking (field 47, after `related_bug_ids=46`):
```protobuf
string linked_commit_sha = 47;
```

### 1.2 Regenerate bindings

```bash
cd proto && ./make-bindings.sh
```

### 1.3 Verify

```bash
grep "RepositoryUrl\|LinkedCommitSha\|WebhookSecret" go/types/l8bugs/bugs.pb.go
```

---

## Step 2: MCP Server (Stdio Binary)

Standalone Go binary. AI agents (Claude Code, Cursor) launch it and communicate over stdin/stdout JSON-RPC 2.0.

### 2.1 `go/bugs/mcp/protocol.go` (~100 lines)

JSON-RPC 2.0 message types. No external dependencies — raw `encoding/json`:

```go
type Request struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      json.RawMessage `json:"id,omitempty"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      json.RawMessage `json:"id,omitempty"`
    Result  interface{} `json:"result,omitempty"`
    Error   *RPCError   `json:"error,omitempty"`
}

// MCP types: InitializeResult, Capabilities, ServerInfo, ToolDef, CallToolParams, CallToolResult, ContentBlock
```

### 2.2 `go/bugs/mcp/server.go` (~120 lines)

Stdin/stdout loop with method dispatch:

```go
type Server struct {
    vnic     ifs.IVNic
    tools    map[string]ToolHandler
    toolDefs []ToolDef
}

func NewServer(vnic ifs.IVNic) *Server
func (s *Server) Run()  // bufio.Scanner on stdin, json.Marshal to stdout
```

Methods handled: `initialize`, `initialized`, `tools/list`, `tools/call`, `ping`.

All log output goes to **stderr** (stdout reserved for JSON-RPC).

### 2.3 `go/bugs/mcp/tools.go` (~150 lines)

Tool definitions with JSON Schema input schemas:

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_issues` | Query bugs/features with filters | `type`, `status`, `priority`, `assignee_id`, `project_id`, `limit` |
| `read_issue` | Get full issue details + comments | `issue_id` |
| `create_issue` | Create bug or feature | `type`, `project_id`, `title`, `description`, + optional fields |
| `update_issue` | Update fields on existing issue | `issue_id`, + any updatable fields |
| `add_comment` | Add comment to issue | `issue_id`, `body`, `author_type` |
| `search_issues` | Text search across title/description | `query`, `type`, `project_id` |

### 2.4 `go/bugs/mcp/handlers.go` (~250 lines)

Tool handler implementations using `common.GetEntity`, `PostEntity`, `PutEntity`, `GetEntities` via VNic.

- `handleListIssues` — builds filter from params, calls `GetEntities[Bug]` or `GetEntities[Feature]`
- `handleReadIssue` — tries Bug by ID, falls back to Feature. Returns full proto JSON via `protojson`
- `handleCreateIssue` — constructs Bug or Feature from params, calls `PostEntity`
- `handleUpdateIssue` — fetches existing entity, applies field updates, calls `PutEntity`
- `handleAddComment` — fetches entity, appends Comment to `comments` field, calls `PutEntity`
- `handleSearchIssues` — fetches entities, substring match on title+description

Entity serialization to JSON uses `google.golang.org/protobuf/encoding/protojson`.

### 2.5 `go/bugs/mcp/main1/main.go` (~40 lines)

Entry point — connects VNic to VNet, starts MCP server:

```go
func main() {
    fmt.Fprintln(os.Stderr, "[mcp] L8Bugs MCP Server starting...")
    nic := website.CreateVnic(common.BUGS_VNET)
    server := mcp.NewServer(nic)
    server.Run()
}
```

Uses `website.CreateVnic()` — same pattern as the web server in `website/main1/main.go`.

---

## Step 3: Git Integration (Webhook Handler)

Webhook handler registered on the existing web server (port 2883).

### 3.1 `go/bugs/webhook/webhook.go` (~200 lines)

```go
func Register(vnic ifs.IVNic)  // Registers POST /bugs/webhook/github on http.DefaultServeMux
```

Handler flow:
1. Read request body + `X-Hub-Signature-256` and `X-GitHub-Event` headers
2. Find project by `repository_url` from webhook payload → get `webhook_secret`
3. Verify HMAC-SHA256 signature
4. **PR merged event** (`pull_request`, `action=closed`, `merged=true`):
   - Extract issue refs from PR title + body (e.g., "Fixes L8B-42")
   - For each ref: fetch Bug/Feature → set `linked_pr_url`, `linked_commit_sha` (merge commit)
   - Auto-transition: Bug `In Review(4) → Resolved(5)`, Feature `In Review(5) → Done(6)`
5. **Push event** (`push`):
   - Extract issue refs from commit messages
   - Set `linked_commit_sha` on referenced issues

### 3.2 `go/bugs/webhook/parser.go` (~100 lines)

Issue reference extraction:
```go
func ExtractIssueRefs(text string) []string
```

Matches patterns: `Fixes #42`, `Closes L8B-42`, `Resolves <uuid>`, case-insensitive.

### 3.3 Integration

In `go/bugs/website/main1/main.go`, after VNic creation:
```go
webhook.Register(nic1)
```

---

## Step 4: AI Root Cause Analysis Enhancement

### 4.1 `go/bugs/triage/rootcause.go` (~180 lines)

**Stack trace parser** — extracts structured info from Go, Java, Python, JS stack traces:

```go
type StackFrame struct {
    File, Function, Package string
    Line                    int
}

type StackTraceInfo struct {
    Language  string        // go, java, python, javascript, unknown
    Frames    []StackFrame
    ErrorType string        // e.g., "NullPointerException", "panic"
    ErrorMsg  string
}

func ParseStackTrace(trace string) *StackTraceInfo
```

**Enhanced root cause prompt** — when stack trace is present, builds a structured prompt:
```go
func BuildRootCausePrompt(bug *l8bugs.Bug, stackInfo *StackTraceInfo, repoURL string) string
```

Returns JSON with: `root_cause`, `likely_files[]`, `error_category`, `is_regression_likely`, `confidence`, `suggested_fix`.

**Standalone analysis function** (for re-analysis or MCP tool use):
```go
func (t *Triager) AnalyzeRootCause(bug *l8bugs.Bug) (*RootCauseResult, error)
```

### 4.2 Integrate into existing triage flow

In `triage.go` `triageBug()`, after basic triage completes — if `bug.StackTrace != ""`:
- Call `AnalyzeRootCause(bug)`
- Store enhanced result in `bug.AiRootCause` (combines root cause, category, likely files, suggested fix)
- Put updated bug

---

## Step 5: Desktop UI Updates

### 5.1 `l8tracking-forms.js`

**BugsProject** — add to "Project Details" section:
```javascript
...f.text('repositoryUrl', 'Repository URL'),
...f.text('webhookSecret', 'Webhook Secret'),
```

**Bug** — add to "Resolution" section:
```javascript
...f.text('linkedCommitSha', 'Linked Commit'),
```

**Feature** — add to "Links" section:
```javascript
...f.text('linkedCommitSha', 'Linked Commit'),
```

### 5.2 `l8tracking-columns.js`

**BugsProject** — add column:
```javascript
...col.col('repositoryUrl', 'Repository'),
```

---

## Step 6: Mobile UI Updates (Parity)

### 6.1 `m/js/tracking/l8tracking-forms.js`

Mirror desktop: add `repositoryUrl`/`webhookSecret` to BugsProject, `linkedCommitSha` to Bug and Feature.

### 6.2 `m/js/tracking/l8tracking-columns.js`

Mirror desktop: add `repositoryUrl` to BugsProject.

---

## Step 7: PRD Update

Mark Phase 4 as complete with details in `plans/l8bugs-prd.md`.

---

## File Summary

| Action | Count | Files |
|--------|-------|-------|
| Modify (Proto) | 1 | `proto/bugs.proto` |
| Regen (Proto) | 1 | `go/types/l8bugs/bugs.pb.go` |
| Create (Go MCP) | 5 | `mcp/protocol.go`, `mcp/server.go`, `mcp/tools.go`, `mcp/handlers.go`, `mcp/main1/main.go` |
| Create (Go Webhook) | 2 | `webhook/webhook.go`, `webhook/parser.go` |
| Create (Go Triage) | 1 | `triage/rootcause.go` |
| Modify (Go) | 2 | `triage/triage.go`, `website/main1/main.go` |
| Modify (Desktop) | 2 | `l8tracking-forms.js`, `l8tracking-columns.js` |
| Modify (Mobile) | 2 | `m/js/tracking/l8tracking-forms.js`, `m/js/tracking/l8tracking-columns.js` |
| Modify (PRD) | 1 | `plans/l8bugs-prd.md` |
| **Total** | **17** | 8 new + 9 modified |

---

## Verification

```bash
# Proto
grep "RepositoryUrl\|LinkedCommitSha\|WebhookSecret" go/types/l8bugs/bugs.pb.go

# Go build
cd go && go build ./... && go vet ./...

# MCP binary builds
cd go/bugs/mcp/main1 && go build .

# JS syntax
for f in go/bugs/website/web/l8ui/sys/tracking/*.js; do node -c "$f"; done
for f in go/bugs/website/web/m/js/tracking/*.js; do node -c "$f"; done

# MCP server has tools
grep 'list_issues\|read_issue\|create_issue\|update_issue\|add_comment\|search_issues' go/bugs/mcp/tools.go

# Webhook registered
grep 'webhook.Register' go/bugs/website/main1/main.go

# Root cause analysis exists
grep 'AnalyzeRootCause' go/bugs/triage/rootcause.go

# Root cause integrated into triage
grep 'AnalyzeRootCause' go/bugs/triage/triage.go
```
