package projects

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newProjectServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.BugsProject]("BugsProject",
		func(e *l8bugs.BugsProject) { common.GenerateID(&e.ProjectId) }).
		Require(func(e *l8bugs.BugsProject) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.BugsProject) string { return e.Name }, "Name").
		Require(func(e *l8bugs.BugsProject) string { return e.Key }, "Key").
		Build()
}
