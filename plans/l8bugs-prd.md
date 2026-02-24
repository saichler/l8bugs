# L8Bugs - Product Requirements Document

## AI-First Bug & Feature Tracking System

**Version:** 1.0
**Date:** 2026-02-22
**Status:** Draft - Pending Approval

---

## 1. Vision & Purpose

L8Bugs is a bug and feature reporting system designed from the ground up to be **attended by AI**. Unlike traditional issue trackers where AI is an add-on, L8Bugs treats AI as the primary operator and humans as reviewers/approvers.

The core loop:
1. A user (or automated system) reports a bug or feature request
2. AI triages, classifies, prioritizes, detects duplicates, and suggests an assignee
3. AI coding agents (Claude Code, Cursor, etc.) pick up issues, analyze code, create branches, implement fixes, and open PRs
4. Humans review PRs and approve/reject
5. On merge, the issue auto-closes

**Target Users:**
- Software development teams using AI coding agents
- Open-source projects wanting automated triage
- Internal engineering teams wanting to reduce manual issue management overhead

---

## 2. System Architecture

L8Bugs is built on the **Layer8 framework**, which provides:
- Protobuf-based data model with auto-generated Go types
- Service-oriented backend with CRUD operations
- Desktop and mobile UI via the L8UI component library
- Built-in authentication, RBAC, and audit trail

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Clients                            │
│  Desktop UI  │  Mobile UI  │  REST API  │  MCP Server   │
└──────────────┴─────────────┴────────────┴───────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│                   L8Bugs Backend                        │
│  Bug Service │ Project Service │ AI Engine │ Webhooks   │
└──────────────┴─────────────────┴───────────┴────────────┘
                          │
