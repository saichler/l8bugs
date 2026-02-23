package assignees

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Assignee"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService[l8bugs.BugsAssignee, l8bugs.BugsAssigneeList](common.ServiceConfig{
		ServiceName:   ServiceName,
		ServiceArea:   ServiceArea,
		PrimaryKey:    "AssigneeId",
		Callback:      newAssigneeServiceCallback(),
		Transactional: true,
	}, creds, dbname, vnic)
}

func Assignees(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Assignee(assigneeId string, vnic ifs.IVNic) (*l8bugs.BugsAssignee, error) {
	return common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsAssignee{AssigneeId: assigneeId}, vnic)
}
