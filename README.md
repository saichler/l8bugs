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

L8Bugs provides both desktop and mobile web interfaces built with the L8UI component library.

### Desktop (`go/bugs/website/web/`)

Full-featured interface with:
- Table views with sortable/filterable columns
- Detail popups with forms for create/edit
- Status badge rendering with color-coded workflows
- AI triage fields (confidence, estimated effort)
- View switcher (table, chart, kanban, timeline, calendar, gantt)

### Mobile (`go/bugs/website/web/m/`)

Touch-optimized interface with:
- Card-based layouts
- Responsive navigation
- PWA support

Both platforms share the same enum definitions, column definitions, and form definitions through the L8Tracking namespace.

## Project Structure

```
l8bugs/
├── proto/                    # Protobuf definitions
│   └── bugs.proto            # All message types and enums
├── plans/                    # PRD and phase plans
├── go/
│   ├── bugs/                 # Backend implementation
│   │   ├── common/           # Shared utilities, DB connection, service factory
│   │   ├── services/         # Service activation orchestrator
│   │   ├── bugs/             # Bug service + callbacks
│   │   ├── features/         # Feature service + callbacks
│   │   ├── projects/         # Project service + callbacks
│   │   ├── assignees/        # Assignee service + callbacks
│   │   ├── sprints/          # Sprint service + callbacks
│   │   ├── digests/          # Digest service + callbacks
│   │   ├── triage/           # AI triage engine (Anthropic API)
│   │   ├── webhook/          # GitHub/GitLab webhook handlers
│   │   ├── mcp/              # MCP server (tools, handlers, protocol)
│   │   ├── website/          # Web server setup + UI assets
│   │   └── main/             # Backend entry point
│   ├── types/l8bugs/         # Generated protobuf Go types
│   ├── tests/                # Integration test suite
│   │   └── mocks/            # Mock data generators
│   ├── demo/                 # Auto-generated demo binaries and assets
│   └── vendor/               # Vendored dependencies
└── README.md
```

## Running Locally

### Prerequisites

- Go 1.21+
- PostgreSQL (auto-started via `/start-postgres.sh`)
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
| `PREFIX` | `/bugs/` | REST endpoint prefix |
| `ServiceArea` | 20 | Shared across all tracking services |
| Web port | 2883 | HTTPS web server |

## Testing

```bash
cd go
go test -v ./tests/ -count=1 -timeout 300s
```

The test suite runs a full integration test:

1. Drops and recreates the database schema
2. Activates all 6 services
3. Starts a test web server (HTTPS on port 9443)
4. Generates mock data (5 projects, 8 assignees, 20 bugs, 10 features, 6 sprints, 4 digests)
5. Tests service handlers and getters
6. Tests CRUD lifecycle for all entities
7. Tests validation (required fields, auto-generated IDs)
8. Tests business logic (status transitions, terminal states, date validation)
9. Tests GitHub and GitLab webhook processing
10. Tests MCP server tools

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

### Embedded Child Types (not standalone services)

- **Comment** — Issue comments with author type (human/AI/system)
- **Attachment** — File metadata with URLs
- **ActivityEntry** — Change history tracking
- **Label** — Project-scoped issue labels
- **Component** — Software components with default assignees
- **Milestone** — Release milestones with due dates
- **WebhookConfig** — Webhook endpoint and secret per project

## License

Apache License, Version 2.0 — see individual source files for details.

Copyright (c) 2025 Sharon Aicler (saichler@gmail.com)
