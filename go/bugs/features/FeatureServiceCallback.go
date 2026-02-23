package features

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newFeatureServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.Feature]("Feature",
		func(e *l8bugs.Feature) { common.GenerateID(&e.FeatureId) }).
		Require(func(e *l8bugs.Feature) string { return e.FeatureId }, "FeatureId").
		Require(func(e *l8bugs.Feature) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.Feature) string { return e.Title }, "Title").
		StatusTransition(&common.StatusTransitionConfig[l8bugs.Feature]{
			StatusGetter: func(e *l8bugs.Feature) int32 { return int32(e.Status) },
			StatusSetter: func(e *l8bugs.Feature, s int32) { e.Status = l8bugs.FeatureStatus(s) },
			FilterBuilder: func(e *l8bugs.Feature) *l8bugs.Feature {
				return &l8bugs.Feature{FeatureId: e.FeatureId}
			},
			ServiceName:   ServiceName,
			ServiceArea:   ServiceArea,
			InitialStatus: int32(l8bugs.FeatureStatus_FEATURE_STATUS_PROPOSED),
			Transitions: map[int32][]int32{
				1: {2, 8},
				2: {3, 8, 9},
				3: {4},
				4: {5},
				5: {6, 4},
				6: {7},
				9: {2},
			},
			StatusNames: map[int32]string{
				0: "Unspecified", 1: "Proposed", 2: "Triaged", 3: "Approved",
				4: "In Progress", 5: "In Review", 6: "Done", 7: "Closed",
				8: "Rejected", 9: "Deferred",
			},
		}).
		Build()
}
