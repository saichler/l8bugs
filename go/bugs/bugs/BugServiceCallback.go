package bugs

import (
	l8common "github.com/saichler/l8common/go/common"
	"github.com/saichler/l8bugs/go/bugs/triage"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newBugServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return l8common.NewValidation(&l8bugs.Bug{}, vnic).
		Require(func(e *l8bugs.Bug) string { return e.BugId }, "BugId").
		Require(func(e *l8bugs.Bug) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.Bug) string { return e.Title }, "Title").
		StatusTransition(&l8common.StatusTransitionConfig{
			StatusGetter:  func(e interface{}) int32 { return int32(e.(*l8bugs.Bug).Status) },
			StatusSetter:  func(e interface{}, s int32) { e.(*l8bugs.Bug).Status = l8bugs.BugStatus(s) },
			FilterBuilder: func(e interface{}) interface{} { return &l8bugs.Bug{BugId: e.(*l8bugs.Bug).BugId} },
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
		After(func(entity *l8bugs.Bug, action ifs.Action, _ ifs.IVNic) error {
			if action != ifs.POST {
				return nil
			}
			t := triage.Get()
			if t == nil || !t.Available() {
				return nil
			}
			go t.TriageBug(entity)
			return nil
		}).
		Build()
}
