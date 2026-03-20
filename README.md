# L8Bugs

AI-First Bug & Feature Tracking System built on the [Layer 8 Ecosystem](https://github.com/saichler).

## Vision

L8Bugs treats AI as the primary operator and humans as reviewers. The core loop:

1. A user (or automated system) reports a bug or feature request
2. AI triages, classifies, prioritizes, detects duplicates, and suggests an assignee
3. AI coding agents (Claude Code, Cursor, etc.) pick up issues, implement fixes, and open PRs
4. Humans review PRs and approve/reject
5. On merge, the issue auto-closes via webhook

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                        Clients                           │
│  Desktop UI  │  Mobile UI  │  REST API  │  MCP Server    │
└──────────────┴─────────────┴────────────┴────────────────┘
                            │
┌──────────────────────────────────────────────────────────┐
│                    L8Bugs Backend                        │
│  Bug  │  Feature  │  Project  │  AI Triage  │  Webhooks  │
└───────┴───────────┴───────────┴─────────────┴────────────┘
                            │
┌──────────────────────────────────────────────────────────┐
│                  Layer 8 Framework                       │
│  ORM  │  Introspection  │  Auth  │  Pub/Sub  │  Storage  │
└───────┴─────────────────┴───────┴──────────┴─────────────┘
                            │
                       PostgreSQL
```

## Services

All services share `ServiceArea=20` and are accessible under the `/bugs/` REST prefix.

| Service | Model | Primary Key | Description |
|---------|-------|-------------|-------------|
| Bug | `Bug` | `bugId` | Issue tracking with status workflows, AI triage, effort estimation |
| Feature | `Feature` | `featureId` | Feature requests with approval workflows |
| Project | `BugsProject` | `projectId` | Project containers with labels, components, milestones, webhook configs |
| Assignee | `BugsAssignee` | `assigneeId` | Human or AI agent assignees per project |
| Sprint | `BugsSprint` | `sprintId` | Sprint planning with capacity and completion tracking |
| Digest | `BugsDigest` | `digestId` | AI-generated project summaries (daily/weekly) |

### Status Workflows

**Bug:** Open &rarr; Triaged &rarr; In Progress &rarr; In Review &rarr; Resolved &rarr; Closed

**Feature:** Proposed &rarr; Triaged &rarr; Approved &rarr; In Progress &rarr; In Review &rarr; Done &rarr; Closed

**Sprint:** Planning &rarr; Active &rarr; Completed

## AI Triage Engine

The triage module (`go/bugs/triage/`) uses the Anthropic API to analyze bugs and features:

- **Priority/Severity suggestion** based on description analysis
- **Component assignment** by matching against project components
- **Duplicate detection** across existing issues
- **Root cause analysis** with suggested investigation steps
- **Effort estimation** with confidence scoring

Set `L8BUGS_ANTHROPIC_API_KEY` to enable AI features.

## MCP Server

L8Bugs includes a [Model Context Protocol](https://modelcontextprotocol.io/) server for AI tool integration. Run it as a subprocess — it communicates over stdin/stdout using JSON-RPC 2.0.

**Entry point:** `go/bugs/mcp/main1/main.go`

### Tools

| Tool | Description |
|------|-------------|
| `list_issues` | List bugs or features with filters (status, priority, assignee, project) |
| `read_issue` | Get full issue details including comments and activity |
| `create_issue` | Create a new bug or feature |
| `update_issue` | Update fields on an existing issue |
| `add_comment` | Add a comment to a bug or feature |
| `search_issues` | Full-text search across title and description |
| `assist_writing` | AI writing assistance (improve descriptions, generate acceptance criteria) |
| `generate_digest` | Generate an AI project summary for a time period |

## Webhook Integrations

L8Bugs processes webhook events from GitHub and GitLab to link commits and auto-transition issues.

### Supported Events

| Provider | Event | Action |
|----------|-------|--------|
| GitHub | Pull Request (merged) | Links merge commit, auto-transitions to Resolved/Done |
| GitHub | Push | Links commit SHA to referenced issue |
| GitLab | Merge Request (merged) | Links merge commit, auto-transitions to Resolved/Done |
| GitLab | Push Hook | Links commit SHA to referenced issue |

Issue references are extracted from commit messages and PR titles/bodies using patterns like `fixes BUG-123` or `closes <issue-id>`.

### Endpoints

- GitHub: `POST /bugs/webhook/github`
- GitLab: `POST /bugs/webhook/gitlab`

Webhook secrets are configured per project in the `WebhookConfig` embedded field.

## UI

L8Bugs provides both desktop and mobile web interfaces built with the L8UI component library. The tracking UI is implemented as a generic L8UI system component under `l8ui/sys/tracking/`, using the `L8Tracking` namespace (desktop) and `MobileL8Tracking` namespace (mobile).

The UI lives under a single "System" sidebar section with Tracking as the default tab.

### Desktop (`go/bugs/website/web/`)

Full-featured interface with:
- Table views with sortable/filterable columns
- Detail popups with forms for create/edit
- Status badge rendering with color-coded workflows
- AI triage fields (confidence, estimated effort)
- View switcher (table, chart, kanban, timeline, calendar, gantt)

### Mobile (`go/bugs/website/web/m/`)

Touch-optimized interface with:
- Card-based layouts with module/service navigation
- Responsive navigation with back buttons and breadcrumbs
- Full CRUD support with mobile-optimized forms
- PWA support

### Shared Tracking Components

Both platforms share enum definitions, column definitions, form definitions, and reference registries:

| File | Desktop Location | Mobile Location |
|------|-----------------|-----------------|
| Enums | `l8ui/sys/tracking/l8tracking-enums.js` | `m/js/tracking/l8tracking-enums.js` |
| Columns | `l8ui/sys/tracking/l8tracking-columns.js` | `m/js/tracking/l8tracking-columns.js` |
| Forms | `l8ui/sys/tracking/l8tracking-forms.js` | `m/js/tracking/l8tracking-forms.js` |
| References | `l8ui/sys/tracking/l8tracking-reference.js` | `m/js/tracking/layer8m-reference-registry-tracking.js` |
| Init | `l8ui/sys/tracking/l8tracking.js` | — |

## Project Structure

```
l8bugs/
├── proto/                       # Protobuf definitions
│   ├── bugs.proto               # All message types and enums
│   └── make-bindings.sh         # Proto compiler script
├── plans/                       # PRD and phase plans
├── go/
│   ├── bugs/                    # Backend implementation
│   │   ├── common/              # Shared utilities, DB connection, service factory
│   │   ├── services/            # Service activation orchestrator
│   │   ├── bugs/                # Bug service + callbacks
│   │   ├── features/            # Feature service + callbacks
│   │   ├── projects/            # Project service + callbacks
│   │   ├── assignees/           # Assignee service + callbacks
│   │   ├── sprints/             # Sprint service + callbacks
│   │   ├── digests/             # Digest service + callbacks
│   │   ├── triage/              # AI triage engine (Anthropic API)
│   │   ├── webhook/             # GitHub/GitLab webhook handlers
│   │   ├── mcp/                 # MCP server (tools, handlers, protocol)
│   │   │   └── main1/           # MCP server entry point
│   │   ├── website/             # Web server + UI assets
│   │   │   └── web/             # Desktop & mobile web assets
│   │   ├── vnet/                # Virtual network entry point
│   │   └── main/                # Backend entry point
│   ├── types/l8bugs/            # Generated protobuf Go types
│   ├── tests/                   # Integration test suite (12 test files)
│   │   └── mocks/               # Mock data generators
│   └── vendor/                  # Vendored dependencies
├── run-local.sh                 # Local development startup
├── test.sh                      # Test runner with coverage
└── README.md
```

## Running Locally

### Prerequisites

- Go 1.25+
- PostgreSQL (auto-started via Docker)
- Layer 8 framework dependencies (vendored)

### Build & Run

```bash
cd go
./run-local.sh
```

This builds three binaries in `demo/`:
- `vnet_demo` — Layer 8 network interface (start first)
- `bugs_demo` — Backend services (start after vnet)
- `ui_demo` — Web server on port 2883 (start last)

### Configuration

| Constant | Value | Description |
|----------|-------|-------------|
| `BUGS_VNET` | 35010 | Service network port |
| `BUGS_LOGS_VNET` | 35015 | Log aggregation network port |
| `PREFIX` | `/bugs/` | REST endpoint prefix |
| `ServiceArea` | 20 | Shared across all tracking services |
| Web port | 2883 | HTTPS web server |

## Testing

```bash
cd go
go test -v ./tests/ -count=1 -timeout 300s
```

The test suite (~4,000 lines across 12 test files) runs full integration tests:

1. Drops and recreates the database schema
2. Activates all 6 services
3. Starts a test web server (HTTPS on port 9443)
4. Generates mock data (5 projects, 8 assignees, 20 bugs, 10 features, 6 sprints, 4 digests)
5. Tests service handlers and getters
6. Tests CRUD lifecycle for all entities
7. Tests validation (required fields, auto-generated IDs)
8. Tests business logic (status transitions, terminal states, date validation)
9. Tests bug and feature status transition workflows
10. Tests sprint lifecycle logic
11. Tests GitHub and GitLab webhook processing
12. Tests MCP server tools

### Test Files

| File | Coverage |
|------|----------|
| `TestAllService_test.go` | Service activation and mock data upload |
| `TestServiceHandlers_test.go` | Handler routing and response validation |
| `TestServiceGetters_test.go` | Query and getter functionality |
| `TestCRUD_test.go` | Full CRUD lifecycle for all entities |
| `TestValidation_test.go` | Required field and input validation |
| `TestBusinessLogic_test.go` | Cross-entity business rules |
| `TestBugTransitions_test.go` | Bug status workflow transitions |
| `TestFeatureTransitions_test.go` | Feature status workflow transitions |
| `TestSprintLogic_test.go` | Sprint lifecycle and capacity logic |
| `TestWebhook_test.go` | GitHub webhook event processing |
| `TestWebhookGitLab_test.go` | GitLab webhook event processing |
| `TestMCP_test.go` | MCP server tool execution |

### Coverage

```bash
cd go
./test.sh
# Opens cover.html with coverage report
```

## Data Model

### Enums

| Enum | Values |
|------|--------|
| Priority | Critical, High, Medium, Low |
| Severity | Blocker, Major, Minor, Trivial |
| BugStatus | Open, Triaged, In Progress, In Review, Resolved, Closed, Reopened, Won't Fix, Duplicate, Cannot Reproduce |
| FeatureStatus | Proposed, Triaged, Approved, In Progress, In Review, Done, Closed, Rejected, Deferred |
| Resolution | Fixed, Won't Fix, Duplicate, Cannot Reproduce, By Design, Obsolete |
| SprintStatus | Planning, Active, Completed |
| TriageStatus | Pending, In Progress, Completed, Failed, Skipped |
| AssigneeType | Human, AI Agent, Team |
| AuthorType | Human, AI, System |
| ProjectStatus | Active, Archived, On Hold |
| ProjectVisibility | Public, Private, Internal |
| MilestoneStatus | Open, Closed |
| DigestPeriod | Daily, Weekly |

### Embedded Child Types (not standalone services)

- **Comment** — Issue comments with author type (human/AI/system)
- **Attachment** — File metadata with URLs
- **ActivityEntry** — Change history tracking
- **Watcher** — Issue subscribers
- **Vote** — Issue voting
- **Label** — Project-scoped issue labels
- **Component** — Software components with default assignees
- **Milestone** — Release milestones with due dates
- **WebhookConfig** — Webhook endpoint and secret per project

## Layer 8 Framework Dependencies

| Package | Purpose |
|---------|---------|
| `l8bus` | Message bus and overlay networking |
| `l8orm` | Object-relational mapping and persistence |
| `l8services` | Service lifecycle management |
| `l8types` | Type system and registry |
| `l8web` | HTTPS web server and REST routing |
| `l8reflect` | Introspection and primary key decorators |
| `l8utils` | Shared utilities and caching |
| `l8srlz` | Serialization |
| `l8test` | Test topology infrastructure |
| `l8ql` | Query language (L8Query) |

## License

Apache License, Version 2.0 — see individual source files for details.

Copyright (c) 2025 Sharon Aicler (saichler@gmail.com)
