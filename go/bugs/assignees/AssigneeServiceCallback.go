package assignees

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newAssigneeServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return l8common.NewValidation(&l8bugs.BugsAssignee{}, vnic).
		Require(func(e *l8bugs.BugsAssignee) string { return e.AssigneeId }, "AssigneeId").
		Require(func(e *l8bugs.BugsAssignee) string { return e.Name }, "Name").
		Build()
}
