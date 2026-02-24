# Phase 5: Analytics & Polish

## Context

Phases 1-4 are complete: proto definitions, Go services, desktop+mobile UI, Sprint service, status workflow transitions, Kanban boards, AI triage pipeline, MCP server, GitHub webhook integration, and AI root cause analysis. Phase 5 adds a dashboard with KPI widgets and charts, additional view types for existing services, and AI writing assistance.

---

## Step 1: Dashboard Module (Desktop)

Add a "Dashboard" tab to the System section as the first module tab, before "Tracking". The dashboard renders KPI widgets and charts using existing l8ui components (`Layer8DWidget`, `Layer8DChart`). No new proto or Go changes needed — the dashboard fetches data from existing Bug/Feature/Sprint endpoints and computes metrics client-side.

### 1.1 `go/bugs/website/web/l8ui/sys/dashboard/l8dashboard.css` (~80 lines)

Dashboard layout styles:
```css
.l8-dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
    padding: 16px;
}
.l8-dashboard-charts {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    gap: 16px;
    padding: 0 16px 16px;
}
.l8-dashboard-chart-card {
    background: var(--layer8d-bg-white);
    border: 1px solid var(--layer8d-border);
    border-radius: 8px;
    padding: 16px;
    min-height: 280px;
}
.l8-dashboard-chart-title {
    font-size: 14px;
    font-weight: 600;
    color: var(--layer8d-text-dark);
    margin-bottom: 12px;
}
```

### 1.2 `go/bugs/website/web/l8ui/sys/dashboard/l8dashboard.js` (~300 lines)

Dashboard data fetcher and renderer. Uses `Layer8DConfig` to make fetch calls to Bug/Feature/Sprint endpoints, then renders:

**KPI Widgets (top row, 4-6 cards):**
- Open Bugs (count of bugs with status 1-4)
- Open Features (count of features with status 1-5)
- Resolved This Week (bugs resolved in last 7 days)
- AI Triage Accuracy (% of triage suggestions with confidence > 70)
- Active Sprints (sprints with status=2)
- Overdue Items (bugs/features past due_date)

Each widget uses `Layer8DWidget.render()` with a trend arrow (comparing to previous period) and a sparkline where applicable.

**Charts (below widgets, 2-column grid):**
- **Bugs by Priority** — pie chart from Bug data grouped by `priority`
- **Bugs by Status** — bar chart from Bug data grouped by `status`
- **Creation vs Resolution** — line chart with two series (created_date vs resolved_date over last 30 days)
- **Top Components** — horizontal bar chart of bug count per `component`

Each chart is rendered inside a `.l8-dashboard-chart-card` div using `Layer8DChart` directly (not via the view factory — these are standalone chart instances).

```go
// Pseudo-code for the chart rendering pattern:
const container = document.getElementById('chart-bugs-by-priority');
const chart = new Layer8DChart(container, {
    chartType: 'pie',
    data: bugsByPriorityData,
    title: 'Bugs by Priority'
});
```

**Public API:**
```javascript
window.L8Dashboard = {
    initialize(),      // Called when dashboard tab is activated
    refresh(),         // Re-fetch all data and re-render
    _fetchBugs(),      // GET /bugs/20/Bug
    _fetchFeatures(),  // GET /bugs/20/Feature
    _fetchSprints(),   // GET /bugs/20/Sprint
    _computeKPIs(bugs, features, sprints),
    _renderWidgets(kpis),
    _renderCharts(bugs, features)
};
```

### 1.3 Desktop HTML + Config Changes

**`sections/system.html`** — Add Dashboard tab (first position) and dashboard content div:
```html
<!-- Before tracking tab -->
<button class="l8-module-tab active" data-module="dashboard">
    <span class="tab-icon">&#x1F4CA;</span>
    <span class="tab-label">Dashboard</span>
</button>
<!-- Tracking tab becomes non-active -->
<button class="l8-module-tab" data-module="tracking">...</button>
```

Dashboard content area:
```html
<div class="l8-module-content active" data-module="dashboard">
    <div class="l8-dashboard-grid" id="dashboard-widgets"></div>
    <div class="l8-dashboard-charts" id="dashboard-charts"></div>
</div>
<!-- Tracking becomes non-active -->
<div class="l8-module-content" data-module="tracking">...</div>
```

**`l8sys-config.js`** — Add 'dashboard' module entry (before 'tracking'):
```javascript
'dashboard': {
    label: 'Dashboard',
    icon: '\uD83D\uDCCA',
    services: []  // No services — custom rendering
}
```

**`l8sys-init.js`** — Initialize dashboard when tab is activated:
```javascript
// In the extended initializeL8Sys:
if (window.L8Dashboard) L8Dashboard.initialize();
```

**`app.html`** — Add CSS and JS includes:
```html
<link rel="stylesheet" href="l8ui/sys/dashboard/l8dashboard.css">
<script src="l8ui/sys/dashboard/l8dashboard.js"></script>
```

---

## Step 2: Enhanced View Types for Tracking Services

Add `timeline` and `calendar` views to Bug and Feature services so users can see issues on a timeline or calendar by due date.

### 2.1 `l8sys-config.js` — Update supportedViews

**Bugs:** Add `timeline` and `calendar`:
```javascript
supportedViews: ['table', 'kanban', 'timeline', 'calendar'],
viewConfig: {
    // existing kanban config...
    // Timeline config (uses createdDate):
    dateField: 'createdDate',
    titleField: 'title',
    actorField: 'assigneeId',
    // Calendar config (uses dueDate):
    calendarDateField: 'dueDate'
}
```

**Features:** Same additions:
```javascript
supportedViews: ['table', 'kanban', 'timeline', 'calendar'],
```

