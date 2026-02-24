package tests

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/saichler/l8bugs/go/tests/mocks"
	"github.com/saichler/l8types/go/ifs"
)

const (
	webhookPath   = "/bugs/webhook/github"
	webhookSecret = "test-webhook-secret"
	webhookRepo   = "https://github.com/test/webhook-test"
)

func testWebhook(t *testing.T, client *mocks.BugsClient) {
	// Create a project with repositoryUrl and webhookSecret for all webhook tests.
	projectId := ifs.NewUuid()
	project := map[string]interface{}{
		"projectId":     projectId,
		"name":          "Webhook Test Project",
		"key":           "WHK",
		"repositoryUrl": webhookRepo,
		"webhookSecret": webhookSecret,
	}
	if _, err := client.Post("/bugs/20/Project", project); err != nil {
		t.Fatalf("Failed to create webhook test project: %v", err)
	}

	testWebhookMethodNotAllowed(t, client)
	testWebhookInvalidSignature(t, client)
	testWebhookPushLinksCommit(t, client, projectId)
	testWebhookMergedPRTransitionsBug(t, client, projectId)
	testWebhookMergedPRTransitionsFeature(t, client, projectId)
	testWebhookUnknownRepo(t, client)
	testWebhookNonMergedPRIgnored(t, client, projectId)
	testWebhookPushNoMatchingRef(t, client)
	testWebhookPushNoRefInMessage(t, client, projectId)

	// GitLab webhook tests
	testGitLabWebhook(t, client)
}

// testWebhookMethodNotAllowed verifies GET requests are rejected with 405.
func testWebhookMethodNotAllowed(t *testing.T, client *mocks.BugsClient) {
	url := client.BaseURL() + webhookPath
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.HTTPClient().Do(req)
	if err != nil {
		t.Fatalf("GET webhook request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for GET, got %d", resp.StatusCode)
	}
}

// testWebhookInvalidSignature verifies a bad HMAC signature is rejected with 403.
func testWebhookInvalidSignature(t *testing.T, client *mocks.BugsClient) {
	payload := pushPayload("abc12345", "fixes nothing")
	status := sendWebhookRaw(t, client, "push", "wrong-secret", payload)
	if status != http.StatusForbidden {
		t.Fatalf("expected 403 for invalid signature, got %d", status)
	}
}

// testWebhookPushLinksCommit verifies a push event links the commit SHA to a bug.
func testWebhookPushLinksCommit(t *testing.T, client *mocks.BugsClient, projectId string) {
	// Create a bug.
	bugId := ifs.NewUuid()
	bug := map[string]interface{}{
		"bugId":     bugId,
		"projectId": projectId,
		"title":     "Webhook Push Test Bug",
		"priority":  2,
		"severity":  2,
	}
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug: %v", err)
	}

	// Send a push event referencing the bug by ID.
	commitSHA := "aabbccdd11223344"
	payload := pushPayload(commitSHA, fmt.Sprintf("fixes %s", bugId))
	status := sendWebhookRaw(t, client, "push", webhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for push webhook, got %d", status)
	}

	// Verify the bug was updated.
	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != commitSHA {
		t.Fatalf("expected linkedCommitSha=%s, got %v", commitSHA, entity["linkedCommitSha"])
	}
}

// testWebhookMergedPRTransitionsBug verifies a merged PR auto-transitions
// a bug from In Review to Resolved and links the PR URL.
func testWebhookMergedPRTransitionsBug(t *testing.T, client *mocks.BugsClient, projectId string) {
	// Create a bug and advance to In Review (Open→Triaged→InProgress→InReview).
	bug := newBug(projectId)
	bugId := bug["bugId"].(string)
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug: %v", err)
	}
	advanceStatus(t, client, "/bugs/20/Bug", bug, 2, 3, 4) // Triaged, InProgress, InReview

	// Send a merged PR event referencing the bug.
	mergeCommit := "ff00ff0011223344"
	payload := prPayload(fmt.Sprintf("fixes %s", bugId), "Merged PR body", mergeCommit)
	status := sendWebhookRaw(t, client, "pull_request", webhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for PR webhook, got %d", status)
	}

	// Verify the bug was transitioned to Resolved and has the PR link.
	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != mergeCommit {
		t.Fatalf("expected linkedCommitSha=%s, got %v", mergeCommit, entity["linkedCommitSha"])
	}
	if entity["linkedPrUrl"] != webhookRepo+"/pull/1" {
		t.Fatalf("expected linkedPrUrl=%s/pull/1, got %v", webhookRepo, entity["linkedPrUrl"])
	}
	// Status 5 = BUG_STATUS_RESOLVED. protojson may return string or number.
	assertStatus(t, entity["status"], 5, "BUG_STATUS_RESOLVED")
}

