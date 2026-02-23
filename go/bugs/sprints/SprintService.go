package sprints

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Sprint"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService[l8bugs.BugsSprint, l8bugs.BugsSprintList](common.ServiceConfig{
		ServiceName:   ServiceName,
		ServiceArea:   ServiceArea,
		PrimaryKey:    "SprintId",
		Callback:      newSprintServiceCallback(),
		Transactional: true,
	}, creds, dbname, vnic)
}

func Sprints(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Sprint(sprintId string, vnic ifs.IVNic) (*l8bugs.BugsSprint, error) {
	return common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsSprint{SprintId: sprintId}, vnic)
}
