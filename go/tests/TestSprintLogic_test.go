package tests

import (
	"github.com/saichler/l8bugs/go/tests/mocks"
	"github.com/saichler/l8types/go/ifs"
	"testing"
)

const sprintEndpoint = "/bugs/20/Sprint"

func testSprintLogic(t *testing.T, client *mocks.BugsClient) {
	testSprintValidTransitions(t, client)
	testSprintInvalidTransitions(t, client)
	testSprintDateValidation(t, client)
}

// testSprintValidTransitions tests: Planning->Active->Completed
func testSprintValidTransitions(t *testing.T, client *mocks.BugsClient) {
	sprint := newSprint(testStore.ProjectIDs[0])
	if _, err := client.Post(sprintEndpoint, sprint); err != nil {
		t.Fatalf("POST Sprint failed: %v", err)
	}

	// Planning(1) -> Active(2) -> Completed(3)
	advanceStatus(t, client, sprintEndpoint, sprint, 2, 3)
}

// testSprintInvalidTransitions tests transitions that must be rejected
func testSprintInvalidTransitions(t *testing.T, client *mocks.BugsClient) {
	// Planning(1) -> Completed(3) — must go through Active
	s1 := newSprint(testStore.ProjectIDs[0])
	if _, err := client.Post(sprintEndpoint, s1); err != nil {
		t.Fatalf("POST Sprint failed: %v", err)
	}
	s1["status"] = 3
	if err := putEntity(client, sprintEndpoint, s1); err == nil {
		t.Fatal("Sprint Planning->Completed should have failed")
	}

	// Completed(3) -> Active(2) — terminal state
	s2 := newSprint(testStore.ProjectIDs[0])
	if _, err := client.Post(sprintEndpoint, s2); err != nil {
		t.Fatalf("POST Sprint failed: %v", err)
	}
	advanceStatus(t, client, sprintEndpoint, s2, 2, 3)
	s2["status"] = 2
	if err := putEntity(client, sprintEndpoint, s2); err == nil {
		t.Fatal("Sprint Completed->Active should have failed")
	}

	// Completed(3) -> Planning(1) — terminal state
	s3 := newSprint(testStore.ProjectIDs[0])
	if _, err := client.Post(sprintEndpoint, s3); err != nil {
		t.Fatalf("POST Sprint failed: %v", err)
	}
	advanceStatus(t, client, sprintEndpoint, s3, 2, 3)
	s3["status"] = 1
	if err := putEntity(client, sprintEndpoint, s3); err == nil {
		t.Fatal("Sprint Completed->Planning should have failed")
	}
}

// testSprintDateValidation tests that endDate must be after startDate
func testSprintDateValidation(t *testing.T, client *mocks.BugsClient) {
	// Invalid: endDate before startDate
	badSprint := map[string]interface{}{
		"sprintId":  ifs.NewUuid(),
		"projectId": testStore.ProjectIDs[0],
		"name":      "Bad Date Sprint",
		"capacity":  30,
		"startDate": 1700086400,
		"endDate":   1700000000,
	}
	if _, err := client.Post(sprintEndpoint, badSprint); err == nil {
		t.Fatal("POST Sprint with endDate before startDate should have failed")
	}

	// Valid: endDate after startDate
	goodSprint := map[string]interface{}{
		"sprintId":  ifs.NewUuid(),
		"projectId": testStore.ProjectIDs[0],
		"name":      "Good Date Sprint",
		"capacity":  30,
		"startDate": 1700000000,
		"endDate":   1700086400,
	}
	if _, err := client.Post(sprintEndpoint, goodSprint); err != nil {
		t.Fatalf("POST Sprint with valid dates failed: %v", err)
	}
}
