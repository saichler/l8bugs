package features

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Feature"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService[l8bugs.Feature, l8bugs.FeatureList](common.ServiceConfig{
		ServiceName:   ServiceName,
		ServiceArea:   ServiceArea,
		PrimaryKey:    "FeatureId",
		Callback:      newFeatureServiceCallback(),
		Transactional: true,
	}, creds, dbname, vnic)
}

func Features(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Feature(featureId string, vnic ifs.IVNic) (*l8bugs.Feature, error) {
	return common.GetEntity(ServiceName, ServiceArea, &l8bugs.Feature{FeatureId: featureId}, vnic)
}
