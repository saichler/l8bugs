package features

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
)

func newFeatureServiceCallback() ifs.IServiceCallback {
	return common.NewValidation[l8bugs.Feature]("Feature",
		func(e *l8bugs.Feature) { common.GenerateID(&e.FeatureId) }).
		Require(func(e *l8bugs.Feature) string { return e.FeatureId }, "FeatureId").
		Require(func(e *l8bugs.Feature) string { return e.ProjectId }, "ProjectId").
		Require(func(e *l8bugs.Feature) string { return e.Title }, "Title").
		Build()
}
