package bugs

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Bug"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService[l8bugs.Bug, l8bugs.BugList](common.ServiceConfig{
		ServiceName:   ServiceName,
		ServiceArea:   ServiceArea,
		PrimaryKey:    "BugId",
		Callback:      newBugServiceCallback(),
		Transactional: true,
	}, creds, dbname, vnic)
}

func Bugs(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Bug(bugId string, vnic ifs.IVNic) (*l8bugs.Bug, error) {
	return common.GetEntity(ServiceName, ServiceArea, &l8bugs.Bug{BugId: bugId}, vnic)
}