┌─────────────────────────────────────────────────────────┐
│                Layer8 Framework                         │
│  ORM │ Introspection │ Auth │ Pub/Sub │ Storage         │
└──────┴───────────────┴──────┴─────────┴─────────────────┘
```

---

## 3. Data Model (Prime Objects)

The following are the Prime Objects (independent entities with their own services). Each passes the Prime Object Test: independent existence, own lifecycle, directly queryable, no mandatory parent dependency.

### 3.1 Project

A project is the top-level container for organizing issues.

| Field | Type | Description |
|-------|------|-------------|
| project_id | string | Primary key |
| name | string | Project name |
| key | string | Short key for issue numbering (e.g., "L8B") |
| description | string | Project description |
| owner_id | string | User who owns the project |
| status | ProjectStatus | Active, Archived |
| default_assignee_id | string | Default assignee for new issues |
| visibility | ProjectVisibility | Public, Private, Internal |
| created_date | int64 | Timestamp |
| labels | repeated Label | Project-defined label/tag definitions |
| components | repeated Component | Project-defined software components |
| milestones | repeated Milestone | Project-defined release milestones |
| audit_info | erp.AuditInfo | Standard audit metadata |

### 3.2 Bug

A bug report describing a defect in the software.

| Field | Type | Description |
|-------|------|-------------|
| bug_id | string | Primary key |
| project_id | string | Reference to Project |
| bug_number | string | Human-readable number (e.g., "L8B-42") |
| title | string | Short summary |
| description | string | Detailed description (rich text / markdown) |
| steps_to_reproduce | string | Steps to trigger the bug |
| expected_behavior | string | What should happen |
| actual_behavior | string | What actually happens |
| status | BugStatus | See workflow section |
| priority | Priority | Critical, High, Medium, Low |
| severity | Severity | Blocker, Major, Minor, Trivial |
| reporter_id | string | User who reported |
| assignee_id | string | User or AI agent assigned |
| assignee_type | AssigneeType | Human, AI_Agent |
| labels | repeated string | Free-form tags |
| component | string | Affected component/area |
| environment | string | OS, browser, version info |
| stack_trace | string | Error log or stack trace |
| affected_version | string | Version where bug was found |
| fix_version | string | Version where fix will ship |
| due_date | int64 | Target fix date |
| estimated_effort | int32 | AI-estimated effort (story points) |
| ai_confidence | int32 | AI confidence in triage (0-100) |
| ai_suggested_priority | Priority | AI's priority suggestion |
| ai_suggested_component | string | AI's component suggestion |
| ai_root_cause | string | AI's root cause analysis |
| resolution | Resolution | How the bug was resolved |
| resolved_date | int64 | When it was resolved |
| linked_pr_url | string | URL to the fix PR |
| linked_branch | string | Branch name for the fix |
| duplicate_of_id | string | If duplicate, reference to original |
| parent_bug_id | string | Parent bug (for sub-bugs) |
| vote_count | int32 | Number of upvotes |
| watcher_count | int32 | Number of watchers |
| comments | repeated Comment | Discussion thread |
| attachments | repeated Attachment | Files, screenshots |
| related_bug_ids | repeated string | Related (not duplicate) bugs |
| audit_info | erp.AuditInfo | Standard audit metadata |

### 3.3 Feature

A feature request describing new functionality.

| Field | Type | Description |
|-------|------|-------------|
| feature_id | string | Primary key |
| project_id | string | Reference to Project |
| feature_number | string | Human-readable number (e.g., "L8B-43") |
| title | string | Short summary |
| description | string | Detailed description (rich text / markdown) |
| user_story | string | "As a [user], I want [goal], so that [benefit]" |
| acceptance_criteria | string | Conditions for completion |
| status | FeatureStatus | See workflow section |
| priority | Priority | Critical, High, Medium, Low |
| reporter_id | string | User who requested |
| assignee_id | string | User or AI agent assigned |
| assignee_type | AssigneeType | Human, AI_Agent |
| labels | repeated string | Free-form tags |
| component | string | Affected component/area |
| target_version | string | Version where feature will ship |
| due_date | int64 | Target completion date |
| estimated_effort | int32 | AI-estimated effort (story points) |
| ai_confidence | int32 | AI confidence in estimate (0-100) |
| ai_suggested_priority | Priority | AI's priority suggestion |
| ai_breakdown | string | AI-generated sub-task breakdown |
| linked_pr_url | string | URL to the implementation PR |
| linked_branch | string | Branch name |
| parent_feature_id | string | Parent feature (for sub-features) |
| vote_count | int32 | Number of upvotes |
| watcher_count | int32 | Number of watchers |
| comments | repeated Comment | Discussion thread |
| attachments | repeated Attachment | Files, mockups |
| related_feature_ids | repeated string | Related features |
| related_bug_ids | repeated string | Related bugs |
| audit_info | erp.AuditInfo | Standard audit metadata |

### 3.4 Assignee

An assignee represents a person or AI agent who can be assigned to bugs and features.

| Field | Type | Description |
|-------|------|-------------|
| assignee_id | string | Primary key |
| name | string | Display name |
| email | string | Contact email |
| assignee_type | AssigneeType | Human, AI_Agent |
| project_id | string | Reference to Project (optional scope) |
| active | bool | Whether the assignee is available |
| audit_info | erp.AuditInfo | Standard audit metadata |

### 3.5 Sprint

A time-boxed iteration for organizing work.

**Prime Object justification:** Sprint has `project_id` but still passes the Prime Object test because: (1) it has a rich, independent lifecycle (Planning → Active → Completed) with its own workflow rules; (2) users query sprints directly for burndown charts, velocity tracking, and cross-project portfolio planning ("show all active sprints"); (3) sprint-level reports (burndown, scope change) are first-class views, not sub-views of a project. This is analogous to how Bug/Feature also have `project_id` but are independently meaningful work items. Sprint is a first-class agile planning entity, not project configuration.

| Field | Type | Description |
|-------|------|-------------|
| sprint_id | string | Primary key |
| project_id | string | Reference to Project |
| name | string | Sprint name (e.g., "Sprint 23") |
| goal | string | Sprint goal description |
| status | SprintStatus | Planning, Active, Completed |
| start_date | int64 | Sprint start date |
| end_date | int64 | Sprint end date |
| capacity | int32 | Total story points capacity |
| completed_points | int32 | Points completed so far |
| audit_info | erp.AuditInfo | Standard audit metadata |

---

## 4. Embedded Child Types (Not Prime Objects)

These types are embedded within their parent as `repeated` fields. They do NOT get their own services, config entries, nav entries, or standalone forms.

### 4.1 Comment (embedded in Bug, Feature)

A comment is a child of its parent issue. It is meaningless without a parent bug or feature, is always viewed within the parent context, and is never queried across all issues. Per the Prime Object rules, it fails all four criteria for independence and is therefore embedded as `repeated Comment comments` in both Bug and Feature.

| Field | Type | Description |
|-------|------|-------------|
| comment_id | string | Unique ID within parent |
| author_id | string | User or AI agent who commented |
| author_type | AuthorType | Human, AI_Agent, System |
| body | string | Comment text (markdown) |
| is_internal | bool | Internal note (not visible to reporters) |
| attachments | repeated Attachment | Files attached to comment |
| created_date | int64 | Timestamp |
| edited_date | int64 | Last edit timestamp |

### 4.2 Attachment (embedded in Bug, Feature, Comment)

| Field | Type | Description |
|-------|------|-------------|
| attachment_id | string | Unique ID within parent |
| filename | string | Original filename |
| content_type | string | MIME type |
| size | int64 | File size in bytes |
| url | string | Storage URL |
| thumbnail_url | string | Thumbnail for images |
| uploaded_by | string | User ID |
| uploaded_date | int64 | Timestamp |

### 4.3 ActivityEntry (embedded in Bug, Feature)

| Field | Type | Description |
|-------|------|-------------|
| entry_id | string | Unique ID |
| actor_id | string | Who made the change |
| actor_type | AuthorType | Human, AI_Agent, System |
| action | string | What changed (e.g., "status_changed") |
| field_name | string | Which field changed |
| old_value | string | Previous value |
| new_value | string | New value |
| timestamp | int64 | When it happened |

### 4.4 Watcher (embedded in Bug, Feature)

| Field | Type | Description |
|-------|------|-------------|
| user_id | string | Watching user |
| added_date | int64 | When they started watching |

### 4.5 Vote (embedded in Bug, Feature)

| Field | Type | Description |
|-------|------|-------------|
| user_id | string | Voting user |
| vote_date | int64 | When they voted |

### 4.6 Label (embedded in Project)

A reusable label/tag definition. Labels are project-level configuration — they are meaningless without a project, have no independent lifecycle (created/deleted as project config), are never queried across projects ("show all labels" is not useful), and their identity requires project context ("bug" label means different things in different projects). Per the Prime Object rules, Label fails all four criteria and is embedded as `repeated Label labels` in Project.

| Field | Type | Description |
|-------|------|-------------|
| label_id | string | Unique ID within project |
| name | string | Label name |
| color | string | Hex color code |
| description | string | What the label means |

### 4.7 Component (embedded in Project)

A software component definition. Components are project-level configuration — they are meaningless without a project, have no independent lifecycle, are never queried across projects, and "Auth component" means different things in different projects. Per the Prime Object rules, Component fails all four criteria and is embedded as `repeated Component components` in Project.

| Field | Type | Description |
|-------|------|-------------|
| component_id | string | Unique ID within project |
| name | string | Component name |
| description | string | What the component covers |
| lead_id | string | Component lead/owner |
| default_assignee_id | string | Default assignee for issues in this component |

### 4.8 Milestone (embedded in Project)

A target release or deadline. Milestones are project-level organization — they are meaningless without a project ("v2.0" requires project context), have minimal lifecycle (Open/Closed), and are always viewed within project context. Per the Prime Object rules, Milestone fails independence, direct query need, and identity criteria. It is embedded as `repeated Milestone milestones` in Project.

| Field | Type | Description |
|-------|------|-------------|
| milestone_id | string | Unique ID within project |
| name | string | Milestone name (e.g., "v2.0") |
| description | string | Milestone description |
| status | MilestoneStatus | Open, Closed |
| due_date | int64 | Target date |
| completion_percentage | int32 | Calculated from issues |

---

## 5. Enums

### Priority
```
PRIORITY_UNSPECIFIED = 0
PRIORITY_CRITICAL = 1
PRIORITY_HIGH = 2
PRIORITY_MEDIUM = 3
PRIORITY_LOW = 4
```

### Severity (Bugs only)
```
SEVERITY_UNSPECIFIED = 0
SEVERITY_BLOCKER = 1
SEVERITY_MAJOR = 2
SEVERITY_MINOR = 3
SEVERITY_TRIVIAL = 4
```

### BugStatus
```
BUG_STATUS_UNSPECIFIED = 0
BUG_STATUS_OPEN = 1
BUG_STATUS_TRIAGED = 2
BUG_STATUS_IN_PROGRESS = 3
BUG_STATUS_IN_REVIEW = 4
BUG_STATUS_RESOLVED = 5
BUG_STATUS_CLOSED = 6
BUG_STATUS_REOPENED = 7
BUG_STATUS_WONT_FIX = 8
BUG_STATUS_DUPLICATE = 9
BUG_STATUS_CANNOT_REPRODUCE = 10
```

### FeatureStatus
```
FEATURE_STATUS_UNSPECIFIED = 0
FEATURE_STATUS_PROPOSED = 1
FEATURE_STATUS_TRIAGED = 2
FEATURE_STATUS_APPROVED = 3
FEATURE_STATUS_IN_PROGRESS = 4
FEATURE_STATUS_IN_REVIEW = 5
FEATURE_STATUS_DONE = 6
FEATURE_STATUS_CLOSED = 7
FEATURE_STATUS_REJECTED = 8
FEATURE_STATUS_DEFERRED = 9
```

### Resolution
```
RESOLUTION_UNSPECIFIED = 0
RESOLUTION_FIXED = 1
RESOLUTION_WONT_FIX = 2
RESOLUTION_DUPLICATE = 3
RESOLUTION_CANNOT_REPRODUCE = 4
RESOLUTION_BY_DESIGN = 5
RESOLUTION_OBSOLETE = 6
```

### AssigneeType
```
ASSIGNEE_TYPE_UNSPECIFIED = 0
ASSIGNEE_TYPE_HUMAN = 1
ASSIGNEE_TYPE_AI_AGENT = 2
```

### AuthorType
```
AUTHOR_TYPE_UNSPECIFIED = 0
AUTHOR_TYPE_HUMAN = 1
AUTHOR_TYPE_AI_AGENT = 2
AUTHOR_TYPE_SYSTEM = 3
```

> **Note:** AssigneeType and AuthorType are separate enums because issues can only be assigned to humans or AI agents (never to "System"), while comments and activity entries can be authored by the system (e.g., automated status transitions).

### ProjectStatus
```
PROJECT_STATUS_UNSPECIFIED = 0
PROJECT_STATUS_ACTIVE = 1
PROJECT_STATUS_ARCHIVED = 2
```

### ProjectVisibility
```
PROJECT_VISIBILITY_UNSPECIFIED = 0
PROJECT_VISIBILITY_PUBLIC = 1
PROJECT_VISIBILITY_PRIVATE = 2
PROJECT_VISIBILITY_INTERNAL = 3
```

### SprintStatus
```
SPRINT_STATUS_UNSPECIFIED = 0
SPRINT_STATUS_PLANNING = 1
SPRINT_STATUS_ACTIVE = 2
SPRINT_STATUS_COMPLETED = 3
```

### MilestoneStatus
```
MILESTONE_STATUS_UNSPECIFIED = 0
MILESTONE_STATUS_OPEN = 1
MILESTONE_STATUS_CLOSED = 2
```

---

## 6. Workflows

### 6.1 Bug Lifecycle

```
              ┌──────────────────────────────────────────────────┐
              │                                                  │
              ▼                                                  │
