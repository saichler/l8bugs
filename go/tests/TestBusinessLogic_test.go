package tests

import (
	"github.com/saichler/l8bugs/go/tests/mocks"
	"github.com/saichler/l8types/go/ifs"
	"testing"
)

func testBusinessLogic(t *testing.T, client *mocks.BugsClient) {
	testBugTransitions(t, client)
	testFeatureTransitions(t, client)
	testSprintLogic(t, client)
}

// putEntity PUTs an entity and returns any error.
func putEntity(client *mocks.BugsClient, endpoint string, entity map[string]interface{}) error {
	_, err := client.Put(endpoint, entity)
	return err
}

// advanceStatus advances an entity through the given status sequence via sequential PUTs.
func advanceStatus(t *testing.T, client *mocks.BugsClient, endpoint string, entity map[string]interface{}, statuses ...int) {
	for _, s := range statuses {
		entity["status"] = s
		if err := putEntity(client, endpoint, entity); err != nil {
			t.Fatalf("Advance to status %d on %s failed: %v", s, endpoint, err)
		}
	}
}

// newBug creates a Bug map with a fresh UUID and required fields.
func newBug(projectId string) map[string]interface{} {
	return map[string]interface{}{
		"bugId":     ifs.NewUuid(),
		"projectId": projectId,
		"title":     "Transition Test Bug",
		"priority":  2,
		"severity":  2,
	}
}

// newFeature creates a Feature map with a fresh UUID and required fields.
func newFeature(projectId string) map[string]interface{} {
	return map[string]interface{}{
		"featureId": ifs.NewUuid(),
		"projectId": projectId,
		"title":     "Transition Test Feature",
		"priority":  3,
	}
}

// newSprint creates a BugsSprint map with a fresh UUID and required fields.
func newSprint(projectId string) map[string]interface{} {
	return map[string]interface{}{
		"sprintId":  ifs.NewUuid(),
		"projectId": projectId,
		"name":      "Transition Test Sprint",
		"capacity":  30,
	}
}
