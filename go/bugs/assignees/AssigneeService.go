package assignees

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Assignee"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: ServiceName,
		ServiceArea: ServiceArea,
		PrimaryKey:  "AssigneeId",
		Callback:    newAssigneeServiceCallback(vnic),
	}, &l8bugs.BugsAssignee{}, &l8bugs.BugsAssigneeList{}, creds, dbname, vnic)
}

func Assignees(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return l8common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Assignee(assigneeId string, vnic ifs.IVNic) (*l8bugs.BugsAssignee, error) {
	result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsAssignee{AssigneeId: assigneeId}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.BugsAssignee), nil
}
