package triage

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BugTriageResult struct {
	Priority    int32    `json:"priority"`
	Severity    int32    `json:"severity"`
	Component   string   `json:"component"`
	AssigneeID  string   `json:"assignee_id"`
	Confidence  int32    `json:"confidence"`
	RootCause   string   `json:"root_cause"`
	DuplicateOf string   `json:"duplicate_of"`
	RelatedIDs  []string `json:"related_ids"`
}

type FeatureTriageResult struct {
	Priority   int32    `json:"priority"`
	Component  string   `json:"component"`
	AssigneeID string   `json:"assignee_id"`
	Confidence int32    `json:"confidence"`
	Breakdown  string   `json:"breakdown"`
	RelatedIDs []string `json:"related_ids"`
}

func ParseBugTriageResponse(response string) (*BugTriageResult, error) {
	cleaned := extractJSON(response)
	var result BugTriageResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse bug triage response: %w", err)
	}
	clampRange(&result.Priority, 1, 4)
	clampRange(&result.Severity, 1, 4)
	clampRange(&result.Confidence, 0, 100)
	return &result, nil
}

func ParseFeatureTriageResponse(response string) (*FeatureTriageResult, error) {
	cleaned := extractJSON(response)
	var result FeatureTriageResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse feature triage response: %w", err)
	}
	clampRange(&result.Priority, 1, 4)
	clampRange(&result.Confidence, 0, 100)
	return &result, nil
}

// extractJSON strips markdown fencing if present.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		start, end := 1, len(lines)-1
		if end > start && strings.HasPrefix(lines[end], "```") {
			s = strings.Join(lines[start:end], "\n")
		}
	}
	return strings.TrimSpace(s)
}

func clampRange(v *int32, min, max int32) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}
