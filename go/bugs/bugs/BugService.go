package bugs

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Bug"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: ServiceName,
		ServiceArea: ServiceArea,
		PrimaryKey:  "BugId",
		Callback:    newBugServiceCallback(vnic),
	}, &l8bugs.Bug{}, &l8bugs.BugList{}, creds, dbname, vnic)
}

func Bugs(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return l8common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Bug(bugId string, vnic ifs.IVNic) (*l8bugs.Bug, error) {
	result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.Bug{BugId: bugId}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.Bug), nil
}
