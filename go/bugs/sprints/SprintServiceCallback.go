package sprints

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newSprintServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return l8common.NewValidation(&l8bugs.BugsSprint{}, vnic).
		Require(func(e *l8bugs.BugsSprint) string { return e.SprintId }, "SprintId").
		Require(func(e *l8bugs.BugsSprint) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.BugsSprint) string { return e.Name }, "Name").
		StatusTransition(&l8common.StatusTransitionConfig{
			StatusGetter:  func(e interface{}) int32 { return int32(e.(*l8bugs.BugsSprint).Status) },
			StatusSetter:  func(e interface{}, s int32) { e.(*l8bugs.BugsSprint).Status = l8bugs.SprintStatus(s) },
			FilterBuilder: func(e interface{}) interface{} { return &l8bugs.BugsSprint{SprintId: e.(*l8bugs.BugsSprint).SprintId} },
			ServiceName:   ServiceName,
			ServiceArea:   ServiceArea,
			InitialStatus: int32(l8bugs.SprintStatus_SPRINT_STATUS_PLANNING),
			Transitions: map[int32][]int32{
				1: {2},
				2: {3},
			},
			StatusNames: map[int32]string{
				0: "Unspecified", 1: "Planning", 2: "Active", 3: "Completed",
			},
		}).
		DateAfter(
			func(e *l8bugs.BugsSprint) int64 { return e.EndDate },
			func(e *l8bugs.BugsSprint) int64 { return e.StartDate },
			"EndDate", "StartDate").
		Build()
}
