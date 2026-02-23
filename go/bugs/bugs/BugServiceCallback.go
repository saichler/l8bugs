package bugs

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newBugServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.Bug]("Bug",
		func(e *l8bugs.Bug) { common.GenerateID(&e.BugId) }).
		Require(func(e *l8bugs.Bug) string { return e.BugId }, "BugId").
		Require(func(e *l8bugs.Bug) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.Bug) string { return e.Title }, "Title").
		StatusTransition(&common.StatusTransitionConfig[l8bugs.Bug]{
			StatusGetter: func(e *l8bugs.Bug) int32 { return int32(e.Status) },
			StatusSetter: func(e *l8bugs.Bug, s int32) { e.Status = l8bugs.BugStatus(s) },
			FilterBuilder: func(e *l8bugs.Bug) *l8bugs.Bug {
				return &l8bugs.Bug{BugId: e.BugId}
			},
			ServiceName:   ServiceName,
			ServiceArea:   ServiceArea,
			InitialStatus: int32(l8bugs.BugStatus_BUG_STATUS_OPEN),
			Transitions: map[int32][]int32{
				1:  {2, 8, 9, 10},
				2:  {3, 8, 9, 10},
				3:  {4},
				4:  {5, 3},
				5:  {6, 7},
				7:  {1},
			},
			StatusNames: map[int32]string{
				0: "Unspecified", 1: "Open", 2: "Triaged", 3: "In Progress",
				4: "In Review", 5: "Resolved", 6: "Closed", 7: "Reopened",
				8: "Won't Fix", 9: "Duplicate", 10: "Cannot Reproduce",
			},
		}).
		Build()
}
