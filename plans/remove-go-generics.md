# Plan: Remove Go Generics from l8bugs (Using l8common)

## Context
The l8bugs project uses Go generics extensively in its service framework (`go/bugs/common/`). Per the `no-go-generics` rule, all Layer 8 projects must use `interface{}` instead. The `l8common` project (`../l8common`) already provides non-generic equivalents of all these abstractions. This plan replaces l8bugs's generic `common` package with imports from `l8common`.

## Approach
**Import l8common instead of rewriting.** Before migrating l8bugs, enhance l8common's validation builder so that `Require`, `RequireInt64`, `Enum`, `DateNotZero`, and `DateAfter` accept typed getter functions via reflection — matching the pattern `Custom`, `After`, and `BeforeAction` already use. This eliminates repetitive `e.(*Type)` boilerplate in every getter closure.

### Critical: NewValidation requires registered types
l8common's `NewValidation(&l8bugs.Bug{}, vnic)` auto-derives `setID` by looking up the PK decorator in the introspector. This requires `RegisterType` to have been called on the same vnic BEFORE the callback is constructed. Currently l8bugs only calls `RegisterTypes` in `website/shared.go` (web vnic), not the main server vnic.

**Fix:** Call `website.RegisterTypes(nic.Resources())` at the top of `ActivateBugsServices()` before any `Activate()` calls.

## Scope
- **Enhance 1 file** in l8common (validation_builder.go — typed getter support)
- **Delete 6 files** in `go/bugs/common/` (replaced by l8common)
- **Slim down 1 file** (`defaults.go` — keep only project-specific constants)
- **Modify 6 service files** (ActivateService + GetEntity calls)
- **Modify 6 callback files** (NewValidation + StatusTransitionConfig calls)
- **Modify 1 activation orchestrator** (add RegisterTypes call)
- **Modify 1 website file** (RegisterType calls)
- **Modify 1 webhook file** (QueryEntities -> GetEntitiesByQuery)
- **Modify 1 test mock file** (extractIDs helper)

---

## Traceability Matrix

| # | Source | Action Item | Phase |
|---|--------|-------------|-------|
| 1 | validation_builder.go | Enhance `Require` to accept typed getters via reflection | Phase 0 |
| 2 | validation_builder.go | Enhance `RequireInt64` to accept typed getters via reflection | Phase 0 |
| 3 | validation_builder.go | Enhance `Enum` to accept typed getters via reflection | Phase 0 |
| 4 | validation_builder.go | Enhance `DateNotZero` to accept typed getters via reflection | Phase 0 |
| 5 | validation_builder.go | Enhance `DateAfter` to accept typed getters via reflection | Phase 0 |
| 6 | go.mod/go.sum/vendor | Add l8common dependency, re-vendor | Phase 1 |
| 7 | common/service_factory.go | Delete — replaced by l8common | Phase 2 |
| 8 | common/service_callback.go | Delete — replaced by l8common | Phase 2 |
| 9 | common/validation_builder.go | Delete — replaced by l8common | Phase 2 |
| 10 | common/status_machine.go | Delete — replaced by l8common | Phase 2 |
| 11 | common/validation_static.go | Delete — replaced by l8common | Phase 2 |
| 12 | common/type_registry.go | Delete — replaced by l8common | Phase 2 |
| 13 | common/defaults.go | Remove delegatable functions, keep constants | Phase 2 |
| 14 | services/activate_bugs.go | Add `RegisterTypes(nic.Resources())` before activations | Phase 3 |
| 15 | bugs/BugService.go | Remove type params, add proto instances, add type assertion | Phase 4 |
| 16 | features/FeatureService.go | Remove type params, add proto instances, add type assertion | Phase 4 |
| 17 | assignees/AssigneeService.go | Remove type params, add proto instances, add type assertion | Phase 4 |
| 18 | projects/ProjectService.go | Remove type params, add proto instances, add type assertion | Phase 4 |
| 19 | sprints/SprintService.go | Remove type params, add proto instances, add type assertion | Phase 4 |
| 20 | digests/DigestService.go | Remove type params, add proto instances, add type assertion | Phase 4 |
| 21 | bugs/BugServiceCallback.go | Switch to l8common NewValidation + typed getters | Phase 5 |
| 22 | features/FeatureServiceCallback.go | Switch to l8common NewValidation + typed getters | Phase 5 |
| 23 | assignees/AssigneeServiceCallback.go | Switch to l8common NewValidation + typed getters | Phase 5 |
| 24 | projects/ProjectServiceCallback.go | Switch to l8common NewValidation + typed getters | Phase 5 |
| 25 | sprints/SprintServiceCallback.go | Switch to l8common NewValidation + typed getters | Phase 5 |
| 26 | digests/DigestServiceCallback.go | Switch to l8common NewValidation + typed getters | Phase 5 |
| 27 | website/shared.go | Remove type params from RegisterType calls | Phase 6 |
| 28 | webhook/webhook.go | Add type assertions, switch QueryEntities | Phase 7 |
| 29 | tests/mocks/phases.go | Delete extractIDs generic, inline at 6 call sites | Phase 8 |
| 30 | All modified files | Verify file size stays under 500 lines | All phases |

