package tests

import (
	"fmt"
	"github.com/saichler/l8bugs/go/tests/mocks"
	"github.com/saichler/l8types/go/ifs"
	"strings"
	"testing"
)

func testCRUD(t *testing.T, client *mocks.BugsClient) {
	testCRUDProject(t, client)
	testCRUDBug(t, client)
	testCRUDFeature(t, client)
	testCRUDSprint(t, client)
	testCRUDAssignee(t, client)
	testCRUDDigest(t, client)
}

func testCRUDProject(t *testing.T, client *mocks.BugsClient) {
	projectId := ifs.NewUuid()
	project := map[string]interface{}{
		"projectId":   projectId,
		"name":        "CRUD Test Project",
		"key":         "CRUD",
		"description": "Created by CRUD test",
		"status":      1,
		"visibility":  3,
	}
	_, err := client.Post("/bugs/20/Project", project)
	if err != nil {
		t.Fatalf("POST Project failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from BugsProject where projectId=%s", projectId))
	getResp, err := client.Get("/bugs/20/Project", q)
	if err != nil {
		t.Fatalf("GET Project failed: %v", err)
	}
	if !strings.Contains(getResp, "CRUD Test Project") {
		t.Fatalf("GET Project did not return expected name, got: %s", getResp)
	}

	project["description"] = "Updated by CRUD test"
	_, err = client.Put("/bugs/20/Project", project)
	if err != nil {
		t.Fatalf("PUT Project failed: %v", err)
	}

	delQ := mocks.L8QueryText(fmt.Sprintf("select * from BugsProject where projectId=%s", projectId))
	_, err = client.Delete("/bugs/20/Project", delQ)
	if err != nil {
		t.Fatalf("DELETE Project failed: %v", err)
	}
}

func testCRUDBug(t *testing.T, client *mocks.BugsClient) {
	bugId := ifs.NewUuid()
	bug := map[string]interface{}{
		"bugId":       bugId,
		"projectId":   testStore.ProjectIDs[0],
		"title":       "CRUD Test Bug",
		"description": "Created by CRUD test",
		"status":      1,
		"priority":    2,
		"severity":    2,
	}
	_, err := client.Post("/bugs/20/Bug", bug)
	if err != nil {
		t.Fatalf("POST Bug failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from Bug where bugId=%s", bugId))
	getResp, err := client.Get("/bugs/20/Bug", q)
	if err != nil {
		t.Fatalf("GET Bug failed: %v", err)
	}
	if !strings.Contains(getResp, "CRUD Test Bug") {
		t.Fatalf("GET Bug did not return expected title, got: %s", getResp)
	}

	bug["description"] = "Updated by CRUD test"
	_, err = client.Put("/bugs/20/Bug", bug)
	if err != nil {
		t.Fatalf("PUT Bug failed: %v", err)
	}

	delQ := mocks.L8QueryText(fmt.Sprintf("select * from Bug where bugId=%s", bugId))
	_, err = client.Delete("/bugs/20/Bug", delQ)
	if err != nil {
		t.Fatalf("DELETE Bug failed: %v", err)
	}
}

func testCRUDFeature(t *testing.T, client *mocks.BugsClient) {
	featureId := ifs.NewUuid()
	feature := map[string]interface{}{
		"featureId":   featureId,
		"projectId":   testStore.ProjectIDs[0],
		"title":       "CRUD Test Feature",
		"description": "Created by CRUD test",
		"status":      1,
		"priority":    3,
	}
	_, err := client.Post("/bugs/20/Feature", feature)
	if err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from Feature where featureId=%s", featureId))
	getResp, err := client.Get("/bugs/20/Feature", q)
	if err != nil {
		t.Fatalf("GET Feature failed: %v", err)
	}
	if !strings.Contains(getResp, "CRUD Test Feature") {
		t.Fatalf("GET Feature did not return expected title, got: %s", getResp)
	}

	feature["description"] = "Updated by CRUD test"
	_, err = client.Put("/bugs/20/Feature", feature)
	if err != nil {
		t.Fatalf("PUT Feature failed: %v", err)
	}

	delQ := mocks.L8QueryText(fmt.Sprintf("select * from Feature where featureId=%s", featureId))
	_, err = client.Delete("/bugs/20/Feature", delQ)
	if err != nil {
		t.Fatalf("DELETE Feature failed: %v", err)
	}
}

