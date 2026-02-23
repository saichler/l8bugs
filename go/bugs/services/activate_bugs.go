package services

import (
	"github.com/saichler/l8bugs/go/bugs/assignees"
	"github.com/saichler/l8bugs/go/bugs/bugs"
	"github.com/saichler/l8bugs/go/bugs/features"
	"github.com/saichler/l8bugs/go/bugs/projects"
	"github.com/saichler/l8types/go/ifs"
)

func ActivateBugsServices(creds, dbname string, nic ifs.IVNic) {
	projects.Activate(creds, dbname, nic)
	assignees.Activate(creds, dbname, nic)
	bugs.Activate(creds, dbname, nic)
	features.Activate(creds, dbname, nic)
}