---

## Phase 0: Enhance l8common Validation Builder

**File:** `../l8common/go/common/validation_builder.go`

Make `Require`, `RequireInt64`, `Enum`, `DateNotZero`, and `DateAfter` accept typed getter functions via reflection, matching the pattern `Custom`/`After`/`BeforeAction` already use.

**Current** (requires `func(interface{}) string`):
```go
func (b *VB) Require(getter func(interface{}) string, name string) *VB
```

**New** (accepts both `func(interface{}) string` AND `func(*ConcreteType) string`):
```go
func (b *VB) Require(getter interface{}, name string) *VB {
    if typed, ok := getter.(func(interface{}) string); ok {
        b.validators = append(b.validators, func(e interface{}, _ ifs.IVNic) error {
            return ValidateRequired(typed(e), name)
        })
        return b
    }
    fnVal := reflect.ValueOf(getter)
    b.validators = append(b.validators, func(e interface{}, _ ifs.IVNic) error {
        results := fnVal.Call([]reflect.Value{reflect.ValueOf(e)})
        return ValidateRequired(results[0].String(), name)
    })
    return b
}
```

Apply the same pattern to:
- **`RequireInt64`**: `getter interface{}` — `results[0].Int()`
- **`Enum`**: `getter interface{}` — `int32(results[0].Int())`
- **`DateNotZero`**: `getter interface{}` — `results[0].Int()`
- **`DateAfter`**: `getter1, getter2 interface{}` — both use `results[0].Int()`

This allows callers to write:
```go
Require(func(e *l8bugs.Bug) string { return e.BugId }, "BugId")
```
instead of:
```go
Require(func(e interface{}) string { return e.(*l8bugs.Bug).BugId }, "BugId")
```

## Phase 1: Add l8common Dependency & Re-vendor

```bash
cd go && rm -rf go.sum go.mod vendor && go mod init && GOPROXY=direct GOPRIVATE=github.com go mod tidy && go mod vendor
```

## Phase 2: Delete Replaced Files & Slim Down defaults.go

**Delete** (entirely replaced by l8common):
- `go/bugs/common/service_factory.go`
- `go/bugs/common/service_callback.go`
- `go/bugs/common/validation_builder.go`
- `go/bugs/common/status_machine.go`
- `go/bugs/common/validation_static.go`
- `go/bugs/common/type_registry.go`

**Slim down** `go/bugs/common/defaults.go`:
- Remove `CreateResources` body — delegate to `l8common.CreateResources(alias, "/data/logs/l8bugs", uint32(BUGS_VNET))`
- Remove `OpenDBConection` — use `l8common.OpenDBConection` at call sites
- Remove `WaitForSignal` — use `l8common.WaitForSignal` at call sites
- Keep only: `BUGS_VNET`, `BUGS_LOGS_VNET`, `PREFIX`, `DB_CREDS`, `DB_NAME` constants

## Phase 3: Update Service Activation Orchestrator

**File:** `go/bugs/services/activate_bugs.go`

Add `RegisterTypes` call before any service activation:
```go
func ActivateBugsServices(creds, dbname string, nic ifs.IVNic) {
    website.RegisterTypes(nic.Resources())  // PK decorators needed by NewValidation
    projects.Activate(creds, dbname, nic)
    // ... rest unchanged
}
```

## Phase 4: Update Service Files (6 files)

**Files:** `go/bugs/{bugs,features,assignees,projects,sprints,digests}/*Service.go`

1. Import `l8common "github.com/saichler/l8common/go/common"` instead of local `common`
2. `ActivateService`: Remove type params, pass proto instances, pass `vnic` to callback constructor:
   ```go
   l8common.ActivateService(l8common.ServiceConfig{
       ServiceName: ServiceName, ServiceArea: ServiceArea,
       PrimaryKey: "BugId", Callback: newBugServiceCallback(vnic),
   }, &l8bugs.Bug{}, &l8bugs.BugList{}, creds, dbname, vnic)
   ```
3. Typed getter wrappers: Add type assertion on `interface{}` result:
   ```go
   func Bug(bugId string, vnic ifs.IVNic) (*l8bugs.Bug, error) {
       result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.Bug{BugId: bugId}, vnic)
       if err != nil { return nil, err }
       if result == nil { return nil, nil }
       return result.(*l8bugs.Bug), nil
   }
   ```

## Phase 5: Update Callback Files (6 files)

**Files:** `go/bugs/{bugs,features,assignees,projects,sprints,digests}/*ServiceCallback.go`

Thanks to Phase 0, getter closures stay typed — just remove `[T]` type params:

