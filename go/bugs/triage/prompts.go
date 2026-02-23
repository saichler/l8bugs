package triage

import (
	"fmt"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"strings"
)

const bugSystemPrompt = `You are an AI triage assistant for a bug tracking system.
Analyze the bug report and provide classification in JSON format.
Respond ONLY with a JSON object, no markdown fencing, no explanation.

Priority values: 1=Critical, 2=High, 3=Medium, 4=Low
Severity values: 1=Blocker, 2=Major, 3=Minor, 4=Trivial
Confidence: 0-100 (how confident you are in your analysis)`

const featureSystemPrompt = `You are an AI triage assistant for a feature tracking system.
Analyze the feature request and provide classification in JSON format.
Respond ONLY with a JSON object, no markdown fencing, no explanation.

Priority values: 1=Critical, 2=High, 3=Medium, 4=Low
Confidence: 0-100 (how confident you are in your analysis)`

type IssueSummary struct {
	ID    string
	Title string
}

func BuildBugTriagePrompt(bug *l8bugs.Bug, components []string, assignees []*l8bugs.BugsAssignee, candidates []IssueSummary) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Bug Title: %s\n", bug.Title))
	if bug.Description != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", bug.Description))
	}
	if bug.StepsToReproduce != "" {
		b.WriteString(fmt.Sprintf("Steps to Reproduce: %s\n", bug.StepsToReproduce))
	}
	if bug.ExpectedBehavior != "" {
		b.WriteString(fmt.Sprintf("Expected: %s\n", bug.ExpectedBehavior))
	}
	if bug.ActualBehavior != "" {
		b.WriteString(fmt.Sprintf("Actual: %s\n", bug.ActualBehavior))
	}
	if bug.StackTrace != "" {
		b.WriteString(fmt.Sprintf("Stack Trace: %s\n", bug.StackTrace))
	}
	if bug.Environment != "" {
		b.WriteString(fmt.Sprintf("Environment: %s\n", bug.Environment))
	}

	if len(components) > 0 {
		b.WriteString(fmt.Sprintf("\nKnown components: %s\n", strings.Join(components, ", ")))
	}

	if len(assignees) > 0 {
		b.WriteString("\nAvailable assignees:\n")
		for _, a := range assignees {
			b.WriteString(fmt.Sprintf("- %s (ID: %s)\n", a.Name, a.AssigneeId))
		}
	}

	if len(candidates) > 0 {
		b.WriteString("\nExisting open bugs (check for duplicates):\n")
		for _, c := range candidates {
			b.WriteString(fmt.Sprintf("- [%s] %s\n", c.ID, c.Title))
		}
	}

	b.WriteString(`
Respond with this exact JSON structure:
{
  "priority": <1-4>,
  "severity": <1-4>,
  "component": "<component name or empty string>",
  "assignee_id": "<assignee ID or empty string>",
  "confidence": <0-100>,
  "root_cause": "<brief root cause analysis or empty string>",
  "duplicate_of": "<bug ID if duplicate found or empty string>",
  "related_ids": ["<bug ID>", ...]
}`)

	return b.String()
}

func BuildFeatureTriagePrompt(feature *l8bugs.Feature, components []string, assignees []*l8bugs.BugsAssignee, candidates []IssueSummary) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Feature Title: %s\n", feature.Title))
	if feature.Description != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", feature.Description))
	}
	if feature.UserStory != "" {
		b.WriteString(fmt.Sprintf("User Story: %s\n", feature.UserStory))
	}
	if feature.AcceptanceCriteria != "" {
		b.WriteString(fmt.Sprintf("Acceptance Criteria: %s\n", feature.AcceptanceCriteria))
	}

	if len(components) > 0 {
		b.WriteString(fmt.Sprintf("\nKnown components: %s\n", strings.Join(components, ", ")))
	}

	if len(assignees) > 0 {
		b.WriteString("\nAvailable assignees:\n")
		for _, a := range assignees {
			b.WriteString(fmt.Sprintf("- %s (ID: %s)\n", a.Name, a.AssigneeId))
		}
	}

	if len(candidates) > 0 {
		b.WriteString("\nExisting features (check for duplicates/related):\n")
		for _, c := range candidates {
			b.WriteString(fmt.Sprintf("- [%s] %s\n", c.ID, c.Title))
		}
	}

	b.WriteString(`
Respond with this exact JSON structure:
{
  "priority": <1-4>,
  "component": "<component name or empty string>",
  "assignee_id": "<assignee ID or empty string>",
  "confidence": <0-100>,
  "breakdown": "<brief effort breakdown or empty string>",
  "related_ids": ["<feature or bug ID>", ...]
}`)

	return b.String()
}
