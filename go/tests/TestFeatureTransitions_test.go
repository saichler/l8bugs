package tests

import (
	"github.com/saichler/l8bugs/go/tests/mocks"
	"testing"
)

const featureEndpoint = "/bugs/20/Feature"

func testFeatureTransitions(t *testing.T, client *mocks.BugsClient) {
	testFeatureHappyPath(t, client)
	testFeatureRejectionPath(t, client)
	testFeatureDeferredPath(t, client)
	testFeatureReviewBounce(t, client)
	testFeatureInvalidTransitions(t, client)
}

// testFeatureHappyPath tests: Proposed->Triaged->Approved->InProgress->InReview->Done->Closed
func testFeatureHappyPath(t *testing.T, client *mocks.BugsClient) {
	feature := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, feature); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}

	// Proposed(1) -> Triaged(2) -> Approved(3) -> InProgress(4) -> InReview(5) -> Done(6) -> Closed(7)
	advanceStatus(t, client, featureEndpoint, feature, 2, 3, 4, 5, 6, 7)
}

// testFeatureRejectionPath tests: Proposed->Rejected
func testFeatureRejectionPath(t *testing.T, client *mocks.BugsClient) {
	feature := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, feature); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}

	// Proposed(1) -> Rejected(8)
	feature["status"] = 8
	if err := putEntity(client, featureEndpoint, feature); err != nil {
		t.Fatalf("Feature Proposed->Rejected failed: %v", err)
	}
}

// testFeatureDeferredPath tests: advance to Triaged, Triaged->Deferred, Deferred->Triaged
func testFeatureDeferredPath(t *testing.T, client *mocks.BugsClient) {
	feature := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, feature); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}

	// Proposed(1) -> Triaged(2)
	advanceStatus(t, client, featureEndpoint, feature, 2)

	// Triaged(2) -> Deferred(9)
	feature["status"] = 9
	if err := putEntity(client, featureEndpoint, feature); err != nil {
		t.Fatalf("Feature Triaged->Deferred failed: %v", err)
	}

	// Deferred(9) -> Triaged(2)
	feature["status"] = 2
	if err := putEntity(client, featureEndpoint, feature); err != nil {
		t.Fatalf("Feature Deferred->Triaged failed: %v", err)
	}
}

// testFeatureReviewBounce tests: advance to InReview, InReview->InProgress (send back)
func testFeatureReviewBounce(t *testing.T, client *mocks.BugsClient) {
	feature := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, feature); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}

	// Advance to InReview: Proposed->Triaged->Approved->InProgress->InReview
	advanceStatus(t, client, featureEndpoint, feature, 2, 3, 4, 5)

	// InReview(5) -> InProgress(4)
	feature["status"] = 4
	if err := putEntity(client, featureEndpoint, feature); err != nil {
		t.Fatalf("Feature InReview->InProgress failed: %v", err)
	}
}

// testFeatureInvalidTransitions tests transitions that must be rejected
func testFeatureInvalidTransitions(t *testing.T, client *mocks.BugsClient) {
	// Proposed(1) -> Approved(3) — must go through Triaged
	f1 := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, f1); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}
	f1["status"] = 3
	if err := putEntity(client, featureEndpoint, f1); err == nil {
		t.Fatal("Feature Proposed->Approved should have failed")
	}

	// Proposed(1) -> Done(6) — skipping workflow
	f2 := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, f2); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}
	f2["status"] = 6
	if err := putEntity(client, featureEndpoint, f2); err == nil {
		t.Fatal("Feature Proposed->Done should have failed")
	}

	// Approved(3) -> Done(6) — skipping InProgress/InReview
	f3 := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, f3); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}
	advanceStatus(t, client, featureEndpoint, f3, 2, 3)
	f3["status"] = 6
	if err := putEntity(client, featureEndpoint, f3); err == nil {
		t.Fatal("Feature Approved->Done should have failed")
	}

	// Closed(7) -> Proposed(1) — terminal state
	f4 := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, f4); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}
	advanceStatus(t, client, featureEndpoint, f4, 2, 3, 4, 5, 6, 7)
	f4["status"] = 1
	if err := putEntity(client, featureEndpoint, f4); err == nil {
		t.Fatal("Feature Closed->Proposed should have failed")
	}

	// Rejected(8) -> Triaged(2) — terminal state
	f5 := newFeature(testStore.ProjectIDs[0])
	if _, err := client.Post(featureEndpoint, f5); err != nil {
		t.Fatalf("POST Feature failed: %v", err)
	}
	advanceStatus(t, client, featureEndpoint, f5, 8)
	f5["status"] = 2
	if err := putEntity(client, featureEndpoint, f5); err == nil {
		t.Fatal("Feature Rejected->Triaged should have failed")
	}
}
