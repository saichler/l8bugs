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

		bugs[i] = &l8bugs.Bug{
			BugId:           fmt.Sprintf("bug-%03d", i+1),
			ProjectId:       store.ProjectIDs[i%len(store.ProjectIDs)],
			BugNumber:       fmt.Sprintf("BUG-%d", 100+i),
			Title:           bugTitles[i],
			Description:     bugDescriptions[i%len(bugDescriptions)],
			Status:          status,
			Priority:        priority,
			Severity:        severity,
			AssigneeId:      store.AssigneeIDs[i%len(store.AssigneeIDs)],
			AssigneeType:    l8bugs.AssigneeType_ASSIGNEE_TYPE_HUMAN,
			Component:       componentNames[i%len(componentNames)],
			Environment:     "Production",
			EstimatedEffort: int32((i%5 + 1) * 2),
			CreatedDate:     now - int64((count-i)*3600),
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

		features[i] = &l8bugs.Feature{
			FeatureId:       fmt.Sprintf("feat-%03d", i+1),
			ProjectId:       store.ProjectIDs[i%len(store.ProjectIDs)],
			FeatureNumber:   fmt.Sprintf("FEAT-%d", 200+i),
			Title:           featureTitles[i],
			Description:     featureDescriptions[i%len(featureDescriptions)],
			Status:          status,
			Priority:        priority,
			AssigneeId:      store.AssigneeIDs[i%len(store.AssigneeIDs)],
			AssigneeType:    l8bugs.AssigneeType_ASSIGNEE_TYPE_HUMAN,
			Component:       componentNames[i%len(componentNames)],
			TargetVersion:   fmt.Sprintf("v1.%d.0", i/3+1),
			EstimatedEffort: int32((i%4 + 1) * 3),
			CreatedDate:     now - int64((count-i)*7200),
		}
	}
	return features
}
