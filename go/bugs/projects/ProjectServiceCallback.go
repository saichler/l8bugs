package projects

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newProjectServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return l8common.NewValidation(&l8bugs.BugsProject{}, vnic).
		Require(func(e *l8bugs.BugsProject) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.BugsProject) string { return e.Name }, "Name").
		Require(func(e *l8bugs.BugsProject) string { return e.Key }, "Key").
		Build()
}
