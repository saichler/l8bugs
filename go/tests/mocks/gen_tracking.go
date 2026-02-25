/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package mocks

// Generates:
// - Bug (20 bugs)
// - Feature (10 features)

import (
	"fmt"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"time"
)

func generateBugs(store *MockDataStore) []*l8bugs.Bug {
	count := len(bugTitles)
	bugs := make([]*l8bugs.Bug, count)
	now := time.Now().Unix()

	for i := 0; i < count; i++ {
		// Flavorable status distribution: 30% Open, 25% In Progress, 20% Resolved, rest mixed
		var status l8bugs.BugStatus
		switch {
		case i < 6:
			status = l8bugs.BugStatus_BUG_STATUS_OPEN
		case i < 11:
			status = l8bugs.BugStatus_BUG_STATUS_IN_PROGRESS
		case i < 15:
			status = l8bugs.BugStatus_BUG_STATUS_RESOLVED
		case i < 17:
			status = l8bugs.BugStatus_BUG_STATUS_TRIAGED
		case i < 19:
			status = l8bugs.BugStatus_BUG_STATUS_CLOSED
		default:
			status = l8bugs.BugStatus_BUG_STATUS_IN_REVIEW
		}

		// Flavorable priority distribution
		var priority l8bugs.Priority
		switch {
		case i < 3:
			priority = l8bugs.Priority_PRIORITY_CRITICAL
		case i < 8:
			priority = l8bugs.Priority_PRIORITY_HIGH
		case i < 14:
			priority = l8bugs.Priority_PRIORITY_MEDIUM
		default:
			priority = l8bugs.Priority_PRIORITY_LOW
		}

		// Flavorable severity
		var severity l8bugs.Severity
		switch {
		case i < 2:
			severity = l8bugs.Severity_SEVERITY_BLOCKER
		case i < 6:
			severity = l8bugs.Severity_SEVERITY_MAJOR
		case i < 14:
			severity = l8bugs.Severity_SEVERITY_MINOR
		default:
			severity = l8bugs.Severity_SEVERITY_TRIVIAL
		}

		// Resolution and resolvedDate for resolved/closed bugs.
		var resolution l8bugs.Resolution
		var resolvedDate int64
		if status == l8bugs.BugStatus_BUG_STATUS_RESOLVED || status == l8bugs.BugStatus_BUG_STATUS_CLOSED {
			resolution = l8bugs.Resolution_RESOLUTION_FIXED
			resolvedDate = now - int64((count-i)*1200)
		}

		// AI triage fields — simulate completed triage for most bugs.
		triageStatus := l8bugs.TriageStatus_TRIAGE_STATUS_COMPLETED
		aiConfidence := int32(60 + (i%4)*10)           // 60–90
		aiEstimatedEffort := int32((i%6 + 1) * 3)      // 3–18
		aiEffortConfidence := int32(50 + (i%5)*10)      // 50–90
		actualEffort := int32((i%5 + 1) * 2)            // 2–10
		if i < 4 {
			triageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_PENDING
			aiConfidence = 0
			aiEstimatedEffort = 0
			aiEffortConfidence = 0
			actualEffort = 0
		}

		bugs[i] = &l8bugs.Bug{
			BugId:              fmt.Sprintf("bug-%03d", i+1),
			ProjectId:          store.ProjectIDs[i%len(store.ProjectIDs)],
			BugNumber:          fmt.Sprintf("BUG-%d", 100+i),
			Title:              bugTitles[i],
			Description:        bugDescriptions[i%len(bugDescriptions)],
			Status:             status,
			Priority:           priority,
			Severity:           severity,
			AssigneeId:         store.AssigneeIDs[i%len(store.AssigneeIDs)],
			AssigneeType:       l8bugs.AssigneeType_ASSIGNEE_TYPE_HUMAN,
			Component:          componentNames[i%len(componentNames)],
			Environment:        "Production",
			EstimatedEffort:    int32((i%5 + 1) * 2),
			ActualEffort:       actualEffort,
			Resolution:         resolution,
			ResolvedDate:       resolvedDate,
			TriageStatus:       triageStatus,
			AiConfidence:       aiConfidence,
			AiEstimatedEffort:  aiEstimatedEffort,
			AiEffortConfidence: aiEffortConfidence,
			CreatedDate:        now - int64((count-i)*3600),
			Comments: []*l8bugs.Comment{
				{
					CommentId:   fmt.Sprintf("cmt-%03d-1", i+1),
					AuthorId:    store.AssigneeIDs[(i+1)%len(store.AssigneeIDs)],
					AuthorType:  l8bugs.AuthorType_AUTHOR_TYPE_HUMAN,
					Body:        fmt.Sprintf("Investigating %s", bugTitles[i]),
					CreatedDate: now - int64((count-i)*1800),
				},
			},
		}
	}
	return bugs
}

func generateFeatures(store *MockDataStore) []*l8bugs.Feature {
	count := len(featureTitles)
	features := make([]*l8bugs.Feature, count)
	now := time.Now().Unix()

	for i := 0; i < count; i++ {
		// Flavorable status distribution
		var status l8bugs.FeatureStatus
		switch {
		case i < 2:
			status = l8bugs.FeatureStatus_FEATURE_STATUS_PROPOSED
		case i < 5:
			status = l8bugs.FeatureStatus_FEATURE_STATUS_IN_PROGRESS
		case i < 7:
			status = l8bugs.FeatureStatus_FEATURE_STATUS_APPROVED
		case i < 9:
			status = l8bugs.FeatureStatus_FEATURE_STATUS_DONE
		default:
			status = l8bugs.FeatureStatus_FEATURE_STATUS_IN_REVIEW
		}

		var priority l8bugs.Priority
		switch {
		case i < 2:
			priority = l8bugs.Priority_PRIORITY_HIGH
		case i < 6:
			priority = l8bugs.Priority_PRIORITY_MEDIUM
		default:
			priority = l8bugs.Priority_PRIORITY_LOW
		}

		// AI triage fields for features.
		fTriageStatus := l8bugs.TriageStatus_TRIAGE_STATUS_COMPLETED
		fAiConfidence := int32(65 + (i%3)*10) // 65–85
		fAiEstimatedEffort := int32((i%5 + 2) * 3)
		fActualEffort := int32((i%4 + 1) * 3)
		if i < 2 {
			fTriageStatus = l8bugs.TriageStatus_TRIAGE_STATUS_PENDING
			fAiConfidence = 0
			fAiEstimatedEffort = 0
			fActualEffort = 0
		}

		features[i] = &l8bugs.Feature{
			FeatureId:         fmt.Sprintf("feat-%03d", i+1),
			ProjectId:         store.ProjectIDs[i%len(store.ProjectIDs)],
			FeatureNumber:     fmt.Sprintf("FEAT-%d", 200+i),
			Title:             featureTitles[i],
			Description:       featureDescriptions[i%len(featureDescriptions)],
			Status:            status,
			Priority:          priority,
			AssigneeId:        store.AssigneeIDs[i%len(store.AssigneeIDs)],
			AssigneeType:      l8bugs.AssigneeType_ASSIGNEE_TYPE_HUMAN,
			Component:         componentNames[i%len(componentNames)],
			TargetVersion:     fmt.Sprintf("v1.%d.0", i/3+1),
			EstimatedEffort:   int32((i%4 + 1) * 3),
			ActualEffort:      fActualEffort,
			TriageStatus:      fTriageStatus,
			AiConfidence:      fAiConfidence,
			AiEstimatedEffort: fAiEstimatedEffort,
			CreatedDate:       now - int64((count-i)*7200),
		}
	}
	return features
}
