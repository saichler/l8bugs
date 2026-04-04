package digests

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newDigestServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return l8common.NewValidation(&l8bugs.BugsDigest{}, vnic).
		Require(func(e *l8bugs.BugsDigest) string { return e.DigestId }, "DigestId").
		Require(func(e *l8bugs.BugsDigest) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.BugsDigest) string { return e.Summary }, "Summary").
		Build()
}