┌──────┐  AI triage  ┌─────────┐  assign  ┌─────────────┐      │
│ Open │────────────►│ Triaged │────────►│ In Progress │      │
└──────┘             └─────────┘          └──────┬──────┘      │
   │                                             │              │
   │  ┌─────────────────┐                        │              │
   ├─►│ Cannot Reproduce │               PR opened              │
   │  └─────────────────┘                        │              │
   │  ┌───────────┐                              ▼              │
   ├─►│ Duplicate │                      ┌───────────┐          │
   │  └───────────┘                      │ In Review │          │
   │  ┌───────────┐                      └─────┬─────┘          │
   └─►│ Won't Fix │                            │               │
      └───────────┘                     PR merged               │
                                               │               │
                                               ▼               │
                                        ┌──────────┐   reopen  │
                                        │ Resolved │───────────┘
                                        └─────┬────┘
                                              │
                                        verified / auto
                                              │
                                              ▼
                                        ┌────────┐
                                        │ Closed │
                                        └────────┘
```

**AI-Driven Transitions:**
- Open → Triaged: AI auto-triages on creation (sets priority, severity, component, suggests assignee)
- Triaged → In Progress: When AI agent starts work or human self-assigns
- In Progress → In Review: When a linked PR is opened
- In Review → Resolved: When the linked PR is merged
- Resolved → Closed: Auto-close after verification period (configurable)
- Resolved → Reopened → Open: When the fix is found insufficient

### 6.2 Feature Lifecycle

```
┌──────────┐  AI triage  ┌─────────┐  approve  ┌──────────┐
│ Proposed │────────────►│ Triaged │──────────►│ Approved │
└──────────┘             └─────────┘           └─────┬────┘
     │                        │                      │
     │  ┌──────────┐         │  ┌──────────┐   assign/start
     └─►│ Rejected │         └─►│ Deferred │        │
        └──────────┘            └──────────┘        ▼
                                             ┌─────────────┐
                                             │ In Progress │
                                             └──────┬──────┘
                                                    │
                                              PR opened
                                                    │
                                                    ▼
                                             ┌───────────┐
                                             │ In Review │
                                             └─────┬─────┘
                                                   │
                                              PR merged
                                                   │
                                                   ▼
                                             ┌──────┐
                                             │ Done │
                                             └──┬───┘
                                                │
                                           verified
                                                │
                                                ▼
                                           ┌────────┐
                                           │ Closed │
                                           └────────┘
