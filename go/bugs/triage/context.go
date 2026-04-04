package triage

import (
	"fmt"
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
)

func (t *Triager) fetchComponents(projectId string) []string {
	if projectId == "" {
		return nil
	}
	result, err := l8common.GetEntity(projectServiceName, serviceArea,
		&l8bugs.BugsProject{ProjectId: projectId}, t.vnic)
	if err != nil || result == nil {
		return nil
	}
	project := result.(*l8bugs.BugsProject)
	names := make([]string, 0, len(project.Components))
	for _, c := range project.Components {
		names = append(names, c.Name)
	}
	return names
}

func (t *Triager) fetchAssignees(projectId string) []*l8bugs.BugsAssignee {
	results, err := l8common.GetEntities(assigneeService, serviceArea,
		&l8bugs.BugsAssignee{ProjectId: projectId, Active: true}, t.vnic)
	if err != nil {
		fmt.Println("[triage] failed to fetch assignees:", err.Error())
		return nil
	}
	assignees := make([]*l8bugs.BugsAssignee, len(results))
	for i, item := range results {
		assignees[i] = item.(*l8bugs.BugsAssignee)
	}
	return assignees
}

func (t *Triager) fetchBugCandidates(excludeId string) []IssueSummary {
	results, err := l8common.GetEntities(bugServiceName, serviceArea,
		&l8bugs.Bug{Status: l8bugs.BugStatus_BUG_STATUS_OPEN}, t.vnic)
	if err != nil {
		return nil
	}
	summaries := make([]IssueSummary, 0, len(results))
	for _, item := range results {
		b := item.(*l8bugs.Bug)
		if b.BugId == excludeId {
			continue
		}
		summaries = append(summaries, IssueSummary{ID: b.BugId, Title: b.Title})
		if len(summaries) >= maxCandidates {
			break
		}
	}
	return summaries
}

func (t *Triager) fetchFeatureCandidates(excludeId string) []IssueSummary {
	results, err := l8common.GetEntities(featureServiceName, serviceArea,
		&l8bugs.Feature{Status: l8bugs.FeatureStatus_FEATURE_STATUS_PROPOSED}, t.vnic)
	if err != nil {
		return nil
	}
	summaries := make([]IssueSummary, 0, len(results))
	for _, item := range results {
		f := item.(*l8bugs.Feature)
		if f.FeatureId == excludeId {
			continue
		}
		summaries = append(summaries, IssueSummary{ID: f.FeatureId, Title: f.Title})
		if len(summaries) >= maxCandidates {
			break
		}
	}
	return summaries
}
