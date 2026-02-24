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

func testServiceGetters(t *testing.T, vnic ifs.IVNic) {
	log := vnic.Resources().Logger()

	if _, err := projects.Project("test-id", vnic); err != nil {
		log.Fail(t, "BugsProject getter failed: ", err.Error())
	}
	if _, err := assignees.Assignee("test-id", vnic); err != nil {
		log.Fail(t, "BugsAssignee getter failed: ", err.Error())
	}
	if _, err := bugs.Bug("test-id", vnic); err != nil {
		log.Fail(t, "Bug getter failed: ", err.Error())
	}
	if _, err := features.Feature("test-id", vnic); err != nil {
		log.Fail(t, "Feature getter failed: ", err.Error())
	}
	if _, err := sprints.Sprint("test-id", vnic); err != nil {
		log.Fail(t, "BugsSprint getter failed: ", err.Error())
	}
	if _, err := digests.Digest("test-id", vnic); err != nil {
		log.Fail(t, "BugsDigest getter failed: ", err.Error())
	}
}
