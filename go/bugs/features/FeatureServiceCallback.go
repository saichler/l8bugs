package features

import (
	l8common "github.com/saichler/l8common/go/common"
	"github.com/saichler/l8bugs/go/bugs/triage"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newFeatureServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return l8common.NewValidation(&l8bugs.Feature{}, vnic).
		Require(func(e *l8bugs.Feature) string { return e.FeatureId }, "FeatureId").
		Require(func(e *l8bugs.Feature) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.Feature) string { return e.Title }, "Title").
		StatusTransition(&l8common.StatusTransitionConfig{
			StatusGetter:  func(e interface{}) int32 { return int32(e.(*l8bugs.Feature).Status) },
			StatusSetter:  func(e interface{}, s int32) { e.(*l8bugs.Feature).Status = l8bugs.FeatureStatus(s) },
			FilterBuilder: func(e interface{}) interface{} { return &l8bugs.Feature{FeatureId: e.(*l8bugs.Feature).FeatureId} },
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
		After(func(entity *l8bugs.Feature, action ifs.Action, _ ifs.IVNic) error {
			if action != ifs.POST {
				return nil
			}
			t := triage.Get()
			if t == nil || !t.Available() {
				return nil
			}
			go t.TriageFeature(entity)
			return nil
		}).
		Build()
}
