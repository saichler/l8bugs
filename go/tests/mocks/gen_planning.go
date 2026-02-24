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
// - BugsSprint (6 sprints)
// - BugsDigest (4 digests)

import (
	"fmt"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"time"
)

func generateSprints(store *MockDataStore) []*l8bugs.BugsSprint {
	count := len(sprintNames)
	sprints := make([]*l8bugs.BugsSprint, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		var status l8bugs.SprintStatus
		switch {
		case i < 2:
			status = l8bugs.SprintStatus_SPRINT_STATUS_COMPLETED
		case i == 2:
			status = l8bugs.SprintStatus_SPRINT_STATUS_ACTIVE
		default:
			status = l8bugs.SprintStatus_SPRINT_STATUS_PLANNING
		}

		startDate := now.AddDate(0, 0, (i-2)*14)
		endDate := startDate.AddDate(0, 0, 14)

		sprints[i] = &l8bugs.BugsSprint{
			SprintId:        fmt.Sprintf("spr-%03d", i+1),
			ProjectId:       store.ProjectIDs[i%len(store.ProjectIDs)],
			Name:            sprintNames[i],
			Goal:            sprintGoals[i],
			Status:          status,
			StartDate:       startDate.Unix(),
			EndDate:         endDate.Unix(),
			Capacity:        int32((i + 3) * 10),
			CompletedPoints: int32((i + 1) * 8),
		}
	}
	return sprints
}

func generateDigests(store *MockDataStore) []*l8bugs.BugsDigest {
	count := len(digestSummaries)
	digests := make([]*l8bugs.BugsDigest, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		period := l8bugs.DigestPeriod_DIGEST_PERIOD_WEEKLY
		if i%2 == 0 {
			period = l8bugs.DigestPeriod_DIGEST_PERIOD_DAILY
		}

		startDate := now.AddDate(0, 0, -(count-i)*7)
		endDate := startDate.AddDate(0, 0, 7)

		digests[i] = &l8bugs.BugsDigest{
			DigestId:      fmt.Sprintf("dgst-%03d", i+1),
			ProjectId:     store.ProjectIDs[i%len(store.ProjectIDs)],
			Period:        period,
			StartDate:     startDate.Unix(),
			EndDate:       endDate.Unix(),
			Summary:       digestSummaries[i],
			KeyMetrics:    fmt.Sprintf("Bugs resolved: %d, Features completed: %d", (i+1)*4, (i+1)*2),
			Blockers:      "No major blockers this period.",
			ActionItems:   fmt.Sprintf("Review PR #%d, Update documentation for v1.%d", 100+i*3, i+1),
			GeneratedDate: now.Unix(),
		}
	}
	return digests
}
