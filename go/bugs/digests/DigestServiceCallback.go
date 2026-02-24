package digests

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newDigestServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.BugsDigest]("BugsDigest",
		func(e *l8bugs.BugsDigest) { common.GenerateID(&e.DigestId) }).
		Require(func(e *l8bugs.BugsDigest) string { return e.DigestId }, "DigestId").
		Require(func(e *l8bugs.BugsDigest) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.BugsDigest) string { return e.Summary }, "Summary").
		Build()
}
