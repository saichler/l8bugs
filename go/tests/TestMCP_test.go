package tests

import (
	"strings"
	"testing"

	"github.com/saichler/l8bugs/go/bugs/mcp"
	"github.com/saichler/l8types/go/ifs"
)

func testMCP(t *testing.T, vnic ifs.IVNic) {
	server := mcp.NewServer(vnic)

	testMCPListBugs(t, server)
	testMCPListFeatures(t, server)
	testMCPListWithProjectFilter(t, server)
	testMCPReadBug(t, server)
	testMCPReadNotFound(t, server)
	testMCPReadMissingID(t, server)
	testMCPCreateBug(t, server)
	testMCPCreateFeature(t, server)
	testMCPCreateMissingRequired(t, server)
	testMCPUpdateBug(t, server)
	testMCPUpdateNotFound(t, server)
	testMCPAddComment(t, server)
	testMCPAddCommentMissingBody(t, server)
	testMCPAddCommentNotFound(t, server)
	testMCPSearch(t, server)
	testMCPSearchNoResults(t, server)
}

// --- list_issues ---

func testMCPListBugs(t *testing.T, s *mcp.Server) {
	result, err := s.CallTool("list_issues", map[string]interface{}{"type": "bug"})
	if err != nil {
		t.Fatalf("list_issues(bug) failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("list_issues(bug) returned error: %s", result.Content[0].Text)
	}
	if len(result.Content) == 0 || result.Content[0].Text == "" {
		t.Fatal("list_issues(bug) returned empty content")
	}
}

func testMCPListFeatures(t *testing.T, s *mcp.Server) {
	result, err := s.CallTool("list_issues", map[string]interface{}{"type": "feature"})
	if err != nil {
		t.Fatalf("list_issues(feature) failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("list_issues(feature) returned error: %s", result.Content[0].Text)
	}
	if len(result.Content) == 0 || result.Content[0].Text == "" {
		t.Fatal("list_issues(feature) returned empty content")
	}
}

func testMCPListWithProjectFilter(t *testing.T, s *mcp.Server) {
	if len(testStore.ProjectIDs) == 0 {
		t.Fatal("no project IDs available for filter test")
	}
	result, err := s.CallTool("list_issues", map[string]interface{}{
		"type":       "bug",
		"project_id": testStore.ProjectIDs[0],
	})
	if err != nil {
		t.Fatalf("list_issues with project filter failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("list_issues with project filter returned error: %s", result.Content[0].Text)
	}
}

// --- read_issue ---

func testMCPReadBug(t *testing.T, s *mcp.Server) {
	if len(testStore.BugIDs) == 0 {
		t.Fatal("no bug IDs available for read test")
	}
	result, err := s.CallTool("read_issue", map[string]interface{}{
		"issue_id": testStore.BugIDs[0],
	})
	if err != nil {
		t.Fatalf("read_issue failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("read_issue returned error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, testStore.BugIDs[0]) {
		t.Fatalf("expected bug ID in output, got: %s", result.Content[0].Text)
	}
}

func testMCPReadNotFound(t *testing.T, s *mcp.Server) {
	_, err := s.CallTool("read_issue", map[string]interface{}{
		"issue_id": "nonexistent-id-12345",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent issue")
	}
	if !strings.Contains(err.Error(), "issue not found") {
		t.Fatalf("expected 'issue not found' error, got: %v", err)
	}
}

func testMCPReadMissingID(t *testing.T, s *mcp.Server) {
	_, err := s.CallTool("read_issue", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for missing issue_id")
	}
	if !strings.Contains(err.Error(), "issue_id is required") {
		t.Fatalf("expected 'issue_id is required' error, got: %v", err)
	}
}

// --- create_issue ---

func testMCPCreateBug(t *testing.T, s *mcp.Server) {
	if len(testStore.ProjectIDs) == 0 {
		t.Fatal("no project IDs for create test")
	}
	result, err := s.CallTool("create_issue", map[string]interface{}{
		"type":       "bug",
		"project_id": testStore.ProjectIDs[0],
		"title":      "MCP test bug",
		"priority":   "medium",
		"severity":   "minor",
	})
	if err != nil {
		t.Fatalf("create_issue(bug) failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("create_issue(bug) returned error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Bug created") {
		t.Fatalf("expected 'Bug created' in output, got: %s", result.Content[0].Text)
	}
}

func testMCPCreateFeature(t *testing.T, s *mcp.Server) {
	if len(testStore.ProjectIDs) == 0 {
		t.Fatal("no project IDs for create test")
	}
	result, err := s.CallTool("create_issue", map[string]interface{}{
		"type":       "feature",
		"project_id": testStore.ProjectIDs[0],
		"title":      "MCP test feature",
		"priority":   "high",
	})
	if err != nil {
		t.Fatalf("create_issue(feature) failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("create_issue(feature) returned error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Feature created") {
		t.Fatalf("expected 'Feature created' in output, got: %s", result.Content[0].Text)
	}
}

func testMCPCreateMissingRequired(t *testing.T, s *mcp.Server) {
	_, err := s.CallTool("create_issue", map[string]interface{}{
		"type": "bug",
	})
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

// --- update_issue ---

func testMCPUpdateBug(t *testing.T, s *mcp.Server) {
	if len(testStore.ProjectIDs) == 0 {
		t.Fatal("no project IDs for update test")
	}

	// Create a bug first.
	createResult, err := s.CallTool("create_issue", map[string]interface{}{
		"type":       "bug",
		"project_id": testStore.ProjectIDs[0],
		"title":      "MCP update test bug",
		"priority":   "low",
		"severity":   "trivial",
	})
	if err != nil {
		t.Fatalf("create_issue for update test failed: %v", err)
	}

	// Extract the bug ID from the creation result.
	bugId := extractID(createResult.Content[0].Text, "bugId")
	if bugId == "" {
		t.Fatal("could not extract bugId from create result")
	}

	// Update it.
	result, err := s.CallTool("update_issue", map[string]interface{}{
		"issue_id": bugId,
		"title":    "Updated MCP test bug",
	})
	if err != nil {
		t.Fatalf("update_issue failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("update_issue returned error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Bug updated") {
		t.Fatalf("expected 'Bug updated' in output, got: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Updated MCP test bug") {
		t.Fatalf("expected updated title in output, got: %s", result.Content[0].Text)
	}
}

func testMCPUpdateNotFound(t *testing.T, s *mcp.Server) {
	_, err := s.CallTool("update_issue", map[string]interface{}{
		"issue_id": "nonexistent-id-99999",
		"title":    "Should fail",
	})
	if err == nil {
		t.Fatal("expected error for updating nonexistent issue")
	}
	if !strings.Contains(err.Error(), "issue not found") {
		t.Fatalf("expected 'issue not found' error, got: %v", err)
	}
}

// --- add_comment ---

func testMCPAddComment(t *testing.T, s *mcp.Server) {
	if len(testStore.BugIDs) == 0 {
		t.Fatal("no bug IDs for comment test")
	}
	result, err := s.CallTool("add_comment", map[string]interface{}{
		"issue_id": testStore.BugIDs[0],
		"body":     "MCP test comment",
	})
	if err != nil {
		t.Fatalf("add_comment failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("add_comment returned error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Comment added") {
		t.Fatalf("expected 'Comment added' in output, got: %s", result.Content[0].Text)
	}
}

func testMCPAddCommentMissingBody(t *testing.T, s *mcp.Server) {
	if len(testStore.BugIDs) == 0 {
		t.Fatal("no bug IDs for comment test")
	}
	_, err := s.CallTool("add_comment", map[string]interface{}{
		"issue_id": testStore.BugIDs[0],
	})
	if err == nil {
		t.Fatal("expected error for missing body")
	}
	if !strings.Contains(err.Error(), "body are required") {
		t.Fatalf("expected 'body are required' error, got: %v", err)
	}
}

func testMCPAddCommentNotFound(t *testing.T, s *mcp.Server) {
	_, err := s.CallTool("add_comment", map[string]interface{}{
		"issue_id": "nonexistent-id-77777",
		"body":     "Should fail",
	})
	if err == nil {
		t.Fatal("expected error for comment on nonexistent issue")
	}
	if !strings.Contains(err.Error(), "issue not found") {
		t.Fatalf("expected 'issue not found' error, got: %v", err)
	}
}

// --- search_issues ---

func testMCPSearch(t *testing.T, s *mcp.Server) {
	// Search for a term — verifies the tool executes without error.
	// Note: search uses GetEntities which may return limited results
	// when called through the direct handler path (no L8Query).
	result, err := s.CallTool("search_issues", map[string]interface{}{
		"query": "test",
	})
	if err != nil {
		t.Fatalf("search_issues failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("search_issues returned error: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "Search results") {
		t.Fatalf("expected search results header, got: %s", result.Content[0].Text)
	}
}

func testMCPSearchNoResults(t *testing.T, s *mcp.Server) {
	result, err := s.CallTool("search_issues", map[string]interface{}{
		"query": "zzz_no_match_zzz_98765",
	})
	if err != nil {
		t.Fatalf("search_issues failed: %v", err)
	}
	if !strings.Contains(result.Content[0].Text, "0 match") {
		t.Fatalf("expected 0 matches, got: %s", result.Content[0].Text)
	}
}

// --- helpers ---

// extractID extracts a JSON field value from MCP result text.
// The result text is of the form: "Bug created: {\"bugId\":\"...\", ...}"
func extractID(text, field string) string {
	// Find the field in the JSON portion of the text.
	key := `"` + field + `":"`
	idx := strings.Index(text, key)
	if idx < 0 {
		return ""
	}
	start := idx + len(key)
	end := strings.Index(text[start:], `"`)
	if end < 0 {
		return ""
	}
	return text[start : start+end]
}
