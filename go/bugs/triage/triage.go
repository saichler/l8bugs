package triage

import (
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	bugServiceName     = "Bug"
	featureServiceName = "Feature"
	projectServiceName = "Project"
	assigneeService    = "Assignee"
	serviceArea        = byte(20)
	maxCandidates      = 50
)

type Triager struct {
	client *Client
	vnic   ifs.IVNic
}

var globalTriager *Triager

func Initialize(vnic ifs.IVNic) {
	client := NewClient()
	if client.Available() {
		fmt.Println("[triage] AI triage initialized (model:", client.model+")")
	} else {
		fmt.Println("[triage] AI triage disabled (no L8BUGS_ANTHROPIC_API_KEY)")
	}
	globalTriager = &Triager{client: client, vnic: vnic}
}

func Get() *Triager {
	return globalTriager
}

func (t *Triager) Available() bool {
	return t.client != nil && t.client.Available()
}

func (t *Triager) TriageBug(bug *l8bugs.Bug) {
	if err := t.triageBug(bug); err != nil {
		fmt.Println("[triage] bug triage failed:", err.Error())
		t.markBugFailed(bug, err.Error())
	}
}

func (t *Triager) TriageFeature(feature *l8bugs.Feature) {
	if err := t.triageFeature(feature); err != nil {
		fmt.Println("[triage] feature triage failed:", err.Error())
		t.markFeatureFailed(feature, err.Error())
	}
}

func (t *Triager) triageBug(bug *l8bugs.Bug) error {
	bug.TriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_IN_PROGRESS
	if err := common.PutEntity(bugServiceName, serviceArea, bug, t.vnic); err != nil {
		return fmt.Errorf("failed to set triage in-progress: %w", err)
	}

	components := t.fetchComponents(bug.ProjectId)
	assignees := t.fetchAssignees(bug.ProjectId)
	candidates := t.fetchBugCandidates(bug.BugId)

	prompt := BuildBugTriagePrompt(bug, components, assignees, candidates)
	response, err := t.client.Complete(bugSystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	result, err := ParseBugTriageResponse(response)
	if err != nil {
		return fmt.Errorf("failed to parse LLM response: %w", err)
	}

	bug.AiSuggestedPriority = l8bugs.Priority(result.Priority)
	bug.AiSuggestedSeverity = l8bugs.Severity(result.Severity)
	bug.AiSuggestedComponent = result.Component
	bug.AiSuggestedAssigneeId = result.AssigneeID
	bug.AiConfidence = result.Confidence
	bug.AiRootCause = result.RootCause
	bug.TriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_COMPLETED
	bug.TriageError = ""

	if result.DuplicateOf != "" {
		bug.DuplicateOfId = result.DuplicateOf
	}
	if len(result.RelatedIDs) > 0 {
		bug.RelatedBugIds = result.RelatedIDs
	}

	if err := common.PutEntity(bugServiceName, serviceArea, bug, t.vnic); err != nil {
		return fmt.Errorf("failed to save triage results: %w", err)
	}

	fmt.Printf("[triage] bug %s triaged (confidence: %d%%)\n", bug.BugId, result.Confidence)
	return nil
}

func (t *Triager) triageFeature(feature *l8bugs.Feature) error {
	feature.TriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_IN_PROGRESS
	if err := common.PutEntity(featureServiceName, serviceArea, feature, t.vnic); err != nil {
		return fmt.Errorf("failed to set triage in-progress: %w", err)
	}

	components := t.fetchComponents(feature.ProjectId)
	assignees := t.fetchAssignees(feature.ProjectId)
	candidates := t.fetchFeatureCandidates(feature.FeatureId)

	prompt := BuildFeatureTriagePrompt(feature, components, assignees, candidates)
	response, err := t.client.Complete(featureSystemPrompt, prompt)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	result, err := ParseFeatureTriageResponse(response)
	if err != nil {
		return fmt.Errorf("failed to parse LLM response: %w", err)
	}

	feature.AiSuggestedPriority = l8bugs.Priority(result.Priority)
	feature.AiSuggestedComponent = result.Component
	feature.AiSuggestedAssigneeId = result.AssigneeID
	feature.AiConfidence = result.Confidence
	feature.AiBreakdown = result.Breakdown
	feature.TriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_COMPLETED
	feature.TriageError = ""

	if len(result.RelatedIDs) > 0 {
		feature.RelatedFeatureIds = result.RelatedIDs
	}

	if err := common.PutEntity(featureServiceName, serviceArea, feature, t.vnic); err != nil {
		return fmt.Errorf("failed to save triage results: %w", err)
	}

	fmt.Printf("[triage] feature %s triaged (confidence: %d%%)\n", feature.FeatureId, result.Confidence)
	return nil
}

func (t *Triager) markBugFailed(bug *l8bugs.Bug, errMsg string) {
	bug.TriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_FAILED
	bug.TriageError = errMsg
	_ = common.PutEntity(bugServiceName, serviceArea, bug, t.vnic)
}

func (t *Triager) markFeatureFailed(feature *l8bugs.Feature, errMsg string) {
	feature.TriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_FAILED
	feature.TriageError = errMsg
	_ = common.PutEntity(featureServiceName, serviceArea, feature, t.vnic)
}
