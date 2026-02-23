package projects

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Project"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService[l8bugs.BugsProject, l8bugs.BugsProjectList](common.ServiceConfig{
		ServiceName:   ServiceName,
		ServiceArea:   ServiceArea,
		PrimaryKey:    "ProjectId",
		Callback:      newProjectServiceCallback(),
		Transactional: true,
	}, creds, dbname, vnic)
}

func Projects(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Project(projectId string, vnic ifs.IVNic) (*l8bugs.BugsProject, error) {
	return common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsProject{ProjectId: projectId}, vnic)
}