1. Import `l8common "github.com/saichler/l8common/go/common"`
2. Callback constructors accept `vnic ifs.IVNic` parameter
3. `NewValidation[l8bugs.XXX]("TypeName", setIDFunc)` -> `l8common.NewValidation(&l8bugs.XXX{}, vnic)` (setID auto-derived)
4. **Getters stay typed** (no `interface{}` boilerplate needed):
   ```go
   func newProjectServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
       return l8common.NewValidation(&l8bugs.BugsProject{}, vnic).
           Require(func(e *l8bugs.BugsProject) string { return e.ProjectId }, "ProjectId").
           Require(func(e *l8bugs.BugsProject) string { return e.Name }, "Name").
           Require(func(e *l8bugs.BugsProject) string { return e.Key }, "Key").
           Build()
   }
   ```
5. `StatusTransitionConfig` (Bug, Feature, Sprint) — struct fields remain `func(interface{})` since these are fixed-signature struct fields, not builder methods:
   ```go
   StatusTransition(&l8common.StatusTransitionConfig{
       StatusGetter:  func(e interface{}) int32 { return int32(e.(*l8bugs.Bug).Status) },
       StatusSetter:  func(e interface{}, s int32) { e.(*l8bugs.Bug).Status = l8bugs.BugStatus(s) },
       FilterBuilder: func(e interface{}) interface{} { return &l8bugs.Bug{BugId: e.(*l8bugs.Bug).BugId} },
       // ... transitions/names unchanged
   })
   ```
6. `.After(...)` / `.Custom(...)` calls stay typed (l8common already supports this):
   ```go
   After(func(entity *l8bugs.Bug, action ifs.Action, _ ifs.IVNic) error { ... })
   ```
7. `.DateAfter(...)` getters stay typed:
   ```go
   DateAfter(
       func(e *l8bugs.BugsSprint) int64 { return e.EndDate },
       func(e *l8bugs.BugsSprint) int64 { return e.StartDate },
       "EndDate", "StartDate")
   ```

## Phase 6: Update Website Registration (`go/bugs/website/shared.go`)

Import `l8common` and change all `RegisterType` calls:
```go
l8common.RegisterType(resources, &l8bugs.Bug{}, &l8bugs.BugList{}, "BugId")
```

## Phase 7: Update Webhook (`go/bugs/webhook/webhook.go`)

1. Import `l8common`
2. `GetEntity` results: Add type assertions for `bug` and `feature` variables
3. `QueryEntities[l8bugs.BugsProject]` -> `l8common.GetEntitiesByQuery`:
   ```go
   results, err := l8common.GetEntitiesByQuery(projectService, serviceArea, query, h.vnic)
   if err != nil || len(results) == 0 { return nil }
   return results[0].(*l8bugs.BugsProject)
   ```

## Phase 8: Update Test Mocks (`go/tests/mocks/phases.go`)

Delete generic `extractIDs[T any]`. Inline the ID extraction at each of the 6 call sites:
```go
ids := make([]string, len(projects))
for i, p := range projects { ids[i] = p.ProjectId }
```

---

## Phase 9: End-to-End Verification

### 9A. Compilation & Static Analysis
```bash
# Verify l8common changes first
cd ../l8common/go && go build ./... && go vet ./...

# Then verify l8bugs
cd /home/saichler/proj/src/github.com/saichler/l8bugs/go
go build ./...
go vet ./...
```

### 9B. No Generics Remain
```bash
grep -rn '\[.*any\]\|ProtoMessage\[' --include='*.go' go/bugs/ go/tests/ | grep -v vendor/
# Should return zero matches
```

### 9C. File Size Check
Verify no modified file exceeds 500 lines (per maintainability.md):
```bash
wc -l go/bugs/bugs/BugServiceCallback.go go/bugs/features/FeatureServiceCallback.go \
      go/bugs/sprints/SprintServiceCallback.go go/bugs/webhook/webhook.go \
      go/tests/mocks/phases.go go/bugs/common/defaults.go
# All files must be under 500 lines
```

### 9D. Functional Smoke Test Checklist
- [ ] All 6 services activate without errors (check server logs)
- [ ] Bug CRUD: POST creates with auto-generated ID, GET retrieves, PUT updates, DELETE removes
- [ ] Feature CRUD: same as Bug
- [ ] Project CRUD: same as Bug
- [ ] Sprint CRUD: same as Bug, DateAfter validation rejects EndDate < StartDate
- [ ] Assignee CRUD: same as Bug
- [ ] Digest CRUD: same as Bug
- [ ] Bug status transitions: only valid transitions succeed (Open→Triaged, etc.)
- [ ] Feature status transitions: only valid transitions succeed
- [ ] Sprint status transitions: only valid transitions succeed
- [ ] Webhook: GitHub push triggers bug/feature lookup and project query
- [ ] Validation: required field errors return correctly for each service
- [ ] Mock data upload: all 6 phases complete without errors
