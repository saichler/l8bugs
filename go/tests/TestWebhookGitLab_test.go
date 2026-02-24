package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/saichler/l8bugs/go/tests/mocks"
	"github.com/saichler/l8types/go/ifs"
)

const (
	gitlabWebhookPath   = "/bugs/webhook/gitlab"
	gitlabWebhookSecret = "gitlab-test-secret"
	gitlabWebhookRepo   = "https://gitlab.com/test/webhook-test"
)

func testGitLabWebhook(t *testing.T, client *mocks.BugsClient) {
	// Create a project with GitLab repo URL and webhook secret.
	projectId := ifs.NewUuid()
	project := map[string]interface{}{
		"projectId":     projectId,
		"name":          "GitLab Webhook Test Project",
		"key":           "GLB",
		"repositoryUrl": gitlabWebhookRepo,
		"webhookSecret": gitlabWebhookSecret,
	}
	if _, err := client.Post("/bugs/20/Project", project); err != nil {
		t.Fatalf("Failed to create GitLab webhook test project: %v", err)
	}

	testGitLabWebhookMethodNotAllowed(t, client)
	testGitLabWebhookInvalidToken(t, client)
	testGitLabWebhookPushLinksCommit(t, client, projectId)
	testGitLabWebhookMergedMRTransitionsBug(t, client, projectId)
	testGitLabWebhookNonMergedMRIgnored(t, client, projectId)
}

// testGitLabWebhookMethodNotAllowed verifies GET requests are rejected with 405.
func testGitLabWebhookMethodNotAllowed(t *testing.T, client *mocks.BugsClient) {
	url := client.BaseURL() + gitlabWebhookPath
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.HTTPClient().Do(req)
	if err != nil {
		t.Fatalf("GET GitLab webhook request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for GET on GitLab webhook, got %d", resp.StatusCode)
	}
}

// testGitLabWebhookInvalidToken verifies a wrong X-Gitlab-Token is rejected with 403.
func testGitLabWebhookInvalidToken(t *testing.T, client *mocks.BugsClient) {
	payload := gitlabPushPayload("abc12345", "fixes nothing")
	status := sendGitLabWebhook(t, client, "Push Hook", "wrong-token", payload)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for invalid GitLab token, got %d", status)
	}
}

// testGitLabWebhookPushLinksCommit verifies a GitLab push event links the commit to a bug.
func testGitLabWebhookPushLinksCommit(t *testing.T, client *mocks.BugsClient, projectId string) {
	bugId := ifs.NewUuid()
	bug := map[string]interface{}{
		"bugId":     bugId,
		"projectId": projectId,
		"title":     "GitLab Push Test Bug",
		"priority":  2,
		"severity":  2,
	}
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug for GitLab push test: %v", err)
	}

	commitSHA := "bbcc11dd22ee3344"
	payload := gitlabPushPayload(commitSHA, fmt.Sprintf("fixes %s", bugId))
	status := sendGitLabWebhook(t, client, "Push Hook", gitlabWebhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for GitLab push webhook, got %d", status)
	}

	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != commitSHA {
		t.Fatalf("expected linkedCommitSha=%s, got %v", commitSHA, entity["linkedCommitSha"])
	}
}

// testGitLabWebhookMergedMRTransitionsBug verifies a merged MR auto-transitions
// a bug from In Review to Resolved and links the MR URL.
func testGitLabWebhookMergedMRTransitionsBug(t *testing.T, client *mocks.BugsClient, projectId string) {
	bug := newBug(projectId)
	bugId := bug["bugId"].(string)
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug for GitLab MR test: %v", err)
	}
	advanceStatus(t, client, "/bugs/20/Bug", bug, 2, 3, 4) // Triaged, InProgress, InReview

	mergeCommit := "ddee11ff22334455"
	mrURL := gitlabWebhookRepo + "/-/merge_requests/1"
	payload := gitlabMRPayload(fmt.Sprintf("fixes %s", bugId), "MR description", mergeCommit, mrURL, "merge")
	status := sendGitLabWebhook(t, client, "Merge Request Hook", gitlabWebhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for GitLab MR webhook, got %d", status)
	}

	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != mergeCommit {
		t.Fatalf("expected linkedCommitSha=%s, got %v", mergeCommit, entity["linkedCommitSha"])
	}
	if entity["linkedPrUrl"] != mrURL {
		t.Fatalf("expected linkedPrUrl=%s, got %v", mrURL, entity["linkedPrUrl"])
	}
	assertStatus(t, entity["status"], 5, "BUG_STATUS_RESOLVED")
}

// testGitLabWebhookNonMergedMRIgnored verifies a closed MR without merge
// does not update the bug.
func testGitLabWebhookNonMergedMRIgnored(t *testing.T, client *mocks.BugsClient, projectId string) {
	bugId := ifs.NewUuid()
	bug := map[string]interface{}{
		"bugId":     bugId,
		"projectId": projectId,
		"title":     "GitLab Non-Merged MR Test Bug",
		"priority":  2,
		"severity":  2,
	}
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug for non-merged MR test: %v", err)
	}

	mrURL := gitlabWebhookRepo + "/-/merge_requests/2"
	payload := gitlabMRPayload(fmt.Sprintf("fixes %s", bugId), "Closed MR", "aabb00cc11dd2233", mrURL, "close")
	status := sendGitLabWebhook(t, client, "Merge Request Hook", gitlabWebhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for non-merged GitLab MR, got %d", status)
	}

	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != nil && entity["linkedCommitSha"] != "" {
		t.Fatalf("expected no linkedCommitSha for non-merged GitLab MR, got %v", entity["linkedCommitSha"])
	}
}

// --- GitLab helpers ---

// sendGitLabWebhook sends a GitLab webhook POST request with the appropriate headers.
func sendGitLabWebhook(t *testing.T, client *mocks.BugsClient, eventType, token string, payload []byte) int {
	t.Helper()
	url := client.BaseURL() + gitlabWebhookPath
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create GitLab webhook request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gitlab-Event", eventType)
	if token != "" {
		req.Header.Set("X-Gitlab-Token", token)
	}

	resp, err := client.HTTPClient().Do(req)
	if err != nil {
		t.Fatalf("GitLab webhook request failed: %v", err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

// gitlabPushPayload creates a GitLab push event JSON payload.
func gitlabPushPayload(commitID, message string) []byte {
	ev := map[string]interface{}{
		"ref": "refs/heads/main",
		"commits": []map[string]interface{}{
			{"id": commitID, "message": message},
		},
		"project": map[string]interface{}{
			"web_url":                gitlabWebhookRepo,
			"http_url":              gitlabWebhookRepo + ".git",
			"name":                  "webhook-test",
			"path_with_namespace":   "test/webhook-test",
		},
	}
	data, _ := json.Marshal(ev)
	return data
}

// gitlabMRPayload creates a GitLab merge request event JSON payload.
func gitlabMRPayload(title, description, mergeCommitSHA, mrURL, action string) []byte {
	ev := map[string]interface{}{
		"object_attributes": map[string]interface{}{
			"title":            title,
			"description":      description,
			"state":            "merged",
			"merge_commit_sha": mergeCommitSHA,
			"url":              mrURL,
			"action":           action,
		},
		"project": map[string]interface{}{
			"web_url":                gitlabWebhookRepo,
			"http_url":              gitlabWebhookRepo + ".git",
			"name":                  "webhook-test",
			"path_with_namespace":   "test/webhook-test",
		},
	}
	data, _ := json.Marshal(ev)
	return data
}
