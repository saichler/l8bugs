package sprints

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Sprint"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: ServiceName,
		ServiceArea: ServiceArea,
		PrimaryKey:  "SprintId",
		Callback:    newSprintServiceCallback(vnic),
	}, &l8bugs.BugsSprint{}, &l8bugs.BugsSprintList{}, creds, dbname, vnic)
}

func Sprints(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return l8common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Sprint(sprintId string, vnic ifs.IVNic) (*l8bugs.BugsSprint, error) {
	result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsSprint{SprintId: sprintId}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.BugsSprint), nil
}
