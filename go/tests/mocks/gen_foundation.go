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
// - BugsProject (5 projects)
// - BugsAssignee (8 assignees)

import (
	"fmt"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"time"
)

func generateProjects() []*l8bugs.BugsProject {
	count := len(projectNames)
	projects := make([]*l8bugs.BugsProject, count)
	now := time.Now().Unix()

	for i := 0; i < count; i++ {
		status := l8bugs.ProjectStatus_PROJECT_STATUS_ACTIVE
		visibility := l8bugs.ProjectVisibility_PROJECT_VISIBILITY_INTERNAL

		if i == count-1 {
			visibility = l8bugs.ProjectVisibility_PROJECT_VISIBILITY_PRIVATE
		}

		projects[i] = &l8bugs.BugsProject{
			ProjectId:     fmt.Sprintf("proj-%03d", i+1),
			Name:          projectNames[i],
			Key:           projectKeys[i],
			Description:   projectDescriptions[i%len(projectDescriptions)],
			Status:        status,
			Visibility:    visibility,
			RepositoryUrl: fmt.Sprintf("https://github.com/example/%s", projectKeys[i]),
			CreatedDate:   now - int64((count-i)*86400),
			Labels: []*l8bugs.Label{
				{LabelId: fmt.Sprintf("lbl-%03d-1", i+1), Name: "bug", Color: "#d73a4a"},
				{LabelId: fmt.Sprintf("lbl-%03d-2", i+1), Name: "enhancement", Color: "#a2eeef"},
			},
			Components: []*l8bugs.Component{
				{
					ComponentId: fmt.Sprintf("cmp-%03d-1", i+1),
					Name:        componentNames[i%len(componentNames)],
					Description: fmt.Sprintf("Component for %s", projectNames[i]),
				},
			},
		}
	}
	return projects
}

func generateAssignees(store *MockDataStore) []*l8bugs.BugsAssignee {
	count := len(assigneeNames)
	assignees := make([]*l8bugs.BugsAssignee, count)

	for i := 0; i < count; i++ {
		assigneeType := l8bugs.AssigneeType_ASSIGNEE_TYPE_HUMAN
		if i == count-1 {
			assigneeType = l8bugs.AssigneeType_ASSIGNEE_TYPE_AI_AGENT
		}

		assignees[i] = &l8bugs.BugsAssignee{
			AssigneeId:   fmt.Sprintf("asgn-%03d", i+1),
			Name:         assigneeNames[i],
			Email:        assigneeEmails[i],
			AssigneeType: assigneeType,
			ProjectId:    store.ProjectIDs[i%len(store.ProjectIDs)],
			Active:       true,
		}
	}
	return assignees
}