// testWebhookMergedPRTransitionsFeature verifies a merged PR auto-transitions
// a feature from In Review to Done and links the PR URL.
func testWebhookMergedPRTransitionsFeature(t *testing.T, client *mocks.BugsClient, projectId string) {
	// Create a feature and advance to In Review
	// (Proposed→Triaged→Approved→InProgress→InReview).
	feature := newFeature(projectId)
	featureId := feature["featureId"].(string)
	if _, err := client.Post("/bugs/20/Feature", feature); err != nil {
		t.Fatalf("Failed to create feature: %v", err)
	}
	advanceStatus(t, client, "/bugs/20/Feature", feature, 2, 3, 4, 5) // Triaged, Approved, InProgress, InReview

	// Send a merged PR event referencing the feature.
	mergeCommit := "ee11ee2233445566"
	payload := prPayload(fmt.Sprintf("closes %s", featureId), "Feature PR", mergeCommit)
	status := sendWebhookRaw(t, client, "pull_request", webhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for PR webhook, got %d", status)
	}

	// Verify the feature was transitioned to Done and has the PR link.
	entity := getEntity(t, client, "/bugs/20/Feature", "Feature", "featureId", featureId)
	if entity["linkedCommitSha"] != mergeCommit {
		t.Fatalf("expected linkedCommitSha=%s, got %v", mergeCommit, entity["linkedCommitSha"])
	}
	if entity["linkedPrUrl"] != webhookRepo+"/pull/1" {
		t.Fatalf("expected linkedPrUrl=%s/pull/1, got %v", webhookRepo, entity["linkedPrUrl"])
	}
	// Status 6 = FEATURE_STATUS_DONE.
	assertStatus(t, entity["status"], 6, "FEATURE_STATUS_DONE")
}

// testWebhookUnknownRepo verifies that a webhook for an unknown repo
// returns 200 (silent ignore, no signature check since no project found).
func testWebhookUnknownRepo(t *testing.T, client *mocks.BugsClient) {
	payload := pushPayloadForRepo("abc12345", "fixes nothing", "https://github.com/unknown/repo")
	// No HMAC needed since the project won't be found (secretFunc returns "").
	status := sendWebhookRaw(t, client, "push", "", payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for unknown repo, got %d", status)
	}
}

// testWebhookNonMergedPRIgnored verifies that a closed but non-merged PR
// does not update the bug (no commit linking, no status transition).
func testWebhookNonMergedPRIgnored(t *testing.T, client *mocks.BugsClient, projectId string) {
	// Create a bug at Open status.
	bugId := ifs.NewUuid()
	bug := map[string]interface{}{
		"bugId":     bugId,
		"projectId": projectId,
		"title":     "Non-Merged PR Test Bug",
		"priority":  2,
		"severity":  2,
	}
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug: %v", err)
	}

	// Send a closed-but-not-merged PR event.
	payload := nonMergedPRPayload(fmt.Sprintf("fixes %s", bugId), "Closed without merge", "deadbeef12345678")
	status := sendWebhookRaw(t, client, "pull_request", webhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for non-merged PR, got %d", status)
	}

	// Verify the bug was NOT updated (no linkedCommitSha).
	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != nil && entity["linkedCommitSha"] != "" {
		t.Fatalf("expected no linkedCommitSha for non-merged PR, got %v", entity["linkedCommitSha"])
	}
}

// testWebhookPushNoMatchingRef verifies that a push referencing a nonexistent
// bug/feature ID still returns 200 (no error, just silently unmatched).
func testWebhookPushNoMatchingRef(t *testing.T, client *mocks.BugsClient) {
	fakeId := ifs.NewUuid()
	payload := pushPayload("1111222233334444", fmt.Sprintf("fixes %s", fakeId))
	status := sendWebhookRaw(t, client, "push", webhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for push with no matching ref, got %d", status)
	}
}

