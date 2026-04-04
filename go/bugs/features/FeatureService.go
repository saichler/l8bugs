package features

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Feature"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: ServiceName,
		ServiceArea: ServiceArea,
		PrimaryKey:  "FeatureId",
		Callback:    newFeatureServiceCallback(vnic),
	}, &l8bugs.Feature{}, &l8bugs.FeatureList{}, creds, dbname, vnic)
}

func Features(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return l8common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Feature(featureId string, vnic ifs.IVNic) (*l8bugs.Feature, error) {
	result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.Feature{FeatureId: featureId}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.Feature), nil
}
