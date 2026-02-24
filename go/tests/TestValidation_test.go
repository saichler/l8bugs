package tests

import (
	"github.com/saichler/l8bugs/go/tests/mocks"
	"strings"
	"testing"
)

func testValidation(t *testing.T, client *mocks.BugsClient) {
	testValidationBug(t, client)
	testValidationFeature(t, client)
	testValidationProject(t, client)
	testValidationSprint(t, client)
	testValidationAssignee(t, client)
	testValidationDigest(t, client)
	testValidationAutoID(t, client)
}

// testValidationBug verifies required field validation on Bug
func testValidationBug(t *testing.T, client *mocks.BugsClient) {
	// Missing title — should fail
	bugNoTitle := map[string]interface{}{
		"projectId": testStore.ProjectIDs[0],
	}
	_, err := client.Post("/bugs/20/Bug", bugNoTitle)
	if err == nil {
		t.Fatal("POST Bug without title should have failed")
	}

	// Missing projectId — should fail
	bugNoProject := map[string]interface{}{
		"title": "Bug without project",
	}
	_, err = client.Post("/bugs/20/Bug", bugNoProject)
	if err == nil {
		t.Fatal("POST Bug without projectId should have failed")
	}
}

// testValidationFeature verifies required field validation on Feature
func testValidationFeature(t *testing.T, client *mocks.BugsClient) {
	// Missing title — should fail
	featureNoTitle := map[string]interface{}{
		"projectId": testStore.ProjectIDs[0],
	}
	_, err := client.Post("/bugs/20/Feature", featureNoTitle)
	if err == nil {
		t.Fatal("POST Feature without title should have failed")
	}
	if !strings.Contains(err.Error(), "Title is required") {
		t.Fatalf("Expected 'Title is required' error, got: %v", err)
	}

	// Missing projectId — should fail
	featureNoProject := map[string]interface{}{
		"title": "Feature without project",
	}
	_, err = client.Post("/bugs/20/Feature", featureNoProject)
	if err == nil {
		t.Fatal("POST Feature without projectId should have failed")
	}
	if !strings.Contains(err.Error(), "ProjectId is required") {
		t.Fatalf("Expected 'ProjectId is required' error, got: %v", err)
	}
}

// testValidationProject verifies required field validation on Project
func testValidationProject(t *testing.T, client *mocks.BugsClient) {
	// Missing name — should fail
	projNoName := map[string]interface{}{
		"key": "TEST",
	}
	_, err := client.Post("/bugs/20/Project", projNoName)
	if err == nil {
		t.Fatal("POST Project without name should have failed")
	}

	// Missing key — should fail
	projNoKey := map[string]interface{}{
		"name": "Test Project No Key",
	}
	_, err = client.Post("/bugs/20/Project", projNoKey)
	if err == nil {
		t.Fatal("POST Project without key should have failed")
	}
}

// testValidationSprint verifies required field validation on Sprint
func testValidationSprint(t *testing.T, client *mocks.BugsClient) {
	// Missing name — should fail
	sprintNoName := map[string]interface{}{
		"projectId": testStore.ProjectIDs[0],
	}
	_, err := client.Post("/bugs/20/Sprint", sprintNoName)
	if err == nil {
		t.Fatal("POST Sprint without name should have failed")
	}
	if !strings.Contains(err.Error(), "Name is required") {
		t.Fatalf("Expected 'Name is required' error, got: %v", err)
	}

	// Missing projectId — should fail
	sprintNoProject := map[string]interface{}{
		"name": "Sprint without project",
	}
	_, err = client.Post("/bugs/20/Sprint", sprintNoProject)
	if err == nil {
		t.Fatal("POST Sprint without projectId should have failed")
	}
	if !strings.Contains(err.Error(), "ProjectId is required") {
		t.Fatalf("Expected 'ProjectId is required' error, got: %v", err)
	}
}

// testValidationAssignee verifies required field validation on Assignee
func testValidationAssignee(t *testing.T, client *mocks.BugsClient) {
	// Missing name — should fail
	assigneeNoName := map[string]interface{}{
		"email": "noname@example.com",
	}
	_, err := client.Post("/bugs/20/Assignee", assigneeNoName)
	if err == nil {
		t.Fatal("POST Assignee without name should have failed")
	}
	if !strings.Contains(err.Error(), "Name is required") {
		t.Fatalf("Expected 'Name is required' error, got: %v", err)
	}
}

// testValidationDigest verifies required field validation on Digest
func testValidationDigest(t *testing.T, client *mocks.BugsClient) {
	// Missing summary — should fail
	digestNoSummary := map[string]interface{}{
		"projectId": testStore.ProjectIDs[0],
	}
	_, err := client.Post("/bugs/20/Digest", digestNoSummary)
	if err == nil {
		t.Fatal("POST Digest without summary should have failed")
	}

	// Missing projectId — should fail
	digestNoProject := map[string]interface{}{
		"summary": "Digest without project",
	}
	_, err = client.Post("/bugs/20/Digest", digestNoProject)
	if err == nil {
		t.Fatal("POST Digest without projectId should have failed")
	}
}

// testValidationAutoID verifies that POST auto-generates IDs when not provided
func testValidationAutoID(t *testing.T, client *mocks.BugsClient) {
	// POST a project without explicit ID — should succeed (auto-generated)
	project := map[string]interface{}{
		"name":        "Auto ID Test",
		"key":         "AUTOID",
		"description": "Testing auto ID generation",
	}
	_, err := client.Post("/bugs/20/Project", project)
	if err != nil {
		t.Fatalf("POST Project for auto-ID test failed: %v", err)
	}

	// Verify the entity was created by querying its unique key
	q := mocks.L8QueryText("select * from BugsProject where key=AUTOID")
	getResp, err := client.Get("/bugs/20/Project", q)
	if err != nil {
		t.Fatalf("GET auto-ID project failed: %v", err)
	}
	if !strings.Contains(getResp, "Auto ID Test") {
		t.Fatalf("Auto-ID project not found in GET response: %s", getResp)
	}
}
