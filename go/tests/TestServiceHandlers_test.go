package tests

import (
	"github.com/saichler/l8bugs/go/bugs/assignees"
	"github.com/saichler/l8bugs/go/bugs/bugs"
	"github.com/saichler/l8bugs/go/bugs/digests"
	"github.com/saichler/l8bugs/go/bugs/features"
	"github.com/saichler/l8bugs/go/bugs/projects"
	"github.com/saichler/l8bugs/go/bugs/sprints"
	"github.com/saichler/l8types/go/ifs"
	"testing"
)

func testServiceHandlers(t *testing.T, vnic ifs.IVNic) {
	log := vnic.Resources().Logger()

	if h, ok := projects.Projects(vnic); !ok || h == nil {
		log.Fail(t, "BugsProject service handler not found")
	}
	if h, ok := assignees.Assignees(vnic); !ok || h == nil {
		log.Fail(t, "BugsAssignee service handler not found")
	}
	if h, ok := bugs.Bugs(vnic); !ok || h == nil {
		log.Fail(t, "Bug service handler not found")
	}
	if h, ok := features.Features(vnic); !ok || h == nil {
		log.Fail(t, "Feature service handler not found")
	}
	if h, ok := sprints.Sprints(vnic); !ok || h == nil {
		log.Fail(t, "BugsSprint service handler not found")
	}
	if h, ok := digests.Digests(vnic); !ok || h == nil {
		log.Fail(t, "BugsDigest service handler not found")
	}
}