```

---

## 7. AI Features

### 7.1 AI Triage (Automatic on Issue Creation)

When a bug or feature is created, AI immediately:

1. **Classifies type**: Analyzes title and description to confirm whether this is truly a bug or a feature request (or suggest reclassification)
2. **Infers priority and severity**: Based on language signals ("crashes", "data loss" = Critical; "cosmetic", "minor" = Low), affected component criticality, and reporter history
3. **Suggests component**: Maps keywords and file references to known project components
4. **Suggests assignee**: Based on component ownership, current workload, and expertise matching
5. **Detects duplicates**: Semantic similarity search against open issues. If a match is found with >80% confidence, flags as potential duplicate and links the original
6. **Finds related issues**: Surfaces issues with related symptoms, affected areas, or root causes
7. **Estimates effort**: Predicts story points based on issue description complexity, historical data for similar issues, and component difficulty

All suggestions are auto-applied but shown with confidence scores. Humans can override any AI decision.

### 7.2 AI Root Cause Analysis (Bugs)

When a bug includes a stack trace or error log:

1. Parse the stack trace to identify the crash location (file, function, line)
2. Cross-reference with the linked Git repository to find the relevant code
3. Identify recent commits that modified the crash location (regression detection)
4. Suggest the probable root cause in plain language
5. If confidence is high enough (configurable threshold), auto-assign to the author of the likely-regressing commit

### 7.3 AI Agent Delegation

Issues can be assigned to AI coding agents via MCP (Model Context Protocol):

1. Human or AI triage assigns an issue to an AI agent
2. The agent receives the issue details via MCP `read_issue` tool
3. Agent analyzes the codebase, creates a branch, implements the fix/feature
4. Agent opens a PR and links it to the issue via MCP `update_issue` tool
5. Issue status auto-transitions to "In Review"
6. Human reviews the PR
7. On merge, issue auto-transitions to "Resolved"

**MCP Tools Exposed:**
- `list_issues` - Query issues with filters
- `read_issue` - Get full issue details including comments
- `create_issue` - Create a new bug or feature
- `update_issue` - Update fields (status, assignee, linked PR, etc.)
- `add_comment` - Add a comment to an issue
- `search_issues` - Semantic search across all issues

### 7.4 AI Writing Assistance

- **Auto-suggest templates**: When creating a bug, pre-populate "Steps to Reproduce", "Expected Behavior", "Actual Behavior" structure
- **Improve descriptions**: Suggest clearer wording, add missing context
- **Generate acceptance criteria**: From feature description text, produce testable acceptance criteria
- **Summarize threads**: Condense long comment threads into key decisions and action items

### 7.5 AI-Powered Search

- **Natural language queries**: "Show all critical bugs in the auth module from the last week"
- **Semantic search**: Find issues by meaning, not just keywords ("login doesn't work" matches "authentication failure on signin page")
- **Conversational Q&A**: "What's blocking the v2.0 milestone?" → AI aggregates blockers and summarizes status

### 7.6 AI Analytics & Predictions

- **Sprint completion prediction**: "At current velocity, Sprint 23 will complete 85% of planned work"
- **Bug trend forecasting**: "The payments component has a rising bug trend — investigate"
- **Anomaly detection**: Alert when bug creation rate, resolution time, or reopening rate deviates significantly from norms
- **Daily/weekly digests**: AI-generated project summaries sent to stakeholders

---

## 8. UI Requirements

### 8.1 Views

| View | Description | Primary Use |
|------|-------------|-------------|
| **Table** | Sortable, filterable list of issues | Default issue browsing |
| **Board (Kanban)** | Cards organized by status columns | Visual workflow tracking |
| **Timeline** | Horizontal bars showing issue lifecycles | Sprint/milestone planning |
| **Calendar** | Issues on a calendar by due date | Deadline management |
| **Dashboard** | Configurable widgets and charts | Project overview & metrics |
| **Triage Inbox** | AI-processed queue of new issues | Review AI triage suggestions |
| **Detail View** | Full issue detail with comments and activity | Issue editing and discussion |

### 8.2 Dashboard Widgets

- Open bugs by priority (pie chart)
- Bug creation vs. resolution over time (line chart)
- Sprint burndown (bar chart)
- Mean time to resolution (KPI card)
- AI triage accuracy (KPI card)
- Top components by bug count (bar chart)
- Assignee workload distribution (bar chart)
- Upcoming milestones (timeline)

### 8.3 Desktop / Mobile Parity

All features must have functional parity between desktop and mobile:
- Desktop: Full L8UI component library (tables, forms, charts, kanban, gantt, popups)
- Mobile: Touch-optimized card-based navigation with equivalent functionality

---

## 9. Integration Points

### 9.1 Git Integration

- **Commit linking**: Mention issue number in commit message (e.g., "Fix L8B-42") to auto-link
- **Branch creation**: Create a branch from an issue with a standardized naming convention
- **PR linking**: Auto-detect PRs that reference issue numbers
- **Auto-transition on merge**: When a linked PR merges, transition issue to Resolved

### 9.2 Webhook System

**Outbound webhooks** (L8Bugs → external systems):
- Issue created, updated, status changed, assigned, commented
- Configurable per project with URL, secret, and event filter

**Inbound webhooks** (external systems → L8Bugs):
- GitHub/GitLab PR events → update issue status
- CI/CD build results → attach to linked issues
- Chat commands → create/update issues

### 9.3 REST API

Full CRUD API for all Prime Objects, following Layer8 patterns:
- `POST /erp/{serviceArea}/{serviceName}` - Create
- `PUT /erp/{serviceArea}/{serviceName}` - Update
- `GET /erp/{serviceArea}/{serviceName}` - List/Query (with L8Query for paging, filtering, sorting)
- `DELETE /erp/{serviceArea}/{serviceName}` - Delete

### 9.4 Notifications

| Channel | Trigger | Content |
|---------|---------|---------|
| In-app | All events | Real-time bell notifications |
| Email | Configurable per user | Instant or daily digest |
| Webhook | Configurable per project | JSON event payload |

**AI-enhanced notifications:**
- Smart batching: Group related changes into a single notification
- Priority-based delivery: Critical issues notify immediately; low-priority batch into digests
- AI daily summary: Auto-generated project status email

---

## 10. Access Control

### Roles

| Role | Permissions |
|------|-------------|
| **Admin** | Full access: manage projects, users, settings |
| **Manager** | Manage project settings, sprints, milestones; all issue operations |
| **Developer** | Create, edit, assign issues; comment; link PRs |
| **Reporter** | Create issues, add comments, vote, watch |
| **Viewer** | Read-only access to issues and dashboards |

### AI Access Scoping

- AI agents inherit the permissions of the user who assigned them
- Projects can opt-out of AI triage entirely
- Security-sensitive issues can be marked as "AI-restricted" to prevent AI access
- All AI actions are logged in the audit trail with `author_type = AI_AGENT`

---

## 11. Metrics & Reporting

### Built-in Metrics

| Metric | Description |
|--------|-------------|
| **Velocity** | Story points completed per sprint |
| **Cycle Time** | Time from In Progress to Resolved |
| **Lead Time** | Time from Open to Resolved |
| **Throughput** | Issues completed per time period |
| **MTTR** | Mean time to resolution |
| **Bug Escape Rate** | Bugs found in production vs. in development |
| **Reopen Rate** | Percentage of resolved issues that get reopened |
| **AI Triage Accuracy** | Percentage of AI suggestions accepted by humans |
| **AI Agent Success Rate** | Percentage of AI-assigned issues resolved without human rework |
| **Burndown** | Work remaining vs. time in a sprint |

### Custom Reporting

- Filter-based query builder for ad-hoc reports
- Export to CSV
- Dashboard widgets for any metric
- Natural language queries via AI

---

## 12. Service Area & Service Names

L8Bugs uses **ServiceArea = 20**. All services share this area. (L8ERP uses 30-130; 20 is unused.)

| Service | ServiceName (max 10 chars) | Model | PrimaryKey |
|---------|---------------------------|-------|------------|
| Projects | Project | BugsProject | ProjectId |
| Assignees | Assignee | BugsAssignee | AssigneeId |
| Bugs | Bug | Bug | BugId |
| Features | Feature | Feature | FeatureId |
| Sprints | Sprint | BugsSprint | SprintId |

**Not services (embedded children):** Comment, Attachment, ActivityEntry, Watcher, and Vote are embedded as `repeated` fields in Bug/Feature. Label, Component, and Milestone are embedded as `repeated` fields in Project. All are managed through the parent's CRUD operations and displayed as inline tables within the parent's detail form.

### Endpoint Examples

```
POST /erp/20/Project     # Create a project
GET  /erp/20/Bug          # List/query bugs (with L8Query)
PUT  /erp/20/Feature      # Update a feature
```

---

## 13. Implementation Reference (based on l8erp)

L8Bugs follows the exact patterns established in `../l8erp`. The L8UI component library (already present in `go/bugs/website/web/l8ui/`) provides all behavioral UI — module files contain **only configuration data**.

### Directory Structure

```
l8bugs/
├── proto/
│   ├── bugs.proto                    # All messages, enums, List types
│   └── make-bindings.sh              # Proto compilation (update for bugs.proto)
├── go/
│   ├── go.mod                        # Module definition + Layer8 dependencies
│   ├── types/bugs/                   # Generated .pb.go files (from make-bindings.sh)
│   ├── bugs/
│   │   ├── common/                   # Shared service utilities (ActivateService, etc.)
│   │   ├── projects/                 # ProjectService.go + ProjectServiceCallback.go
│   │   ├── bugs/                     # BugService.go + BugServiceCallback.go
│   │   ├── features/                 # FeatureService.go + FeatureServiceCallback.go
│   │   ├── sprints/                  # SprintService.go + SprintServiceCallback.go (Phase 2)
│   │   └── website/web/              # UI assets
│   │       ├── app.html              # Desktop entry point (single "System" sidebar section)
│   │       ├── js/
│   │       │   ├── sections.js       # Section loader (system only)
│   │       │   └── app.js            # App bootstrap
│   │       ├── sections/
│   │       │   └── system.html       # System section (Tracking, Health, Security, Modules, Logs)
│   │       ├── m/                    # Mobile UI
│   │       │   ├── app.html
│   │       │   ├── sections/
│   │       │   │   └── system.html   # Mobile system section (Tracking, Health, Security, Modules)
│   │       │   └── js/
│   │       │       ├── tracking/     # Mobile tracking data files
│   │       │       │   ├── l8tracking-enums.js     # MobileL8Tracking namespace
│   │       │       │   ├── l8tracking-columns.js
│   │       │       │   ├── l8tracking-forms.js
│   │       │       │   └── layer8m-reference-registry-tracking.js
│   │       │       ├── app-core.js
│   │       │       ├── mobile-config-bugs.js       # Reference picker config
│   │       │       └── layer8m-nav-config-bugs.js  # Nav config (system only)
│   │       └── l8ui/                 # Shared component library
│   │           └── sys/
│   │               ├── l8sys-config.js     # System module config (tracking + security)
│   │               ├── l8sys-init.js       # System module init
│   │               ├── tracking/           # Generic tracking sub-module (l8ui component)
│   │               │   ├── l8tracking-enums.js     # L8Tracking namespace
│   │               │   ├── l8tracking-columns.js
│   │               │   ├── l8tracking-forms.js
│   │               │   ├── l8tracking.js           # Verification entry point
│   │               │   └── l8tracking-reference.js # Reference registry
│   │               ├── security/           # Security sub-module
│   │               ├── health/             # Health sub-module
│   │               ├── modules/            # Modules sub-module
│   │               └── logs/               # Logs sub-module
│   └── tests/mocks/                  # Mock data generation (Phase 2+)
```

**Note:** Bug Tracking is not a standalone sidebar section. It is integrated into the System section as the "Tracking" sub-module tab (the default tab), alongside Health, Security, Modules, and Logs. The tracking data files live under `l8ui/sys/tracking/` as a generic l8ui component, following the same pattern as `l8ui/sys/security/`.

### Shared Proto Dependency

L8Bugs imports `erp-common.proto` (for `erp.AuditInfo`, `erp.Money`, etc.) and `api.proto` (for `l8api.L8MetaData`). These are Layer8 framework proto files downloaded automatically by `make-bindings.sh` during the binding generation step. The script fetches them from the Layer8 framework repositories, so l8bugs does not need to vendor or copy them — they are resolved at proto compilation time, the same way l8erp handles them.

### Proto File Pattern (following l8erp)

```protobuf
syntax = "proto3";
package l8bugs;
option go_package = "./types/bugs";
import "erp-common.proto";
import "api.proto";

