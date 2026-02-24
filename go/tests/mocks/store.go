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

// MockDataStore holds generated IDs for all l8bugs entities
type MockDataStore struct {
	// Phase 1: Foundation
	ProjectIDs  []string
	AssigneeIDs []string

	// Phase 2: Tracking
	BugIDs     []string
	FeatureIDs []string

	// Phase 3: Planning
	SprintIDs []string
	DigestIDs []string
}
