package sprints

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newSprintServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.BugsSprint]("BugsSprint",
		func(e *l8bugs.BugsSprint) { common.GenerateID(&e.SprintId) }).
		Require(func(e *l8bugs.BugsSprint) string { return e.SprintId }, "SprintId").
		Require(func(e *l8bugs.BugsSprint) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.BugsSprint) string { return e.Name }, "Name").
		StatusTransition(&common.StatusTransitionConfig[l8bugs.BugsSprint]{
			StatusGetter: func(e *l8bugs.BugsSprint) int32 { return int32(e.Status) },
			StatusSetter: func(e *l8bugs.BugsSprint, s int32) { e.Status = l8bugs.SprintStatus(s) },
			FilterBuilder: func(e *l8bugs.BugsSprint) *l8bugs.BugsSprint {
				return &l8bugs.BugsSprint{SprintId: e.SprintId}
			},
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
