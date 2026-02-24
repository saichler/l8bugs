# Phase 6: Digests, Webhooks & Effort Tracking

## Context

Phases 1-5 are complete: proto definitions, Go services, desktop+mobile UI, Sprint service, AI triage, MCP server, GitHub webhooks, AI root cause analysis, dashboard, enhanced views, and AI writing assistance. Phase 6 adds the remaining PRD features: AI-generated project digests (section 7.6), outbound webhook configuration (section 9.2), and effort estimation fields (section 7.1 item 7).

The proto changes are already in place (uncommitted in `proto/bugs.proto`):
- New Prime Object: `BugsDigest` + `BugsDigestList`
- New embedded child: `WebhookConfig` in `BugsProject.outbound_webhooks`
- New enums: `WebhookEventType`, `DigestPeriod`
- New fields on `Bug`: `actual_effort`, `ai_estimated_effort`, `ai_effort_confidence`
- New fields on `Feature`: `actual_effort`, `ai_estimated_effort`

---

## Step 1: Generate Proto Bindings

```bash
cd proto && ./make-bindings.sh
```

Verify generated types:
```bash
grep -A 15 "type BugsDigest struct" go/types/l8bugs/*.pb.go
grep "WebhookConfig" go/types/l8bugs/*.pb.go
grep "actual_effort\|ai_estimated_effort\|ai_effort_confidence" go/types/l8bugs/*.pb.go
```

---

## Step 2: Digest Go Service

### 2.1 `go/bugs/digests/DigestService.go` (~30 lines)

New service for BugsDigest. Follows exact pattern from `go/bugs/sprints/SprintService.go`.

```go
const (
    ServiceName = "Digest"
    ServiceArea = byte(20)
)
```

Primary key: `DigestId`. Activate via `common.ActivateService[*l8bugs.BugsDigest, *l8bugs.BugsDigestList]()`.

### 2.2 `go/bugs/digests/DigestServiceCallback.go` (~50 lines)

Follows pattern from `go/bugs/sprints/SprintServiceCallback.go`:
- `Before()`: auto-generate ID on POST (`common.GenerateID(&entity.DigestId)`)
- `validate()`: `ValidateRequired` for `projectId`, `period`, `summary`

### 2.3 `go/bugs/services/activate_bugs.go` — Add Digest activation

Add `digests.Activate(serviceConfig)` call after Sprint activation.

### 2.4 `go/bugs/website/shared.go` — Register Digest type

Add `registerDigestType()` following existing `registerSprintType()` pattern:
- `introspect.AddPrimaryKeyDecorator` for `digestId`
- `registry.Register(&l8bugs.BugsDigest{})`

---

## Step 3: Digest UI (Desktop)

### 3.1 `go/bugs/website/web/l8ui/sys/l8sys-config.js` — Add Digest service

Add to the tracking module's services array:
```javascript
{ key: 'digests', label: 'Digests', model: 'BugsDigest', endpoint: 'Digest' }
```

### 3.2 `go/bugs/website/web/l8ui/sys/tracking/l8tracking-enums.js` — Add Digest enums

Add `DIGEST_PERIOD` and `WEBHOOK_EVENT_TYPE` enum definitions + renderers. Export in `L8Tracking.enums` and `L8Tracking.render`.

### 3.3 `go/bugs/website/web/l8ui/sys/tracking/l8tracking-columns.js` — Add Digest columns

Add `BugsDigest` column definitions:
- `digestId` (hidden), `projectId` (reference), `period` (enum), `startDate` (date), `endDate` (date), `summary` (text), `generatedDate` (date)

Also add effort columns to existing Bug and Feature column definitions:
- `actualEffort`, `aiEstimatedEffort`, `aiEffortConfidence` (Bug only)

### 3.4 `go/bugs/website/web/l8ui/sys/tracking/l8tracking-forms.js` — Add Digest form

Add `BugsDigest` form definition:
- Section "Digest Details": `projectId` (reference), `period` (select), `startDate` (date), `endDate` (date), `generatedDate` (date, readonly)
- Section "Content": `summary` (textarea), `keyMetrics` (textarea), `blockers` (textarea), `actionItems` (textarea)

Also add effort fields to existing Bug and Feature forms:
- Section "Effort" (new section): `estimatedEffort` (existing), `actualEffort`, `aiEstimatedEffort`, `aiEffortConfidence` (Bug only)

Add `BugsDigest` to `L8Tracking.primaryKeys`: `BugsDigest: 'digestId'`

### 3.5 `go/bugs/website/web/l8ui/sys/tracking/l8tracking-reference.js` — Add Digest reference

Add `BugsDigest` reference entry for any forms that might reference digests.

### 3.6 `go/bugs/website/web/js/reference-registry-bugs.js` — Add Digest reference

Add `BugsDigest` entry to the desktop reference registry.

### 3.7 WebhookConfig inline table in BugsProject form

Add a new section "Outbound Webhooks" to the `BugsProject` form with `f.inlineTable('outboundWebhooks', ...)` containing:
- `webhookConfigId` (hidden), `url` (text, required), `secret` (text), `description` (text), `active` (checkbox), `events` (multi-select/text for now), `createdDate` (date, readonly)

---

## Step 4: Digest UI (Mobile)

### 4.1 `go/bugs/website/web/m/js/tracking/l8tracking-enums.js` — Add Digest enums (mobile)

Add `DIGEST_PERIOD` and `WEBHOOK_EVENT_TYPE` with mobile renderer pattern (`Layer8MRenderers`).

