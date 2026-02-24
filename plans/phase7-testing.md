# Phase 7: Testing Implementation Plan

## Context
L8Bugs has 6 services (Project, Assignee, Bug, Feature, Sprint, Digest) with zero test coverage. This plan implements comprehensive tests following the l8erp test patterns: topology setup, mock data generation, service handler/getter verification, and CRUD lifecycle tests via HTTP.

**Credentials**: `operator/operator` (no access control yet)

## Files to Create

### Test Infrastructure (`go/tests/`)

#### 1. `TestInit.go` — Topology setup/teardown
- Import `l8test` framework (`t_topology`, `t_resources`)
- Global vars: `topo *TestTopology`, `testStore *MockDataStore`
- `setup()` → `NewTestTopology(4, []int{20000, 30000, 40000}, Info_Level)`
- `tear()` → `topo.Shutdown()`
- **Requires**: Add `github.com/saichler/l8test` dependency to `go.mod`

#### 2. `StartWebserver.go` — Web server for HTTP tests
- Follow `go/bugs/website/main1/main.go` pattern
- Port: 9443, Prefix: `common.PREFIX` ("/bugs/"), CertName: "/data/l8bugs"
- Register types via `website.RegisterTypes()`
- Register health service, activate webpoints service
- `go svr.Start()` (non-blocking)

#### 3. `TestAllService_test.go` — Main test orchestrator
- `TestMain(m)` → `setup()` / `m.Run()` / `tear()`
- `TestAllServices(t)`:
  1. `dropAllTables(t)` — clean slate via `common.OpenDBConection("operator","operator","operator")`
  2. Activate services: `services.ActivateBugsServices("operator", "operator", erpVnic)`
  3. Start web server on webVnic
  4. Create `mocks.NewBugsClient()`, authenticate with `operator/operator`
  5. `mocks.RunAllPhases(client, testStore)`
  6. Verify key counts (ProjectIDs, BugIDs not empty)
  7. Call `testServiceHandlers(t, erpVnic)`
  8. Call `testServiceGetters(t, erpVnic)`
  9. Call `testCRUD(t, client)`
  10. Call `testValidation(t, client)`

#### 4. `TestServiceHandlers_test.go` — Handler registration
```
testServiceHandlers(t, vnic):
  projects.Projects(vnic) → not nil
  assignees.Assignees(vnic) → not nil
  bugs.Bugs(vnic) → not nil
  features.Features(vnic) → not nil
  sprints.Sprints(vnic) → not nil
  digests.Digests(vnic) → not nil
```

#### 5. `TestServiceGetters_test.go` — Getter functions
```
testServiceGetters(t, vnic):
  projects.Project("test-id", vnic) → no error
  assignees.Assignee("test-id", vnic) → no error
  bugs.Bug("test-id", vnic) → no error
  features.Feature("test-id", vnic) → no error
  sprints.Sprint("test-id", vnic) → no error
  digests.Digest("test-id", vnic) → no error
```

#### 6. `TestCRUD_test.go` — End-to-end HTTP CRUD
For each service, test full lifecycle via HTTP client:
- **POST**: Create entity with required fields, verify 200/201
- **GET**: Fetch via L8Query (`select * from ModelName where id='x'`), verify fields match
- **PUT**: Update a field, verify persisted
- **DELETE**: Remove entity, verify gone
- Models to test: `BugsProject`, `Bug`, `Feature`, `BugsSprint`, `BugsAssignee`, `BugsDigest`
- L8Query model names must use protobuf type names (per global rule)

#### 7. `TestValidation_test.go` — Validation rules
- POST Bug without `Title` → error
- POST Bug without `ProjectId` → error
- POST Project without `Name` → error
- POST Project without `Key` → error
- POST Digest without `Summary` → error
- POST Digest without `ProjectId` → error
- POST Bug with empty body → ID still auto-generated (auto-gen validates)

### Mock Data (`go/tests/mocks/`)

#### 8. `client.go` — HTTP test client
- `BugsClient` struct with `baseURL`, `token`, `client`
- `NewBugsClient(baseURL, httpClient)` constructor
- `Authenticate(user, password)` → POST `/auth`
- `Post(endpoint, data)` → POST with Bearer token
- `Get(endpoint)` → GET with Bearer token
- `Put(endpoint, data)` → PUT with Bearer token
- `Delete(endpoint, query)` → DELETE with Bearer token

#### 9. `store.go` — MockDataStore
```go
type MockDataStore struct {
    ProjectIDs  []string
    AssigneeIDs []string
    BugIDs      []string
    FeatureIDs  []string
    SprintIDs   []string
    DigestIDs   []string
}
```

#### 10. `data.go` — Curated name arrays
- Project names (10): "Backend API", "Frontend App", "Mobile Client", etc.
- Project keys (10): "BAPI", "FAPP", "MOBC", etc.
- Bug titles (15): "Login page crashes on Safari", "API returns 500...", etc.
- Feature titles (10): "Dark mode support", "Export to CSV", etc.
- Component names (8): "Authentication", "API Gateway", "Database", etc.

#### 11. `gen_foundation.go` — Project & Assignee generators
- `generateProjects()` → 5 BugsProject with name, key, status, visibility
- `generateAssignees(store)` → 8 BugsAssignee linked to projects, mix of Human/AI Agent

#### 12. `gen_tracking.go` — Bug & Feature generators
- `generateBugs(store)` → 20 Bug with flavorable status distribution, linked to projects/assignees
- `generateFeatures(store)` → 10 Feature with status distribution, linked to projects/assignees

#### 13. `gen_planning.go` — Sprint & Digest generators
- `generateSprints(store)` → 6 BugsSprint linked to projects
- `generateDigests(store)` → 4 BugsDigest linked to projects

#### 14. `phases.go` — Phase orchestration + RunAllPhases + PrintSummary
- Phase 1: Foundation (projects, assignees) → POST `/bugs/20/Project`, `/bugs/20/Assignee`
- Phase 2: Tracking (bugs, features) → POST `/bugs/20/Bug`, `/bugs/20/Feature`
- Phase 3: Planning (sprints, digests) → POST `/bugs/20/Sprint`, `/bugs/20/Digest`

## Existing Code to Reuse
- `common.OpenDBConection()` — `go/bugs/common/defaults.go`
- `common.PREFIX` ("/bugs/") — `go/bugs/common/defaults.go`
- `services.ActivateBugsServices()` — `go/bugs/services/activate_bugs.go`
- `website.RegisterTypes()` — `go/bugs/website/shared.go`
- Service handler/getter functions from each `*Service.go`
- `l8test` framework: `NewTestTopology`, `VnicByVnetNum`, `Shutdown`
- `l8web/go/web/server` for RestServer
- Phase helper pattern from l8erp: `runOp()`, `extractIDs()`, `runPhase()`

## Dependencies to Add
- `github.com/saichler/l8test` — test topology framework
- Run `go mod tidy` after adding import

## Verification
```bash
# Build
cd go && go build ./tests/...

# Vet
cd go && go vet ./tests/...

# Run tests (requires PostgreSQL + L8 vnet running)
cd go && go test -v ./tests/ -count=1 -timeout 300s
```

## File Count: 14 files
- `go/tests/`: 7 files (TestInit, StartWebserver, 5 test files)
- `go/tests/mocks/`: 7 files (client, store, data, 3 generators, phases)