// testWebhookPushNoRefInMessage verifies that a push with no issue reference
// keywords in the commit message returns 200 and doesn't modify any entities.
func testWebhookPushNoRefInMessage(t *testing.T, client *mocks.BugsClient, projectId string) {
	// Create a bug to verify it stays unchanged.
	bugId := ifs.NewUuid()
	bug := map[string]interface{}{
		"bugId":     bugId,
		"projectId": projectId,
		"title":     "No-Ref Message Test Bug",
		"priority":  2,
		"severity":  2,
	}
	if _, err := client.Post("/bugs/20/Bug", bug); err != nil {
		t.Fatalf("Failed to create bug: %v", err)
	}

	// Push with a commit message that has no issue ref keywords.
	payload := pushPayload("5555666677778888", "refactored logging module")
	status := sendWebhookRaw(t, client, "push", webhookSecret, payload)
	if status != http.StatusOK {
		t.Fatalf("expected 200 for push with no refs, got %d", status)
	}

	// Bug should remain untouched.
	entity := getEntity(t, client, "/bugs/20/Bug", "Bug", "bugId", bugId)
	if entity["linkedCommitSha"] != nil && entity["linkedCommitSha"] != "" {
		t.Fatalf("expected no linkedCommitSha for no-ref push, got %v", entity["linkedCommitSha"])
	}
}

// --- Helpers ---

// sendWebhookRaw sends a webhook POST request and returns the HTTP status code.
func sendWebhookRaw(t *testing.T, client *mocks.BugsClient, eventType, secret string, payload []byte) int {
	t.Helper()
	url := client.BaseURL() + webhookPath
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create webhook request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", eventType)

	if secret != "" {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(payload)
		sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Hub-Signature-256", sig)
	}

	resp, err := client.HTTPClient().Do(req)
	if err != nil {
		t.Fatalf("Webhook request failed: %v", err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

// pushPayload creates a GitHub push event JSON payload for the test repo.
func pushPayload(commitID, message string) []byte {
	return pushPayloadForRepo(commitID, message, webhookRepo)
}

// pushPayloadForRepo creates a GitHub push event JSON payload for a given repo.
func pushPayloadForRepo(commitID, message, repo string) []byte {
	ev := map[string]interface{}{
		"ref": "refs/heads/main",
		"commits": []map[string]interface{}{
			{"id": commitID, "message": message},
		},
		"repository": map[string]interface{}{
			"html_url":  repo,
			"clone_url": repo + ".git",
			"full_name": "test/webhook-test",
		},
	}
	data, _ := json.Marshal(ev)
	return data
}

// nonMergedPRPayload creates a GitHub closed-but-not-merged PR event payload.
func nonMergedPRPayload(title, body, mergeCommitSHA string) []byte {
	ev := map[string]interface{}{
		"action": "closed",
		"pull_request": map[string]interface{}{
			"title":            title,
			"body":             body,
			"merged":           false,
			"merge_commit_sha": mergeCommitSHA,
			"html_url":         webhookRepo + "/pull/2",
		},
		"repository": map[string]interface{}{
			"html_url":  webhookRepo,
			"clone_url": webhookRepo + ".git",
			"full_name": "test/webhook-test",
		},
	}
	data, _ := json.Marshal(ev)
	return data
}

// prPayload creates a GitHub merged PR event JSON payload.
func prPayload(title, body, mergeCommitSHA string) []byte {
	ev := map[string]interface{}{
		"action": "closed",
		"pull_request": map[string]interface{}{
			"title":            title,
			"body":             body,
			"merged":           true,
			"merge_commit_sha": mergeCommitSHA,
			"html_url":         webhookRepo + "/pull/1",
		},
		"repository": map[string]interface{}{
			"html_url":  webhookRepo,
			"clone_url": webhookRepo + ".git",
			"full_name": "test/webhook-test",
		},
	}
	data, _ := json.Marshal(ev)
	return data
}

// getEntity queries an entity by primary key and returns it as a map.
func getEntity(t *testing.T, client *mocks.BugsClient, endpoint, model, idField, id string) map[string]interface{} {
	t.Helper()
	q := mocks.L8QueryText(fmt.Sprintf("select * from %s where %s=%s", model, idField, id))
	resp, err := client.Get(endpoint, q)
	if err != nil {
		t.Fatalf("Failed to GET %s %s: %v", model, id, err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &data); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	list, ok := data["list"]
	if !ok {
		t.Fatalf("No 'list' field in %s response", model)
	}
	items, ok := list.([]interface{})
	if !ok || len(items) == 0 {
		t.Fatalf("Empty list in %s response for id=%s", model, id)
	}
	return items[0].(map[string]interface{})
}

// assertStatus checks that a status value matches the expected numeric or string representation.
func assertStatus(t *testing.T, got interface{}, expectedNum float64, expectedStr string) {
	t.Helper()
	switch v := got.(type) {
	case float64:
		if v != expectedNum {
			t.Fatalf("expected status %v, got %v", expectedNum, v)
		}
	case string:
		if v != expectedStr {
			t.Fatalf("expected status %s, got %s", expectedStr, v)
		}
	default:
		t.Fatalf("unexpected status type %T: %v", got, got)
	}
}