### 4.2 `go/bugs/website/web/m/js/tracking/l8tracking-columns.js` — Add Digest columns (mobile)

Add `BugsDigest` columns + effort columns on Bug/Feature.

### 4.3 `go/bugs/website/web/m/js/tracking/l8tracking-forms.js` — Add Digest form (mobile)

Add `BugsDigest` form + effort fields on Bug/Feature + webhook inline table on BugsProject.

### 4.4 `go/bugs/website/web/m/js/mobile-config-bugs.js` — Add Digest reference (mobile)

Register BugsDigest in mobile reference config.

### 4.5 `go/bugs/website/web/m/js/layer8m-nav-config-bugs.js` — Add Digest to nav

Add Digest service to the tracking module's services list in mobile nav config.

---

## Step 5: MCP Tools for Digests

### 5.1 `go/bugs/mcp/tools.go` — Add `generate_digest` tool

New MCP tool definition:
```json
{
    "name": "generate_digest",
    "description": "Generate an AI project digest/summary for a time period",
    "inputSchema": {
        "properties": {
            "project_id": { "type": "string" },
            "period": { "type": "string", "enum": ["daily", "weekly", "custom"] },
            "start_date": { "type": "string", "description": "ISO date (for custom period)" },
            "end_date": { "type": "string", "description": "ISO date (for custom period)" }
        },
        "required": ["project_id", "period"]
    }
}
```

### 5.2 `go/bugs/mcp/handlers.go` — Add `handleGenerateDigest`

Handler that:
1. Fetches bugs and features for the project within the date range
2. Calls triage AI to generate a summary, key metrics, blockers, and action items
3. POSTs the BugsDigest to the Digest service
4. Returns the digest content

### 5.3 `go/bugs/triage/digest.go` (~120 lines)

AI digest generator:
- `func (t *Triager) GenerateDigest(bugs []*l8bugs.Bug, features []*l8bugs.Feature, period string) (*DigestResult, error)`
- Uses `Client.Complete()` with a system prompt that summarizes activity, extracts metrics, identifies blockers, and suggests action items

---

## Step 6: AI Effort Estimation

### 6.1 `go/bugs/triage/effort.go` (~100 lines)

AI effort estimator:
- `func (t *Triager) EstimateEffort(title, description, component string) (*EffortResult, error)`
- Returns estimated effort (story points) and confidence percentage
- Called during triage (existing `TriageBug`/`TriageFeature` functions)

### 6.2 Update `go/bugs/triage/triage.go`

In `TriageBug()` and `TriageFeature()`, after existing triage logic, call `EstimateEffort()` and set `ai_estimated_effort` and `ai_effort_confidence` on the entity.

---

## File Summary

| Action | Count | Files |
|--------|-------|-------|
| Create (Go Digest service) | 2 | `digests/DigestService.go`, `digests/DigestServiceCallback.go` |
| Create (Go AI digest) | 1 | `triage/digest.go` |
| Create (Go AI effort) | 1 | `triage/effort.go` |
| Modify (Go activate) | 1 | `services/activate_bugs.go` |
| Modify (Go type registry) | 1 | `website/shared.go` |
| Modify (Desktop config) | 1 | `l8ui/sys/l8sys-config.js` |
| Modify (Desktop enums) | 1 | `l8ui/sys/tracking/l8tracking-enums.js` |
| Modify (Desktop columns) | 1 | `l8ui/sys/tracking/l8tracking-columns.js` |
| Modify (Desktop forms) | 1 | `l8ui/sys/tracking/l8tracking-forms.js` |
| Modify (Desktop reference) | 2 | `l8tracking-reference.js`, `reference-registry-bugs.js` |
| Modify (Mobile enums) | 1 | `m/js/tracking/l8tracking-enums.js` |
| Modify (Mobile columns) | 1 | `m/js/tracking/l8tracking-columns.js` |
| Modify (Mobile forms) | 1 | `m/js/tracking/l8tracking-forms.js` |
| Modify (Mobile config) | 1 | `m/js/mobile-config-bugs.js` |
| Modify (Mobile nav) | 1 | `m/js/layer8m-nav-config-bugs.js` |
| Modify (MCP tools) | 1 | `mcp/tools.go` |
| Modify (MCP handlers) | 1 | `mcp/handlers.go` |
| Modify (Triage) | 1 | `triage/triage.go` |
| **Total** | **20** | 4 new + 16 modified |

---

## Verification

```bash
# Generate proto bindings
cd proto && ./make-bindings.sh

# Go build + vet
cd go && go build ./... && go vet ./...

# JS syntax check
for f in go/bugs/website/web/l8ui/sys/tracking/*.js; do node -c "$f"; done
for f in go/bugs/website/web/m/js/tracking/*.js; do node -c "$f"; done

# Digest service exists
grep "ServiceName" go/bugs/digests/DigestService.go

# Digest type registered
grep "BugsDigest" go/bugs/website/shared.go

# Digest in UI config
grep "digests" go/bugs/website/web/l8ui/sys/l8sys-config.js

# Effort fields in forms
grep "actualEffort\|aiEstimatedEffort" go/bugs/website/web/l8ui/sys/tracking/l8tracking-forms.js

# Webhook inline table in project form
grep "outboundWebhooks" go/bugs/website/web/l8ui/sys/tracking/l8tracking-forms.js

# MCP digest tool
grep "generate_digest" go/bugs/mcp/tools.go

# Mobile parity
grep "digests" go/bugs/website/web/m/js/layer8m-nav-config-bugs.js
grep "BugsDigest" go/bugs/website/web/m/js/tracking/l8tracking-columns.js
```
