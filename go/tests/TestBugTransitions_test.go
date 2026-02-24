package tests

import (
	"encoding/json"
	"fmt"
	"github.com/saichler/l8bugs/go/tests/mocks"
	"testing"
)

const bugEndpoint = "/bugs/20/Bug"

func testBugTransitions(t *testing.T, client *mocks.BugsClient) {
	testBugHappyPath(t, client)
	testBugReopenPath(t, client)
	testBugTerminalShortcuts(t, client)
	testBugInvalidTransitions(t, client)
	testBugAutoSetStatus(t, client)
}

// testBugHappyPath tests the full lifecycle: Open->Triaged->InProgress->InReview->Resolved->Closed
func testBugHappyPath(t *testing.T, client *mocks.BugsClient) {
	bug := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}

	// Open(1) -> Triaged(2) -> InProgress(3) -> InReview(4) -> Resolved(5) -> Closed(6)
	advanceStatus(t, client, bugEndpoint, bug, 2, 3, 4, 5, 6)
}

// testBugReopenPath tests: advance to Resolved, then Resolved->Reopened->Open
func testBugReopenPath(t *testing.T, client *mocks.BugsClient) {
	bug := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}

	// Advance to Resolved
	advanceStatus(t, client, bugEndpoint, bug, 2, 3, 4, 5)

	// Resolved(5) -> Reopened(7)
	bug["status"] = 7
	if err := putEntity(client, bugEndpoint, bug); err != nil {
		t.Fatalf("Bug Resolved->Reopened failed: %v", err)
	}

	// Reopened(7) -> Open(1)
	bug["status"] = 1
	if err := putEntity(client, bugEndpoint, bug); err != nil {
		t.Fatalf("Bug Reopened->Open failed: %v", err)
	}
}

// testBugTerminalShortcuts tests shortcuts from Open/Triaged to terminal states
func testBugTerminalShortcuts(t *testing.T, client *mocks.BugsClient) {
	// Open(1) -> Won't Fix(8)
	bug1 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug1); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	bug1["status"] = 8
	if err := putEntity(client, bugEndpoint, bug1); err != nil {
		t.Fatalf("Bug Open->WontFix failed: %v", err)
	}

	// Open(1) -> Duplicate(9)
	bug2 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug2); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	bug2["status"] = 9
	if err := putEntity(client, bugEndpoint, bug2); err != nil {
		t.Fatalf("Bug Open->Duplicate failed: %v", err)
	}

	// Triaged(2) -> Cannot Reproduce(10)
	bug3 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug3); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	advanceStatus(t, client, bugEndpoint, bug3, 2)
	bug3["status"] = 10
	if err := putEntity(client, bugEndpoint, bug3); err != nil {
		t.Fatalf("Bug Triaged->CannotReproduce failed: %v", err)
	}
}

// testBugInvalidTransitions tests transitions that must be rejected
func testBugInvalidTransitions(t *testing.T, client *mocks.BugsClient) {
	// Open(1) -> Closed(6) — skipping workflow
	bug1 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug1); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	bug1["status"] = 6
	if err := putEntity(client, bugEndpoint, bug1); err == nil {
		t.Fatal("Bug Open->Closed should have failed")
	}

	// Open(1) -> In Progress(3) — must go through Triaged
	bug2 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug2); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	bug2["status"] = 3
	if err := putEntity(client, bugEndpoint, bug2); err == nil {
		t.Fatal("Bug Open->InProgress should have failed")
	}

	// In Progress(3) -> Resolved(5) — must go through In Review
	bug3 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug3); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	advanceStatus(t, client, bugEndpoint, bug3, 2, 3)
	bug3["status"] = 5
	if err := putEntity(client, bugEndpoint, bug3); err == nil {
		t.Fatal("Bug InProgress->Resolved should have failed")
	}

	// Closed(6) -> Open(1) — terminal state
	bug4 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug4); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	advanceStatus(t, client, bugEndpoint, bug4, 2, 3, 4, 5, 6)
	bug4["status"] = 1
	if err := putEntity(client, bugEndpoint, bug4); err == nil {
		t.Fatal("Bug Closed->Open should have failed")
	}

	// Won't Fix(8) -> Open(1) — terminal state
	bug5 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug5); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	advanceStatus(t, client, bugEndpoint, bug5, 8)
	bug5["status"] = 1
	if err := putEntity(client, bugEndpoint, bug5); err == nil {
		t.Fatal("Bug WontFix->Open should have failed")
	}

	// Duplicate(9) -> Triaged(2) — terminal state
	bug6 := newBug(testStore.ProjectIDs[0])
	if _, err := client.Post(bugEndpoint, bug6); err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}
	advanceStatus(t, client, bugEndpoint, bug6, 9)
	bug6["status"] = 2
	if err := putEntity(client, bugEndpoint, bug6); err == nil {
		t.Fatal("Bug Duplicate->Triaged should have failed")
	}
}

// testBugAutoSetStatus verifies that POST without explicit status auto-sets to Open(1)
func testBugAutoSetStatus(t *testing.T, client *mocks.BugsClient) {
	bug := newBug(testStore.ProjectIDs[0])
	bugId := bug["bugId"].(string)

	if _, err := client.Post(bugEndpoint, bug); err != nil {
		t.Fatalf("POST Bug for auto-set test failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from Bug where bugId=%s", bugId))
	getResp, err := client.Get(bugEndpoint, q)
	if err != nil {
		t.Fatalf("GET Bug for auto-set test failed: %v", err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(getResp), &resp); err != nil {
		t.Fatalf("Failed to parse GET response: %v", err)
	}
	list, ok := resp["list"].([]interface{})
	if !ok || len(list) == 0 {
		t.Fatal("GET Bug returned empty list for auto-set test")
	}
	entity := list[0].(map[string]interface{})
	status := int(entity["status"].(float64))
	if status != 1 {
		t.Fatalf("Expected auto-set status=1 (Open), got %d", status)
	}
}
