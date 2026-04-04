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

import (
	"fmt"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"os"
	"strings"
)

// RunAllPhases executes all mock data phases in dependency order
func RunAllPhases(client *BugsClient, store *MockDataStore) {
	runPhase("Phase 1: Foundation", func() error {
		return runFoundation(client, store)
	})
	runPhase("Phase 2: Tracking", func() error {
		return runTracking(client, store)
	})
	runPhase("Phase 3: Planning", func() error {
		return runPlanning(client, store)
	})
}

// Phase 1: Foundation (Projects + Assignees — no dependencies)
func runFoundation(client *BugsClient, store *MockDataStore) error {
	projects := generateProjects()
	ids := extractProjectIDs(projects)
	err := runOp(client, "Projects", "/bugs/20/Project", &l8bugs.BugsProjectList{List: projects}, ids, &store.ProjectIDs)
	if err != nil {
		return err
	}

	assignees := generateAssignees(store)
	aIDs := extractAssigneeIDs(assignees)
	return runOp(client, "Assignees", "/bugs/20/Assignee", &l8bugs.BugsAssigneeList{List: assignees}, aIDs, &store.AssigneeIDs)
}

// Phase 2: Tracking (Bugs + Features — depend on Projects, Assignees)
func runTracking(client *BugsClient, store *MockDataStore) error {
	bugs := generateBugs(store)
	bIDs := extractBugIDs(bugs)
	err := runOp(client, "Bugs", "/bugs/20/Bug", &l8bugs.BugList{List: bugs}, bIDs, &store.BugIDs)
	if err != nil {
		return err
	}

	features := generateFeatures(store)
	fIDs := extractFeatureIDs(features)
	return runOp(client, "Features", "/bugs/20/Feature", &l8bugs.FeatureList{List: features}, fIDs, &store.FeatureIDs)
}

// Phase 3: Planning (Sprints + Digests — depend on Projects)
func runPlanning(client *BugsClient, store *MockDataStore) error {
	sprints := generateSprints(store)
	sIDs := extractSprintIDs(sprints)
	err := runOp(client, "Sprints", "/bugs/20/Sprint", &l8bugs.BugsSprintList{List: sprints}, sIDs, &store.SprintIDs)
	if err != nil {
		return err
	}

	digests := generateDigests(store)
	dIDs := extractDigestIDs(digests)
	return runOp(client, "Digests", "/bugs/20/Digest", &l8bugs.BugsDigestList{List: digests}, dIDs, &store.DigestIDs)
}

// PrintSummary prints the final summary of all generated mock data
func PrintSummary(store *MockDataStore) {
	fmt.Printf("\n=======================\n")
	fmt.Printf("L8Bugs Mock Data Generation Complete!\n")
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  - Projects:  %d\n", len(store.ProjectIDs))
	fmt.Printf("  - Assignees: %d\n", len(store.AssigneeIDs))
	fmt.Printf("  - Bugs:      %d\n", len(store.BugIDs))
	fmt.Printf("  - Features:  %d\n", len(store.FeatureIDs))
	fmt.Printf("  - Sprints:   %d\n", len(store.SprintIDs))
	fmt.Printf("  - Digests:   %d\n", len(store.DigestIDs))
}

// --- Helpers ---

// runOp executes a single phase operation: post + store IDs
func runOp(client *BugsClient, label, endpoint string, list interface{}, ids []string, storeIDs *[]string) error {
	fmt.Printf("  Creating %s...", label)
	resp, err := client.Post(endpoint, list)
	if err != nil {
		fmt.Printf(" FAILED\n")
		fmt.Printf("    Error: %v\n", err)
		if resp != "" {
			fmt.Printf("    Response: %s\n", resp)
		}
		return fmt.Errorf("%s: %w", label, err)
	}
	if storeIDs != nil && ids != nil {
		*storeIDs = append(*storeIDs, ids...)
	}
	fmt.Printf(" %d created\n", len(ids))
	return nil
}

// extractProjectIDs extracts project IDs from a slice
func extractProjectIDs(items []*l8bugs.BugsProject) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ProjectId
	}
	return ids
}

// extractAssigneeIDs extracts assignee IDs from a slice
func extractAssigneeIDs(items []*l8bugs.BugsAssignee) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.AssigneeId
	}
	return ids
}

// extractBugIDs extracts bug IDs from a slice
func extractBugIDs(items []*l8bugs.Bug) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.BugId
	}
	return ids
}

// extractFeatureIDs extracts feature IDs from a slice
func extractFeatureIDs(items []*l8bugs.Feature) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.FeatureId
	}
	return ids
}

// extractSprintIDs extracts sprint IDs from a slice
func extractSprintIDs(items []*l8bugs.BugsSprint) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.SprintId
	}
	return ids
}

// extractDigestIDs extracts digest IDs from a slice
func extractDigestIDs(items []*l8bugs.BugsDigest) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.DigestId
	}
	return ids
}

// runPhase runs a phase function with header formatting
func runPhase(label string, fn func() error) {
	fmt.Printf("\n%s\n", label)
	fmt.Printf("%s\n", strings.Repeat("-", len(label)))
	if err := fn(); err != nil {
		fmt.Printf("%s failed: %v\n", label, err)
		os.Exit(1)
	}
}
