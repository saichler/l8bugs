package mcp

import (
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/common"
	"github.com/saichler/l8bugs/go/bugs/triage"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

const (
	bugService     = "Bug"
	featureService = "Feature"
	serviceArea    = byte(20)
	defaultLimit   = 20
)

var marshaler = protojson.MarshalOptions{EmitUnpopulated: false, UseProtoNames: false}

func textResult(text string) *CallToolResult {
	return &CallToolResult{Content: []ContentBlock{{Type: "text", Text: text}}}
}

func protoToText(msg proto.Message) string {
	data, err := marshaler.Marshal(msg)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"marshal failed: %s\"}", err)
	}
	return string(data)
}

func getStr(args map[string]interface{}, key string) string {
	v, ok := args[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func getInt(args map[string]interface{}, key string, def int) int {
	v, ok := args[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	}
	return def
}

// --- list_issues ---

func (s *Server) handleListIssues(args map[string]interface{}) (*CallToolResult, error) {
	issueType := getStr(args, "type")
	limit := getInt(args, "limit", defaultLimit)

	if issueType == "feature" {
		return s.listFeatures(args, limit)
	}
	return s.listBugs(args, limit)
}

func (s *Server) listBugs(args map[string]interface{}, limit int) (*CallToolResult, error) {
	filter := &l8bugs.Bug{}
	if v := getStr(args, "project_id"); v != "" {
		filter.ProjectId = v
	}
	if v := getStr(args, "assignee_id"); v != "" {
		filter.AssigneeId = v
	}
	if v := getStr(args, "status"); v != "" {
		filter.Status = parseBugStatus(v)
	}
	if v := getStr(args, "priority"); v != "" {
		filter.Priority = parsePriority(v)
	}

	bugs, err := common.GetEntities(bugService, serviceArea, filter, s.vnic)
	if err != nil {
		return nil, fmt.Errorf("failed to list bugs: %w", err)
	}

	if len(bugs) > limit {
		bugs = bugs[:limit]
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Found %d bug(s):\n\n", len(bugs))
	for _, bug := range bugs {
		fmt.Fprintf(&b, "- **%s** [%s] %s (priority=%s, status=%s)\n",
			bug.BugId, bug.BugNumber, bug.Title,
			bug.Priority.String(), bug.Status.String())
	}
	return textResult(b.String()), nil
}

func (s *Server) listFeatures(args map[string]interface{}, limit int) (*CallToolResult, error) {
	filter := &l8bugs.Feature{}
	if v := getStr(args, "project_id"); v != "" {
		filter.ProjectId = v
	}
	if v := getStr(args, "assignee_id"); v != "" {
		filter.AssigneeId = v
	}
	if v := getStr(args, "status"); v != "" {
		filter.Status = parseFeatureStatus(v)
	}
	if v := getStr(args, "priority"); v != "" {
		filter.Priority = parsePriority(v)
	}

	features, err := common.GetEntities(featureService, serviceArea, filter, s.vnic)
	if err != nil {
		return nil, fmt.Errorf("failed to list features: %w", err)
	}

	if len(features) > limit {
		features = features[:limit]
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Found %d feature(s):\n\n", len(features))
	for _, f := range features {
		fmt.Fprintf(&b, "- **%s** [%s] %s (priority=%s, status=%s)\n",
			f.FeatureId, f.FeatureNumber, f.Title,
			f.Priority.String(), f.Status.String())
	}
	return textResult(b.String()), nil
}

// --- read_issue ---

func (s *Server) handleReadIssue(args map[string]interface{}) (*CallToolResult, error) {
	id := getStr(args, "issue_id")
	if id == "" {
		return nil, fmt.Errorf("issue_id is required")
	}

	bug, err := common.GetEntity(bugService, serviceArea, &l8bugs.Bug{BugId: id}, s.vnic)
	if err == nil && bug != nil {
		return textResult(protoToText(bug)), nil
	}

	feature, err := common.GetEntity(featureService, serviceArea, &l8bugs.Feature{FeatureId: id}, s.vnic)
	if err == nil && feature != nil {
		return textResult(protoToText(feature)), nil
	}

	return nil, fmt.Errorf("issue not found: %s", id)
}

// --- create_issue ---

func (s *Server) handleCreateIssue(args map[string]interface{}) (*CallToolResult, error) {
	issueType := getStr(args, "type")
	if issueType == "feature" {
		return s.createFeature(args)
	}
	return s.createBug(args)
}

func (s *Server) createBug(args map[string]interface{}) (*CallToolResult, error) {
	bug := &l8bugs.Bug{
		ProjectId:   getStr(args, "project_id"),
		Title:       getStr(args, "title"),
		Description: getStr(args, "description"),
		Component:   getStr(args, "component"),
		AssigneeId:  getStr(args, "assignee_id"),
		Status:      l8bugs.BugStatus_BUG_STATUS_OPEN,
		CreatedDate: time.Now().Unix(),
	}
	if v := getStr(args, "priority"); v != "" {
		bug.Priority = parsePriority(v)
	}
	if v := getStr(args, "severity"); v != "" {
		bug.Severity = parseSeverity(v)
	}

	created, err := common.PostEntity(bugService, serviceArea, bug, s.vnic)
	if err != nil {
		return nil, fmt.Errorf("failed to create bug: %w", err)
	}
	return textResult(fmt.Sprintf("Bug created: %s", protoToText(created))), nil
}

func (s *Server) createFeature(args map[string]interface{}) (*CallToolResult, error) {
	feature := &l8bugs.Feature{
		ProjectId:   getStr(args, "project_id"),
		Title:       getStr(args, "title"),
		Description: getStr(args, "description"),
		Component:   getStr(args, "component"),
		AssigneeId:  getStr(args, "assignee_id"),
		Status:      l8bugs.FeatureStatus_FEATURE_STATUS_PROPOSED,
		CreatedDate: time.Now().Unix(),
	}
	if v := getStr(args, "priority"); v != "" {
		feature.Priority = parsePriority(v)
	}

	created, err := common.PostEntity(featureService, serviceArea, feature, s.vnic)
	if err != nil {
		return nil, fmt.Errorf("failed to create feature: %w", err)
	}
	return textResult(fmt.Sprintf("Feature created: %s", protoToText(created))), nil
}

// --- update_issue ---

func (s *Server) handleUpdateIssue(args map[string]interface{}) (*CallToolResult, error) {
	id := getStr(args, "issue_id")
	if id == "" {
		return nil, fmt.Errorf("issue_id is required")
	}

	bug, _ := common.GetEntity(bugService, serviceArea, &l8bugs.Bug{BugId: id}, s.vnic)
	if bug != nil {
		return s.updateBug(bug, args)
	}

	feature, _ := common.GetEntity(featureService, serviceArea, &l8bugs.Feature{FeatureId: id}, s.vnic)
	if feature != nil {
		return s.updateFeature(feature, args)
	}

	return nil, fmt.Errorf("issue not found: %s", id)
}

func (s *Server) updateBug(bug *l8bugs.Bug, args map[string]interface{}) (*CallToolResult, error) {
	if v := getStr(args, "title"); v != "" {
		bug.Title = v
	}
	if v := getStr(args, "description"); v != "" {
		bug.Description = v
	}
	if v := getStr(args, "status"); v != "" {
		bug.Status = parseBugStatus(v)
	}
	if v := getStr(args, "priority"); v != "" {
		bug.Priority = parsePriority(v)
	}
	if v := getStr(args, "severity"); v != "" {
		bug.Severity = parseSeverity(v)
	}
	if v := getStr(args, "assignee_id"); v != "" {
		bug.AssigneeId = v
	}
	if v := getStr(args, "component"); v != "" {
		bug.Component = v
	}
	if v := getStr(args, "resolution"); v != "" {
		bug.Resolution = parseResolution(v)
	}

	if err := common.PutEntity(bugService, serviceArea, bug, s.vnic); err != nil {
		return nil, fmt.Errorf("failed to update bug: %w", err)
	}
	return textResult(fmt.Sprintf("Bug updated: %s", protoToText(bug))), nil
}

func (s *Server) updateFeature(feature *l8bugs.Feature, args map[string]interface{}) (*CallToolResult, error) {
	if v := getStr(args, "title"); v != "" {
		feature.Title = v
	}
	if v := getStr(args, "description"); v != "" {
		feature.Description = v
	}
	if v := getStr(args, "status"); v != "" {
		feature.Status = parseFeatureStatus(v)
	}
	if v := getStr(args, "priority"); v != "" {
		feature.Priority = parsePriority(v)
	}
	if v := getStr(args, "assignee_id"); v != "" {
		feature.AssigneeId = v
	}
	if v := getStr(args, "component"); v != "" {
		feature.Component = v
	}

	if err := common.PutEntity(featureService, serviceArea, feature, s.vnic); err != nil {
		return nil, fmt.Errorf("failed to update feature: %w", err)
	}
	return textResult(fmt.Sprintf("Feature updated: %s", protoToText(feature))), nil
}

// --- add_comment ---

func (s *Server) handleAddComment(args map[string]interface{}) (*CallToolResult, error) {
	id := getStr(args, "issue_id")
	body := getStr(args, "body")
	if id == "" || body == "" {
		return nil, fmt.Errorf("issue_id and body are required")
	}

	authorType := l8bugs.AuthorType_AUTHOR_TYPE_AI_AGENT
	if v := getStr(args, "author_type"); v != "" {
		authorType = parseAuthorType(v)
	}

	comment := &l8bugs.Comment{
		AuthorType:  authorType,
		Body:        body,
		CreatedDate: time.Now().Unix(),
	}

	bug, _ := common.GetEntity(bugService, serviceArea, &l8bugs.Bug{BugId: id}, s.vnic)
	if bug != nil {
		bug.Comments = append(bug.Comments, comment)
		if err := common.PutEntity(bugService, serviceArea, bug, s.vnic); err != nil {
			return nil, fmt.Errorf("failed to add comment: %w", err)
		}
		return textResult("Comment added to bug " + id), nil
	}

	feature, _ := common.GetEntity(featureService, serviceArea, &l8bugs.Feature{FeatureId: id}, s.vnic)
	if feature != nil {
		feature.Comments = append(feature.Comments, comment)
		if err := common.PutEntity(featureService, serviceArea, feature, s.vnic); err != nil {
			return nil, fmt.Errorf("failed to add comment: %w", err)
		}
		return textResult("Comment added to feature " + id), nil
	}

	return nil, fmt.Errorf("issue not found: %s", id)
}

// --- search_issues ---

func (s *Server) handleSearchIssues(args map[string]interface{}) (*CallToolResult, error) {
	query := strings.ToLower(getStr(args, "query"))
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	issueType := getStr(args, "type")
	projectID := getStr(args, "project_id")

	var b strings.Builder
	count := 0

	if issueType == "" || issueType == "bug" {
		filter := &l8bugs.Bug{}
		if projectID != "" {
			filter.ProjectId = projectID
		}
		bugs, _ := common.GetEntities(bugService, serviceArea, filter, s.vnic)
		for _, bug := range bugs {
			if matchesQuery(query, bug.Title, bug.Description) {
				fmt.Fprintf(&b, "- [Bug] **%s** %s (status=%s)\n", bug.BugId, bug.Title, bug.Status.String())
				count++
			}
		}
	}

	if issueType == "" || issueType == "feature" {
		filter := &l8bugs.Feature{}
		if projectID != "" {
			filter.ProjectId = projectID
		}
		features, _ := common.GetEntities(featureService, serviceArea, filter, s.vnic)
		for _, f := range features {
			if matchesQuery(query, f.Title, f.Description) {
				fmt.Fprintf(&b, "- [Feature] **%s** %s (status=%s)\n", f.FeatureId, f.Title, f.Status.String())
				count++
			}
		}
	}

	header := fmt.Sprintf("Search results for \"%s\": %d match(es)\n\n", getStr(args, "query"), count)
	return textResult(header + b.String()), nil
}

func matchesQuery(query, title, description string) bool {
	return strings.Contains(strings.ToLower(title), query) ||
		strings.Contains(strings.ToLower(description), query)
}

// --- assist_writing ---

func (s *Server) handleAssistWriting(args map[string]interface{}) (*CallToolResult, error) {
	action := getStr(args, "action")
	input := getStr(args, "input")
	title := getStr(args, "title")

	if action == "" {
		return nil, fmt.Errorf("action is required")
	}
	if input == "" {
		return nil, fmt.Errorf("input is required")
	}

	triager := triage.Get()
	if triager == nil || !triager.Available() {
		return nil, fmt.Errorf("AI writing assistance unavailable")
	}

	result, err := triager.AssistWriting(&triage.WriteRequest{
		Action: action,
		Input:  input,
		Title:  title,
	})
	if err != nil {
		return nil, err
	}

	return textResult(result.Output), nil
}

