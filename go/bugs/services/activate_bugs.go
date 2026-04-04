package services

import (
	"github.com/saichler/l8bugs/go/bugs/assignees"
	"github.com/saichler/l8bugs/go/bugs/bugs"
	"github.com/saichler/l8bugs/go/bugs/digests"
	"github.com/saichler/l8bugs/go/bugs/features"
	"github.com/saichler/l8bugs/go/bugs/projects"
	"github.com/saichler/l8bugs/go/bugs/sprints"
	"github.com/saichler/l8bugs/go/bugs/triage"
	"github.com/saichler/l8bugs/go/bugs/website"
	"github.com/saichler/l8types/go/ifs"
)

func ActivateBugsServices(creds, dbname string, nic ifs.IVNic) {
	website.RegisterTypes(nic.Resources())
	projects.Activate(creds, dbname, nic)
	assignees.Activate(creds, dbname, nic)
	bugs.Activate(creds, dbname, nic)
	features.Activate(creds, dbname, nic)
	sprints.Activate(creds, dbname, nic)
	digests.Activate(creds, dbname, nic)
	triage.Initialize(nic)
}
