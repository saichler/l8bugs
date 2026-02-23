package mcp

func (s *Server) registerTools() {
	s.toolDefs = []ToolDef{
		{
			Name:        "list_issues",
			Description: "List bugs or features with optional filters. Returns a summary of matching issues.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"type":        {Type: "string", Description: "Issue type to list", Enum: []string{"bug", "feature"}},
					"status":      {Type: "string", Description: "Filter by status name (e.g. 'open', 'in_progress', 'resolved')"},
					"priority":    {Type: "string", Description: "Filter by priority (critical, high, medium, low)"},
					"assignee_id": {Type: "string", Description: "Filter by assignee ID"},
					"project_id":  {Type: "string", Description: "Filter by project ID"},
					"limit":       {Type: "number", Description: "Max results to return (default 20)"},
				},
				Required: []string{"type"},
			},
		},
		{
			Name:        "read_issue",
			Description: "Get full details of a bug or feature by ID, including comments and activity.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"issue_id": {Type: "string", Description: "The bug or feature ID"},
				},
				Required: []string{"issue_id"},
			},
		},
		{
			Name:        "create_issue",
			Description: "Create a new bug or feature.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"type":        {Type: "string", Description: "Issue type", Enum: []string{"bug", "feature"}},
					"project_id":  {Type: "string", Description: "Project ID"},
					"title":       {Type: "string", Description: "Issue title"},
					"description": {Type: "string", Description: "Detailed description"},
					"priority":    {Type: "string", Description: "Priority (critical, high, medium, low)"},
					"severity":    {Type: "string", Description: "Bug severity (blocker, major, minor, trivial)"},
					"component":   {Type: "string", Description: "Component name"},
					"assignee_id": {Type: "string", Description: "Assignee ID"},
				},
				Required: []string{"type", "project_id", "title"},
			},
		},
		{
			Name:        "update_issue",
			Description: "Update fields on an existing bug or feature.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"issue_id":    {Type: "string", Description: "The bug or feature ID to update"},
					"title":       {Type: "string", Description: "New title"},
					"description": {Type: "string", Description: "New description"},
					"status":      {Type: "string", Description: "New status name"},
					"priority":    {Type: "string", Description: "New priority"},
					"severity":    {Type: "string", Description: "New severity (bugs only)"},
					"assignee_id": {Type: "string", Description: "New assignee ID"},
					"component":   {Type: "string", Description: "New component"},
					"resolution":  {Type: "string", Description: "Resolution (bugs only: fixed, wont_fix, duplicate, etc.)"},
				},
				Required: []string{"issue_id"},
			},
		},
		{
			Name:        "add_comment",
			Description: "Add a comment to a bug or feature.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"issue_id":    {Type: "string", Description: "The bug or feature ID"},
					"body":        {Type: "string", Description: "Comment body text"},
					"author_type": {Type: "string", Description: "Author type", Enum: []string{"human", "ai_agent", "system"}},
				},
				Required: []string{"issue_id", "body"},
			},
		},
		{
			Name:        "search_issues",
			Description: "Search bugs and features by text query across title and description.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"query":      {Type: "string", Description: "Search text"},
					"type":       {Type: "string", Description: "Limit to bug or feature", Enum: []string{"bug", "feature"}},
					"project_id": {Type: "string", Description: "Limit to a specific project"},
				},
				Required: []string{"query"},
			},
		},
	}

	s.tools["list_issues"] = s.handleListIssues
	s.tools["read_issue"] = s.handleReadIssue
	s.tools["create_issue"] = s.handleCreateIssue
	s.tools["update_issue"] = s.handleUpdateIssue
	s.tools["add_comment"] = s.handleAddComment
	s.tools["search_issues"] = s.handleSearchIssues
}