// @PrimeObject
message BugsProject {
    string project_id = 1;
    // ... scalar fields ...
    repeated Label labels = 20;            // Embedded child (project config)
    repeated Component components = 21;    // Embedded child (project config)
    repeated Milestone milestones = 22;    // Embedded child (project org)
    erp.AuditInfo audit_info = 30;
}

message BugsProjectList {
    repeated BugsProject list = 1;         // MUST be named "list"
    l8api.L8MetaData metadata = 2;         // MUST include metadata
}

// @PrimeObject
message Bug {
    string bug_id = 1;
    string project_id = 2;
    // ... all fields ...
    repeated Comment comments = 25;        // Embedded child
    repeated Attachment attachments = 26;   // Embedded child
    erp.AuditInfo audit_info = 30;
}

message BugList {
    repeated Bug list = 1;                  // MUST be named "list"
    l8api.L8MetaData metadata = 2;         // MUST include metadata
}

// Child types — no List type, no service
message Comment {
    string comment_id = 1;
    // ... fields ...
}
message Label {
    string label_id = 1;
    // ... fields ...
}
message Component {
    string component_id = 1;
    // ... fields ...
}
message Milestone {
    string milestone_id = 1;
    // ... fields ...
}
```

### Go Service Pattern (following l8erp)

```go
// go/bugs/bugs/BugService.go
package bugs

