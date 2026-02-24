package triage

import (
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"strings"
	"time"
)

const digestSystemPrompt = `You are a project manager. Given a list of bugs and features with their statuses, generate a concise project digest.

Output four sections, each prefixed with its header on its own line:
SUMMARY:
<2-3 sentence project status overview>

KEY_METRICS:
<bullet points with counts: opened, resolved, in progress, etc.>

BLOCKERS:
<bullet points listing any blocking issues or risks, or "None identified" if none>

ACTION_ITEMS:
<numbered action items for the team based on the current state>`

// DigestResult holds the parsed AI digest output.
type DigestResult struct {
	Summary     string
	KeyMetrics  string
	Blockers    string
	ActionItems string
}

// GenerateDigest creates an AI-powered project summary for a time period.
func (t *Triager) GenerateDigest(projectID string, period l8bugs.DigestPeriod, startDate, endDate int64) (*l8bugs.BugsDigest, error) {
	if !t.Available() {
		return nil, fmt.Errorf("AI digest unavailable: API key not configured")
	}

	bugs, _ := common.GetEntities(bugServiceName, serviceArea, &l8bugs.Bug{ProjectId: projectID}, t.vnic)
	features, _ := common.GetEntities(featureServiceName, serviceArea, &l8bugs.Feature{ProjectId: projectID}, t.vnic)

	// Filter to date range if specified.
	if startDate > 0 {
		bugs = filterBugsByDate(bugs, startDate, endDate)
		features = filterFeaturesByDate(features, startDate, endDate)
	}

	prompt := buildDigestPrompt(bugs, features)
	response, err := t.client.Complete(digestSystemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("digest LLM call failed: %w", err)
	}

	result := parseDigestResponse(response)

	digest := &l8bugs.BugsDigest{
		ProjectId:     projectID,
		Period:        period,
		StartDate:     startDate,
		EndDate:       endDate,
		Summary:       result.Summary,
		KeyMetrics:    result.KeyMetrics,
		Blockers:      result.Blockers,
		ActionItems:   result.ActionItems,
		GeneratedDate: time.Now().Unix(),
	}

	created, err := common.PostEntity("Digest", serviceArea, digest, t.vnic)
	if err != nil {
		return nil, fmt.Errorf("failed to save digest: %w", err)
	}
	fmt.Printf("[triage] digest generated for project %s\n", projectID)
	return created, nil
}

func buildDigestPrompt(bugs []*l8bugs.Bug, features []*l8bugs.Feature) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Project has %d bugs and %d features.\n\n", len(bugs), len(features))

	if len(bugs) > 0 {
		b.WriteString("BUGS:\n")
		for _, bug := range bugs {
			fmt.Fprintf(&b, "- [%s] %s (status=%s, priority=%s)\n",
				bug.BugNumber, bug.Title, bug.Status.String(), bug.Priority.String())
		}
		b.WriteString("\n")
	}

	if len(features) > 0 {
		b.WriteString("FEATURES:\n")
		for _, f := range features {
			fmt.Fprintf(&b, "- [%s] %s (status=%s, priority=%s)\n",
				f.FeatureNumber, f.Title, f.Status.String(), f.Priority.String())
		}
	}

	return b.String()
}

func parseDigestResponse(response string) *DigestResult {
	result := &DigestResult{}
	sections := map[string]*string{
		"SUMMARY:":      &result.Summary,
		"KEY_METRICS:":  &result.KeyMetrics,
		"BLOCKERS:":     &result.Blockers,
		"ACTION_ITEMS:": &result.ActionItems,
	}

	order := []string{"SUMMARY:", "KEY_METRICS:", "BLOCKERS:", "ACTION_ITEMS:"}
	for i, header := range order {
		idx := strings.Index(response, header)
		if idx < 0 {
			continue
		}
		start := idx + len(header)
		end := len(response)
		for j := i + 1; j < len(order); j++ {
			nextIdx := strings.Index(response, order[j])
			if nextIdx > start {
				end = nextIdx
				break
			}
		}
		*sections[header] = strings.TrimSpace(response[start:end])
	}

	return result
}

func filterBugsByDate(bugs []*l8bugs.Bug, start, end int64) []*l8bugs.Bug {
	var filtered []*l8bugs.Bug
	for _, b := range bugs {
		if b.CreatedDate >= start && (end == 0 || b.CreatedDate <= end) {
			filtered = append(filtered, b)
		}
	}
	return filtered
}

func filterFeaturesByDate(features []*l8bugs.Feature, start, end int64) []*l8bugs.Feature {
	var filtered []*l8bugs.Feature
	for _, f := range features {
		if f.CreatedDate >= start && (end == 0 || f.CreatedDate <= end) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}
