package projects

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Project"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: ServiceName,
		ServiceArea: ServiceArea,
		PrimaryKey:  "ProjectId",
		Callback:    newProjectServiceCallback(vnic),
	}, &l8bugs.BugsProject{}, &l8bugs.BugsProjectList{}, creds, dbname, vnic)
}

func Projects(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return l8common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Project(projectId string, vnic ifs.IVNic) (*l8bugs.BugsProject, error) {
	result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsProject{ProjectId: projectId}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.BugsProject), nil
}
