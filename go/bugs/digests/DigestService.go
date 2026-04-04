package digests

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Digest"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: ServiceName,
		ServiceArea: ServiceArea,
		PrimaryKey:  "DigestId",
		Callback:    newDigestServiceCallback(vnic),
	}, &l8bugs.BugsDigest{}, &l8bugs.BugsDigestList{}, creds, dbname, vnic)
}

func Digests(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return l8common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Digest(digestId string, vnic ifs.IVNic) (*l8bugs.BugsDigest, error) {
	result, err := l8common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsDigest{DigestId: digestId}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.BugsDigest), nil
}
