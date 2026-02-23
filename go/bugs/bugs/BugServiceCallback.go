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
		Build()
}