const (
    ServiceName = "Bug"           // max 10 chars
    ServiceArea = byte(20)        // shared across all l8bugs services
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
    common.ActivateService[l8bugs.Bug, l8bugs.BugList](
        common.ServiceConfig{
            ServiceName: ServiceName,
            ServiceArea: ServiceArea,
            PrimaryKey:  "BugId",
            Callback:    newBugServiceCallback(),
            Transactional: true,
        }, creds, dbname, vnic)
}
```

```go
// go/bugs/bugs/BugServiceCallback.go
func newBugServiceCallback() ifs.IServiceCallback {
    return common.NewServiceCallback("Bug",
        func(e *l8bugs.Bug) {
            common.GenerateID(&e.BugId)   // Auto-generate ID on POST
        },
        validate)
}

func validate(bug *l8bugs.Bug, vnic ifs.IVNic) error {
    // Required → Enum → Date → Business rules → References
    return nil
}
```

### UI Type Registration Pattern (following l8erp)

```go
// In main.go or shared.go
func registerBugsTypes(resources ifs.IResources) {
    common.RegisterType[l8bugs.BugsProject, l8bugs.BugsProjectList](resources, "ProjectId")
    common.RegisterType[l8bugs.Bug, l8bugs.BugList](resources, "BugId")
    common.RegisterType[l8bugs.Feature, l8bugs.FeatureList](resources, "FeatureId")
    common.RegisterType[l8bugs.BugsSprint, l8bugs.BugsSprintList](resources, "SprintId")
    // Note: Label, Component, Milestone are embedded children — NOT registered
}
```

### Desktop UI Config Pattern (System Section Integration)

Bug Tracking is configured as a sub-module of the System section in `l8ui/sys/l8sys-config.js`, not as a standalone module. The tracking data files (`l8tracking-enums.js`, `l8tracking-columns.js`, `l8tracking-forms.js`) follow the same pattern as `l8ui/sys/security/` — they are generic l8ui components reusable by any Layer8 app.

```javascript
// l8ui/sys/l8sys-config.js — tracking is the FIRST module (default tab)
L8Sys.modules = {
    tracking: {
        label: 'Tracking', icon: '🐛',
        services: [
            { key: 'bugs', label: 'Bugs', icon: '🐛',
              endpoint: '/20/Bug', model: 'Bug',
              supportedViews: ['table', 'kanban'] },
            { key: 'features', label: 'Features', icon: '✨',
              endpoint: '/20/Feature', model: 'Feature',
              supportedViews: ['table', 'kanban'] },
            { key: 'projects', label: 'Projects', icon: '📁',
              endpoint: '/20/Project', model: 'BugsProject' }
        ]
    },
    health: { ... },
    security: { ... },
    modules: { ... },
    logs: { ... }
};
L8Sys.submodules = ['L8Security', 'L8Tracking'];
```

```javascript
// l8ui/sys/l8sys-init.js — defaults to tracking tab
Layer8DModuleFactory.create({
    namespace: 'L8Sys',
    defaultModule: 'tracking',
    defaultService: 'bugs',
    sectionSelector: 'tracking',       // MUST match defaultModule
    initializerName: 'initializeL8Sys',
    requiredNamespaces: ['L8Security', 'L8Tracking']
});
```

### Section HTML Pattern

The System section (`sections/system.html`) contains static HTML with module tabs (Tracking, Health, Security, Modules, Logs). Tracking is the first/active tab with subnav for Bugs, Features, and Projects.

```html
<!-- sections/system.html (relevant excerpt) -->
<div class="l8-module-tabs">
    <button class="l8-module-tab active" data-module="tracking">
        <span class="tab-icon">🐛</span><span class="tab-label">Tracking</span>
    </button>
    <button class="l8-module-tab" data-module="health">...</button>
    <button class="l8-module-tab" data-module="security">...</button>
    <!-- ... -->