func testCRUDSprint(t *testing.T, client *mocks.BugsClient) {
	sprintId := ifs.NewUuid()
	sprint := map[string]interface{}{
		"sprintId":  sprintId,
		"projectId": testStore.ProjectIDs[0],
		"name":      "CRUD Test Sprint",
		"goal":      "Test sprint lifecycle",
		"status":    1,
		"capacity":  50,
	}
	_, err := client.Post("/bugs/20/Sprint", sprint)
	if err != nil {
		t.Fatalf("POST Sprint failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from BugsSprint where sprintId=%s", sprintId))
	getResp, err := client.Get("/bugs/20/Sprint", q)
	if err != nil {
		t.Fatalf("GET Sprint failed: %v", err)
	}
	if !strings.Contains(getResp, "CRUD Test Sprint") {
		t.Fatalf("GET Sprint did not return expected name, got: %s", getResp)
	}

	sprint["goal"] = "Updated by CRUD test"
	_, err = client.Put("/bugs/20/Sprint", sprint)
	if err != nil {
		t.Fatalf("PUT Sprint failed: %v", err)
	}

	delQ := mocks.L8QueryText(fmt.Sprintf("select * from BugsSprint where sprintId=%s", sprintId))
	_, err = client.Delete("/bugs/20/Sprint", delQ)
	if err != nil {
		t.Fatalf("DELETE Sprint failed: %v", err)
	}
}

func testCRUDAssignee(t *testing.T, client *mocks.BugsClient) {
	assigneeId := ifs.NewUuid()
	assignee := map[string]interface{}{
		"assigneeId":   assigneeId,
		"name":         "CRUD Test User",
		"email":        "crudtest@example.com",
		"assigneeType": 1,
		"projectId":    testStore.ProjectIDs[0],
		"active":       true,
	}
	_, err := client.Post("/bugs/20/Assignee", assignee)
	if err != nil {
		t.Fatalf("POST Assignee failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from BugsAssignee where assigneeId=%s", assigneeId))
	getResp, err := client.Get("/bugs/20/Assignee", q)
	if err != nil {
		t.Fatalf("GET Assignee failed: %v", err)
	}
	if !strings.Contains(getResp, "CRUD Test User") {
		t.Fatalf("GET Assignee did not return expected name, got: %s", getResp)
	}

	assignee["email"] = "updated@example.com"
	_, err = client.Put("/bugs/20/Assignee", assignee)
	if err != nil {
		t.Fatalf("PUT Assignee failed: %v", err)
	}

	delQ := mocks.L8QueryText(fmt.Sprintf("select * from BugsAssignee where assigneeId=%s", assigneeId))
	_, err = client.Delete("/bugs/20/Assignee", delQ)
	if err != nil {
		t.Fatalf("DELETE Assignee failed: %v", err)
	}
}

func testCRUDDigest(t *testing.T, client *mocks.BugsClient) {
	digestId := ifs.NewUuid()
	digest := map[string]interface{}{
		"digestId":  digestId,
		"projectId": testStore.ProjectIDs[0],
		"period":    1,
		"summary":   "CRUD Test Digest Summary",
		"startDate": 1700000000,
		"endDate":   1700086400,
	}
	_, err := client.Post("/bugs/20/Digest", digest)
	if err != nil {
		t.Fatalf("POST Digest failed: %v", err)
	}

	q := mocks.L8QueryText(fmt.Sprintf("select * from BugsDigest where digestId=%s", digestId))
	getResp, err := client.Get("/bugs/20/Digest", q)
	if err != nil {
		t.Fatalf("GET Digest failed: %v", err)
	}
	if !strings.Contains(getResp, "CRUD Test Digest Summary") {
		t.Fatalf("GET Digest did not return expected summary, got: %s", getResp)
	}

	digest["summary"] = "Updated by CRUD test"
	_, err = client.Put("/bugs/20/Digest", digest)
	if err != nil {
		t.Fatalf("PUT Digest failed: %v", err)
	}

	delQ := mocks.L8QueryText(fmt.Sprintf("select * from BugsDigest where digestId=%s", digestId))
	_, err = client.Delete("/bugs/20/Digest", delQ)
	if err != nil {
		t.Fatalf("DELETE Digest failed: %v", err)
	}
}
