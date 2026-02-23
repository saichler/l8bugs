package assignees

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newAssigneeServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.BugsAssignee]("BugsAssignee",
		func(e *l8bugs.BugsAssignee) { common.GenerateID(&e.AssigneeId) }).
		Require(func(e *l8bugs.BugsAssignee) string { return e.AssigneeId }, "AssigneeId").
		Require(func(e *l8bugs.BugsAssignee) string { return e.Name }, "Name").
		Build()
}