**Sprints:** Add `gantt` view:
```javascript
supportedViews: ['table', 'gantt'],
viewConfig: {
    startDateField: 'startDate',
    endDateField: 'endDate',
    titleField: 'name',
    progressField: 'completedPoints'
}
```

---

## Step 3: Dashboard Module (Mobile)

### 3.1 `m/sections/system.html` — Add Dashboard tab and content

Add "Dashboard" as the first tab, before "Tracking". Add a dashboard content div:
```html
<button class="section-tab active" data-module="dashboard">Dashboard</button>
<button class="section-tab" data-module="tracking">Tracking</button>
```

Dashboard content area:
```html
<div id="system-dashboard-content">
    <div class="l8-dashboard-grid" id="mobile-dashboard-widgets"></div>
    <div class="l8-dashboard-charts" id="mobile-dashboard-charts"></div>
</div>
```

### 3.2 `m/sections/system.html` — Update switchSystemModule

Add `dashboard` case in the module switcher function, calling the same `L8Dashboard` module (shared with desktop):
```javascript
} else if (moduleKey === 'dashboard') {
    dashboardContent.style.display = '';
    if (window.L8Dashboard) L8Dashboard.initialize('mobile-dashboard-widgets', 'mobile-dashboard-charts');
}
```

### 3.3 `m/app.html` — Add CSS and JS includes

```html
<link rel="stylesheet" href="../l8ui/sys/dashboard/l8dashboard.css">
<script src="../l8ui/sys/dashboard/l8dashboard.js"></script>
```

---

## Step 4: AI Writing Assistance (Go Backend)

Add a lightweight endpoint that AI coding agents (or the UI) can call to get writing suggestions for issue fields.

### 4.1 `go/bugs/triage/writer.go` (~120 lines)

AI writing assistant that generates or improves issue text:

```go
type WriteRequest struct {
    Action string `json:"action"` // "suggest_steps", "improve_description", "generate_acceptance_criteria", "summarize_comments"
    Input  string `json:"input"`
    Title  string `json:"title"`
}

type WriteResult struct {
    Output string `json:"output"`
}

func (t *Triager) AssistWriting(req *WriteRequest) (*WriteResult, error)
```

Actions:
- `suggest_steps` — Given a bug title/description, generate "Steps to Reproduce" structure
- `improve_description` — Suggest clearer wording for a bug/feature description
- `generate_acceptance_criteria` — From feature description, produce testable acceptance criteria
- `summarize_comments` — Condense comment bodies into key decisions

Uses the existing `Client.Complete()` with action-specific system prompts.

### 4.2 MCP Tool Addition

Add a 7th tool to the MCP server:

**`mcp/tools.go`** — Add `assist_writing` tool definition:
```javascript
{
    name: "assist_writing",
    description: "AI writing assistance for issue fields",
    inputSchema: {
        type: "object",
        properties: {
            action: { type: "string", enum: ["suggest_steps", "improve_description", "generate_acceptance_criteria", "summarize_comments"] },
            input: { type: "string", description: "The text to improve or use as context" },
            title: { type: "string", description: "Issue title for context" }
        },
        required: ["action", "input"]
    }
}
```

**`mcp/handlers.go`** — Add `handleAssistWriting` handler that delegates to `triage.Get().AssistWriting()`.

---

## Step 5: PRD Update

Mark Phase 5 as complete with details in `plans/l8bugs-prd.md`.

---

## File Summary

| Action | Count | Files |
|--------|-------|-------|
| Create (Dashboard CSS) | 1 | `l8ui/sys/dashboard/l8dashboard.css` |
| Create (Dashboard JS) | 1 | `l8ui/sys/dashboard/l8dashboard.js` |
| Create (Go Writer) | 1 | `triage/writer.go` |
| Modify (Desktop HTML) | 1 | `sections/system.html` |
| Modify (Desktop Config) | 1 | `l8ui/sys/l8sys-config.js` |
| Modify (Desktop Init) | 1 | `l8ui/sys/l8sys-init.js` |
| Modify (Desktop app.html) | 1 | `app.html` |
| Modify (Mobile HTML) | 1 | `m/sections/system.html` |
| Modify (Mobile app.html) | 1 | `m/app.html` |
| Modify (MCP tools) | 1 | `mcp/tools.go` |
| Modify (MCP handlers) | 1 | `mcp/handlers.go` |
| Modify (PRD) | 1 | `plans/l8bugs-prd.md` |
| **Total** | **12** | 3 new + 9 modified |

---

## Verification

```bash
# Go build
cd go && go build ./... && go vet ./...

# JS syntax
for f in go/bugs/website/web/l8ui/sys/dashboard/*.js; do node -c "$f"; done
for f in go/bugs/website/web/l8ui/sys/tracking/*.js; do node -c "$f"; done

# Dashboard JS exists
grep 'L8Dashboard' go/bugs/website/web/l8ui/sys/dashboard/l8dashboard.js

# Dashboard initialized
grep 'L8Dashboard' go/bugs/website/web/l8ui/sys/l8sys-init.js

# Dashboard CSS included
grep 'l8dashboard.css' go/bugs/website/web/app.html

# Dashboard JS included
grep 'l8dashboard.js' go/bugs/website/web/app.html

# Dashboard tab in HTML
grep 'data-module="dashboard"' go/bugs/website/web/sections/system.html

# Mobile parity
grep 'data-module="dashboard"' go/bugs/website/web/m/sections/system.html
grep 'l8dashboard' go/bugs/website/web/m/app.html

# Enhanced views
grep 'timeline\|calendar\|gantt' go/bugs/website/web/l8ui/sys/l8sys-config.js

# Writing assistant
grep 'AssistWriting' go/bugs/triage/writer.go

# MCP tool
grep 'assist_writing' go/bugs/mcp/tools.go
```
