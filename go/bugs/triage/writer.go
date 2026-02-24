package triage

import (
	"fmt"
)

// WriteRequest represents an AI writing assistance request.
type WriteRequest struct {
	Action string `json:"action"` // suggest_steps, improve_description, generate_acceptance_criteria, summarize_comments
	Input  string `json:"input"`
	Title  string `json:"title"`
}

// WriteResult represents the AI writing output.
type WriteResult struct {
	Output string `json:"output"`
}

const (
	ActionSuggestSteps       = "suggest_steps"
	ActionImproveDescription = "improve_description"
	ActionGenAcceptance      = "generate_acceptance_criteria"
	ActionSummarizeComments  = "summarize_comments"
)

var writerSystemPrompts = map[string]string{
	ActionSuggestSteps: `You are a QA expert. Given a bug title and description, generate clear, numbered "Steps to Reproduce" that another developer can follow. Be specific about inputs, clicks, and expected state at each step. Output only the steps, no preamble.`,

	ActionImproveDescription: `You are a technical writer. Given an issue title and description, rewrite the description to be clearer, more structured, and actionable. Preserve all technical details. Use short paragraphs. Output only the improved description.`,

	ActionGenAcceptance: `You are a product manager. Given a feature title and description, produce testable acceptance criteria in "Given/When/Then" format. Each criterion should be independently verifiable. Output only the acceptance criteria, numbered.`,

	ActionSummarizeComments: `You are a project coordinator. Given a series of comments from an issue thread, produce a concise summary of key decisions, action items, and open questions. Use bullet points. Output only the summary.`,
}

// AssistWriting generates or improves issue text using AI.
func (t *Triager) AssistWriting(req *WriteRequest) (*WriteResult, error) {
	if !t.Available() {
		return nil, fmt.Errorf("AI writing assistance unavailable: API key not configured")
	}

	systemPrompt, ok := writerSystemPrompts[req.Action]
	if !ok {
		return nil, fmt.Errorf("unknown writing action: %s", req.Action)
	}

	userMessage := buildWriterPrompt(req)

	response, err := t.client.Complete(systemPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("AI writing failed: %w", err)
	}

	return &WriteResult{Output: response}, nil
}

func buildWriterPrompt(req *WriteRequest) string {
	prompt := ""
	if req.Title != "" {
		prompt += "Title: " + req.Title + "\n\n"
	}
	switch req.Action {
	case ActionSuggestSteps:
		prompt += "Bug description:\n" + req.Input
	case ActionImproveDescription:
		prompt += "Current description:\n" + req.Input
	case ActionGenAcceptance:
		prompt += "Feature description:\n" + req.Input
	case ActionSummarizeComments:
		prompt += "Comments thread:\n" + req.Input
	default:
		prompt += req.Input
	}
	return prompt
}