</div>
<div class="l8-module-content active" data-module="tracking">
    <nav class="l8-subnav">
        <a class="l8-subnav-item active" data-service="bugs">Bugs</a>
        <a class="l8-subnav-item" data-service="features">Features</a>
        <a class="l8-subnav-item" data-service="projects">Projects</a>
    </nav>
    <!-- Container IDs follow {moduleKey}-{serviceKey}-table-container -->
    <div class="l8-service-view active" data-service="bugs">
        <div class="l8-table-container" id="tracking-bugs-table-container"></div>
    </div>
    <!-- ... -->
</div>
```

### Reference Registry Pattern

```javascript
// l8ui/sys/tracking/l8tracking-reference.js
const ref = window.Layer8RefFactory;
Layer8DReferenceRegistry.register({
    ...ref.simple('BugsProject', 'projectId', 'name', 'Project'),
    ...ref.simple('Bug', 'bugId', 'title', 'Bug'),           // For duplicate_of_id, parent_bug_id lookups
    ...ref.simple('Feature', 'featureId', 'title', 'Feature'), // For parent_feature_id lookups
    ...ref.simple('BugsSprint', 'sprintId', 'name', 'Sprint'),
    // Note: Label, Component, Milestone are embedded children of Project — not registered
});
```

---

## 14. Implementation Phases

### Phase 1: Foundation ✅ Complete
1. **Proto**: Define `bugs.proto` with all Prime Object messages, embedded child types, enums, and List types. Update `make-bindings.sh` and generate Go types.
2. **Go services**: Implement Project, Bug, Feature services following l8erp pattern:
   - `*Service.go` (Activate, ServiceName, ServiceArea=20)
   - `*ServiceCallback.go` (GenerateID on POST, validate chain)
   - Type registration in `main.go`
3. **Desktop UI** — Tracking integrated into System section as generic l8ui component:
   - `l8ui/sys/l8sys-config.js` — tracking added as first module with 3 services (bugs, features, projects)
   - `l8ui/sys/tracking/l8tracking-enums.js` — `Layer8EnumFactory` calls for all enums (`L8Tracking` namespace)
   - `l8ui/sys/tracking/l8tracking-columns.js` — `Layer8ColumnFactory` calls (field names verified against .pb.go)
   - `l8ui/sys/tracking/l8tracking-forms.js` — `Layer8FormFactory` calls (Comment as `f.inlineTable()` in Bug/Feature; Label, Component, Milestone as `f.inlineTable()` in Project)
   - `l8ui/sys/tracking/l8tracking.js` — verification entry point
   - `l8ui/sys/tracking/l8tracking-reference.js` — `Layer8RefFactory` registrations
   - `l8ui/sys/l8sys-init.js` — defaults to tracking tab, `requiredNamespaces: ['L8Security', 'L8Tracking']`
   - `sections/system.html` — Tracking tab (active) with Bugs/Features/Projects subnav
   - `app.html` — single "System" sidebar section, tracking scripts in SYS Module block
   - `js/sections.js` — system section only
4. **Mobile UI** (per Mobile Parity rule):
   - `m/js/tracking/l8tracking-enums.js`, `l8tracking-columns.js`, `l8tracking-forms.js` (`MobileL8Tracking` namespace)
   - `m/js/tracking/layer8m-reference-registry-tracking.js` — mobile reference registry
   - `m/sections/system.html` — Tracking tab (active) with service sub-tabs and inline `loadTrackingService()` function
   - `m/app.html` — single "System" sidebar section, tracking scripts
   - `m/js/app-core.js` — defaults to system section

### Phase 2: Workflow & Core Features ✓
- ✓ Sprint Go service implemented (`go/bugs/sprints/SprintService.go`, `SprintServiceCallback.go`)
- ✓ Status workflow transition rules added to Bug, Feature, and Sprint ServiceCallbacks
- ✓ Kanban board viewConfig added to Bug and Feature services (laneField, lanes, cardTitle, cardSubtitle, cardFields)
- ✓ Activity Log inline table added to Bug and Feature forms (desktop + mobile)
- ✓ Sprint service registered in type system, activated in service bootstrap
- ✓ Desktop + mobile UI complete: config, HTML, enums, columns, forms, reference registry

### Phase 3: AI Triage ✓
- AI classification (bug vs. feature, priority, severity, component)
- Duplicate detection via semantic similarity
- Related issue surfacing
- Triage inbox UI
- AI confidence scoring

**Completed:**
- ✓ Protobuf additions: TriageStatus enum, ai_suggested_severity/assignee_id/triage_status/triage_error fields on Bug and Feature
- ✓ Triage Go package: thin Anthropic HTTP client (client.go), prompt templates (prompts.go), JSON parser (parser.go), orchestrator (triage.go), context fetchers (context.go)
- ✓ After hook extended to fire on POST; Bug and Feature callbacks spawn async triage goroutine on creation
- ✓ Triage initialized on service startup via triage.Initialize(nic)
- ✓ Desktop UI: TRIAGE_STATUS enum, triage columns in Bug/Feature tables, expanded AI Analysis form sections, Triage Inbox subnav + service view
- ✓ Mobile UI: TRIAGE_STATUS enum, triage columns, expanded AI Analysis forms, Triage tab + TRACKING_SERVICES entry

### Phase 4: AI Agent Integration ✓
- MCP server implementation
- Agent delegation workflow
- Git integration (commit/PR linking, auto-transitions)
- AI root cause analysis for bugs with stack traces

**Completed:**
- ✓ Protobuf additions: repository_url/webhook_secret on BugsProject, linked_commit_sha on Bug and Feature
- ✓ MCP Server: Standalone stdio binary (mcp/protocol.go, server.go, tools.go, handlers.go, main1/main.go) with 6 tools: list_issues, read_issue, create_issue, update_issue, add_comment, search_issues — JSON-RPC 2.0 over stdin/stdout
- ✓ Git Integration: GitHub webhook handler (webhook/webhook.go, webhook/parser.go) registered on existing web server; handles PR merge (auto-transition In Review → Resolved/Done) and push events (commit SHA linking); HMAC-SHA256 signature verification
- ✓ AI Root Cause Analysis: Stack trace parser (triage/rootcause.go) supporting Go, Java, Python, JavaScript; enhanced LLM prompt with structured stack frames; integrated into triage flow — when bug has stack trace, runs secondary root cause analysis after basic triage
- ✓ Desktop UI: repositoryUrl/webhookSecret on BugsProject form, linkedCommitSha on Bug/Feature forms, repositoryUrl column on BugsProject
- ✓ Mobile UI: Parity with desktop — same form and column additions

### Phase 5: Analytics & Polish ✓
- **Dashboard module** (desktop + mobile) with 6 KPI widgets (Open Bugs, Open Features, Resolved This Week, AI Triage Accuracy, Active Sprints, Overdue Items) and 4 charts (Bugs by Priority/Status, Features by Status, Top Components) using Layer8DWidget and Layer8DChart
- Dashboard is the new default tab in System section (before Tracking)
- **Enhanced view types**: Bug and Feature services now support timeline and calendar views; Sprint service supports gantt view
- **AI writing assistance** backend (`triage/writer.go`) with 4 actions: suggest_steps, improve_description, generate_acceptance_criteria, summarize_comments
- **MCP assist_writing tool** (7th tool) exposes AI writing to coding agents
- Desktop/mobile parity for all dashboard features

### Phase 6: Digests, Webhooks & Effort Tracking ✓
- **BugsDigest Prime Object**: New service (ServiceName="Digest", ServiceArea=20) for AI-generated project summaries with period, date range, summary, key metrics, blockers, and action items
- **Outbound webhook configuration**: WebhookConfig embedded child added to BugsProject (url, secret, events, active flag) with inline table in project form
- **AI effort estimation**: New `triage/effort.go` estimates story points (1-13 Fibonacci) with confidence percentage; integrated into triage flow for both Bug and Feature
- **AI digest generation**: New `triage/digest.go` generates project summaries from bug/feature data using LLM; auto-saves as BugsDigest records
- **MCP generate_digest tool** (8th tool): Allows AI coding agents to generate daily/weekly/custom digests for any project
- **New proto additions**: DigestPeriod enum, WebhookEventType enum, WebhookConfig message, BugsDigest/BugsDigestList messages, actual_effort/ai_estimated_effort/ai_effort_confidence fields on Bug, actual_effort/ai_estimated_effort on Feature
- **Desktop UI**: Digest service in tracking config, DIGEST_PERIOD/WEBHOOK_EVENT_TYPE enums, BugsDigest columns/form, effort columns/fields on Bug/Feature, webhook inline table on BugsProject
- **Mobile UI**: Full parity — same enums, columns, forms, reference config additions
- **MCP handlers refactored**: Extracted enum parsers and digest handler into `parsers.go` to keep handlers.go under 500 lines

### Phase 7: Advanced Features (Future)
- Natural language search and Q&A
- Inbound webhook system (GitHub/GitLab PR events, CI/CD results)
- Import from Jira/GitHub/Linear
- AI effort estimation learning feedback loop (compare actual vs estimated)
- Sprint completion prediction and anomaly detection

---

## 15. Non-Functional Requirements

| Requirement | Target |
|-------------|--------|
| Page load time | < 2 seconds |
| API response time (p95) | < 500ms |
| AI triage latency | < 5 seconds per issue |
| Concurrent users | 100+ per instance |
| Data retention | Configurable, default unlimited |
| Availability | 99.9% uptime target |
| Browser support | Chrome, Firefox, Safari, Edge (latest 2 versions) |
| Mobile support | iOS Safari, Android Chrome |

---

## 16. Research Sources

This PRD was informed by analysis of the following tools:

- **Jira** (Atlassian): Enterprise-grade customizability, Atlassian Intelligence AI, workflow engine
- **Linear**: Opinionated speed, AI-first triage intelligence, MCP server for agent integration
- **GitHub Issues**: Native Git integration, Agentic Workflows for AI automation, sub-issues
- **Azure DevOps**: Template-based workflows, Power BI analytics, Delivery Plans
- **Bugzilla**: Classic lifecycle management, severity/priority separation, open source
- **YouTrack** (JetBrains): JavaScript-based custom workflows, AI assistant, MCP support
- **Shortcut**: Balanced simplicity, Korey AI for story generation
- **Plane.so**: Open-source alternative to Linear, self-hosted
