# Phase 3: AI Triage — Implementation Plan

## Overview

Phase 3 adds AI-powered triage to l8bugs. When a bug or feature is created, an async AI pipeline classifies, prioritizes, detects duplicates, suggests assignees/components, and populates the existing `ai_*` fields. A Triage Inbox UI lets humans review and accept/override AI suggestions.

---

## Architecture Decision: LLM Integration

**Approach**: Direct Anthropic Messages API via HTTP (`net/http`). No SDK dependency — the Messages API is a single POST endpoint.

**Why not an SDK?**: The Anthropic Go SDK adds a dependency tree. The Messages API only needs one HTTP POST with a JSON body. A thin client (~100 lines) is simpler, vendorable, and has no transitive dependencies.

**Configuration**: API key via environment variable `L8BUGS_ANTHROPIC_API_KEY`. Model and triage settings stored in project configuration (future phase — for now, defaults in code).

---

## Step 1: Proto Additions

### 1.1 Add missing AI fields to `proto/bugs.proto`

Bug is missing `ai_suggested_severity` and `ai_suggested_assignee_id`. Feature is missing `ai_suggested_component` and `ai_suggested_assignee_id`. Both need a triage status to track whether AI has processed them.

```protobuf
enum TriageStatus {
    TRIAGE_STATUS_UNSPECIFIED = 0;
    TRIAGE_STATUS_PENDING = 1;
    TRIAGE_STATUS_IN_PROGRESS = 2;
    TRIAGE_STATUS_COMPLETED = 3;
    TRIAGE_STATUS_FAILED = 4;
    TRIAGE_STATUS_SKIPPED = 5;  // AI triage disabled or no API key
}
```

Add to Bug message (after existing ai fields):
```protobuf
Severity ai_suggested_severity = 36;
string ai_suggested_assignee_id = 37;
TriageStatus triage_status = 38;
string triage_error = 39;  // Error message if triage failed
```

Add to Feature message (after existing ai fields):
```protobuf
string ai_suggested_component = 27;
string ai_suggested_assignee_id = 28;
TriageStatus triage_status = 29;
string triage_error = 30;
```

### 1.2 Regenerate protobuf bindings

```bash
cd proto && ./make-bindings.sh
```

### 1.3 Verify generated field names

```bash
grep -A 50 "type Bug struct" go/types/l8bugs/bugs.pb.go | grep 'json:"'
grep -A 40 "type Feature struct" go/types/l8bugs/bugs.pb.go | grep 'json:"'
```

---

## Step 2: Anthropic LLM Client

### 2.1 Create `go/bugs/triage/client.go`

Thin HTTP client for the Anthropic Messages API:

```go
package triage

// Client wraps Anthropic Messages API
type Client struct {
    apiKey  string
    model   string
    baseURL string
    http    *http.Client
}

func NewClient() *Client  // reads L8BUGS_ANTHROPIC_API_KEY from env
func (c *Client) Available() bool  // returns true if API key is set
func (c *Client) Complete(systemPrompt, userMessage string) (string, error)
```

- Model default: `claude-sonnet-4-20250514` (fast, cost-effective for triage)
- Timeout: 30 seconds
- Max tokens: 2048 (triage responses are structured JSON, not long)

### 2.2 Create `go/bugs/triage/prompts.go`

Prompt templates for triage operations:

```go
func BuildBugTriagePrompt(bug *l8bugs.Bug, components []string, assignees []*l8bugs.BugsAssignee) string
func BuildFeatureTriagePrompt(feature *l8bugs.Feature, components []string, assignees []*l8bugs.BugsAssignee) string
func BuildDuplicateSearchPrompt(title, description string, candidates []IssueSummary) string
```

Bug triage prompt asks the LLM to return JSON with:
- `priority` (1-4 mapping to Critical/High/Medium/Low)
- `severity` (1-4 mapping to Blocker/Major/Minor/Trivial)
- `component` (from provided component list, or "unknown")
- `assignee_id` (from provided assignee list, or empty)
- `confidence` (0-100)
- `root_cause` (if stack trace present — brief analysis)
- `duplicate_of` (bug ID if duplicate found, or empty)
- `related_ids` (list of related bug IDs)

Feature triage prompt asks for:
- `priority` (1-4)
- `component` (from provided list)
- `assignee_id` (from provided list)
- `confidence` (0-100)
- `breakdown` (effort breakdown — brief)
- `related_ids` (related feature/bug IDs)

### 2.3 Create `go/bugs/triage/parser.go`

JSON response parsing:

