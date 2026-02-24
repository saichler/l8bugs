package digests

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Digest"
	ServiceArea = byte(20)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService[l8bugs.BugsDigest, l8bugs.BugsDigestList](common.ServiceConfig{
		ServiceName:   ServiceName,
		ServiceArea:   ServiceArea,
		PrimaryKey:    "DigestId",
		Callback:      newDigestServiceCallback(),
		Transactional: true,
	}, creds, dbname, vnic)
}

func Digests(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Digest(digestId string, vnic ifs.IVNic) (*l8bugs.BugsDigest, error) {
	return common.GetEntity(ServiceName, ServiceArea, &l8bugs.BugsDigest{DigestId: digestId}, vnic)
}
