package triage

import (
	"encoding/json"
	"fmt"
)

const effortSystemPrompt = `You are a software estimation expert. Given a bug or feature title, description, and component, estimate the effort in story points (1-13 Fibonacci scale) and provide a confidence percentage (0-100).

Consider:
- Scope and complexity of the change
- Likely number of files affected
- Testing requirements
- Integration risks

Respond in JSON only:
{"estimated_effort": <int>, "confidence": <int>}`

// EffortResult holds the AI effort estimation output.
type EffortResult struct {
	EstimatedEffort int32 `json:"estimated_effort"`
	Confidence      int32 `json:"confidence"`
}

// EstimateEffort uses AI to predict story points for an issue.
func (t *Triager) EstimateEffort(title, description, component string) (*EffortResult, error) {
	if !t.Available() {
		return nil, fmt.Errorf("AI effort estimation unavailable: API key not configured")
	}

	prompt := "Title: " + title + "\n"
	if component != "" {
		prompt += "Component: " + component + "\n"
	}
	prompt += "\nDescription:\n" + description

	response, err := t.client.Complete(effortSystemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("effort estimation LLM call failed: %w", err)
	}

	var result EffortResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse effort response: %w", err)
	}

	return &result, nil
}