```go
type BugTriageResult struct {
    Priority     int32    `json:"priority"`
    Severity     int32    `json:"severity"`
    Component    string   `json:"component"`
    AssigneeID   string   `json:"assignee_id"`
    Confidence   int32    `json:"confidence"`
    RootCause    string   `json:"root_cause"`
    DuplicateOf  string   `json:"duplicate_of"`
    RelatedIDs   []string `json:"related_ids"`
}

type FeatureTriageResult struct { ... }

func ParseBugTriageResponse(response string) (*BugTriageResult, error)
func ParseFeatureTriageResponse(response string) (*FeatureTriageResult, error)
```

---

## Step 3: Triage Orchestrator

### 3.1 Create `go/bugs/triage/triage.go`

The orchestrator coordinates the triage pipeline:

```go
type Triager struct {
    client *Client
    vnic   ifs.IVNic
}

func NewTriager(vnic ifs.IVNic) *Triager
func (t *Triager) TriageBug(bug *l8bugs.Bug) error
func (t *Triager) TriageFeature(feature *l8bugs.Feature) error
```

`TriageBug` flow:
1. Set `triage_status = IN_PROGRESS` via PutEntity
2. Fetch project components and assignees for context
3. Fetch recent open bugs (last 50) for duplicate detection
4. Build prompt with all context
5. Call LLM
6. Parse response
7. Update bug with AI fields + `triage_status = COMPLETED` via PutEntity
8. On error: set `triage_status = FAILED`, `triage_error = err.Error()`

### 3.2 Global triager instance

```go
var globalTriager *Triager

func Initialize(vnic ifs.IVNic) {
    globalTriager = NewTriager(vnic)
}

func Get() *Triager { return globalTriager }
```

---

## Step 4: Hook Triage into Bug/Feature Creation

### 4.1 Extend After hook to support POST

In `go/bugs/common/service_callback.go`, change the After method:

```go
// Before (current):
if (action != ifs.PUT && action != ifs.PATCH) || len(cb.afterActions) == 0 {

// After (new):
if len(cb.afterActions) == 0 {
```

This allows After hooks to fire on POST too. The After hooks themselves decide which actions they care about.

### 4.2 Add triage After hook to Bug callback

In `BugServiceCallback.go`, add an After hook:

```go
return common.NewValidation[l8bugs.Bug]("Bug", ...).
    Require(...).
    StatusTransition(...).
    After(func(entity *l8bugs.Bug, action ifs.Action, vnic ifs.IVNic) error {
        if action != ifs.POST {
            return nil
        }
        triager := triage.Get()
        if triager == nil || !triager.Available() {
            return nil  // AI triage not configured
        }
        // Run async — don't block the response
        go triager.TriageBug(entity)
        return nil
    }).
    Build()
```

### 4.3 Same pattern for Feature callback

---

## Step 5: Initialize Triage on Service Startup

### 5.1 `go/bugs/services/activate_bugs.go`

After activating all services, initialize the triager:

```go
func ActivateBugsServices(creds, dbname string, nic ifs.IVNic) {
    projects.Activate(creds, dbname, nic)
    assignees.Activate(creds, dbname, nic)
    bugs.Activate(creds, dbname, nic)
    features.Activate(creds, dbname, nic)
    sprints.Activate(creds, dbname, nic)
    triage.Initialize(nic)
}
```

---

## Step 6: Desktop UI — Triage Fields & Inbox

### 6.1 Update `l8tracking-enums.js`

Add TRIAGE_STATUS enum:
```javascript
const TRIAGE_STATUS = factory.create([
    ['Unspecified', null, ''],
    ['Pending', 'pending', 'layer8d-status-pending'],
    ['In Progress', 'in-progress', 'layer8d-status-active'],
    ['Completed', 'completed', 'layer8d-status-active'],
    ['Failed', 'failed', 'layer8d-status-terminated'],
    ['Skipped', 'skipped', 'layer8d-status-inactive']
]);
```

### 6.2 Update `l8tracking-columns.js`

Add triage columns to Bug and Feature (after existing columns):
```javascript
...col.enum('triageStatus', 'AI Triage', null, render.triageStatus),
...col.col('aiConfidence', 'AI Conf.'),
```

### 6.3 Update `l8tracking-forms.js`

Update the "AI Analysis" section in Bug form to include new fields:
```javascript
f.section('AI Analysis', [
    ...f.select('triageStatus', 'Triage Status', enums.TRIAGE_STATUS),
    ...f.number('aiConfidence', 'AI Confidence'),
    ...f.select('aiSuggestedPriority', 'AI Suggested Priority', enums.PRIORITY),
    ...f.select('aiSuggestedSeverity', 'AI Suggested Severity', enums.SEVERITY),
    ...f.text('aiSuggestedComponent', 'AI Suggested Component'),
    ...f.reference('aiSuggestedAssigneeId', 'AI Suggested Assignee', 'BugsAssignee'),
    ...f.textarea('aiRootCause', 'AI Root Cause'),
    ...f.text('triageError', 'Triage Error')
])
```

Similar update for Feature form (without severity/root cause, add breakdown).

### 6.4 Triage Inbox — Config View

Add a new service entry to tracking module config for the triage inbox:
```javascript
{ key: 'triage', label: 'Triage Inbox', icon: '\uD83E\uDD16',
  endpoint: '/20/Bug', model: 'Bug',
  supportedViews: ['table'],
  defaultFilter: { triageStatus: 3, status: 1 }  // Completed triage + Open status
}
```

This reuses the Bug service endpoint but with a pre-applied filter showing only AI-triaged issues awaiting human review (status=Open means the human hasn't yet moved it to Triaged).

### 6.5 Update `sections/system.html`

Add "Triage" subnav item and service view container (after Sprints).

---

## Step 7: Mobile UI — Mirror Desktop Changes

### 7.1 Update mobile enums, columns, forms

Same additions as desktop (TRIAGE_STATUS enum, triage columns, updated AI Analysis form section).

### 7.2 Update `m/sections/system.html`

Add "Triage" tab and TRACKING_SERVICES entry.

---

## File Summary

| Action | Count | Files |
|--------|-------|-------|
| Modify (Proto) | 1 | `proto/bugs.proto` (add TriageStatus enum + 4 fields each to Bug/Feature) |
| Regen (Proto) | 1 | `go/types/l8bugs/bugs.pb.go` (regenerated) |
| Create (Go) | 4 | `triage/client.go`, `triage/prompts.go`, `triage/parser.go`, `triage/triage.go` |
| Modify (Go) | 3 | `service_callback.go` (After hook), `BugServiceCallback.go`, `FeatureServiceCallback.go` |
| Modify (Go) | 1 | `activate_bugs.go` (triage init) |
| Modify (Desktop) | 4 | `l8tracking-enums.js`, `l8tracking-columns.js`, `l8tracking-forms.js`, `l8sys-config.js` |
| Modify (Desktop) | 1 | `sections/system.html` |
| Modify (Mobile) | 4 | `l8tracking-enums.js`, `l8tracking-columns.js`, `l8tracking-forms.js`, `m/sections/system.html` |
| Modify (PRD) | 1 | `plans/l8bugs-prd.md` |
| **Total** | **~20** | 4 new Go + ~16 modified |

---

## Verification

```bash
# Proto regeneration
cd proto && ./make-bindings.sh

# Go build
cd go && go build ./... && go vet ./...

# JS syntax
for f in go/bugs/website/web/l8ui/sys/tracking/*.js; do node -c "$f"; done
for f in go/bugs/website/web/m/js/tracking/*.js; do node -c "$f"; done

# Triage package exists
ls go/bugs/triage/*.go

# TriageStatus enum in proto
grep 'TriageStatus' go/types/l8bugs/bugs.pb.go

# After hook fires on POST
grep 'len(cb.afterActions) == 0' go/bugs/common/service_callback.go

# Triage After hook in Bug callback
grep 'triage.Get' go/bugs/bugs/BugServiceCallback.go

# Triage initialized
grep 'triage.Initialize' go/bugs/services/activate_bugs.go

# Desktop: triage inbox in config
grep 'triage' go/bugs/website/web/l8ui/sys/l8sys-config.js

# Desktop: TRIAGE_STATUS enum
grep 'TRIAGE_STATUS' go/bugs/website/web/l8ui/sys/tracking/l8tracking-enums.js

# Mobile: triage tab
grep 'data-service="triage"' go/bugs/website/web/m/sections/system.html
```

---

## Open Questions

1. **Rate limiting**: Should we throttle triage calls if many issues are created in burst? (e.g., bulk import). Suggest: a simple channel-based queue with configurable concurrency (default 2).

2. **Retry on failure**: If the LLM call fails (network, rate limit), should we auto-retry? Suggest: one retry after 5s, then mark as FAILED. User can manually re-trigger.

3. **Cost visibility**: Should we log/surface token usage per triage call? Useful for cost awareness. Suggest: log to stdout for now, dashboard widget in Phase 5.

4. **Triage re-run**: Should editing a bug's title/description re-trigger triage? Suggest: not in Phase 3 — only on creation. Phase 4 can add re-triage on significant edits.
