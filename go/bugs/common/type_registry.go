package common

import (
	"github.com/saichler/l8types/go/ifs"
)

func RegisterType[T any, TList any](resources ifs.IResources, pkField string) {
	resources.Introspector().Decorators().AddPrimaryKeyDecorator(new(T), pkField)
	resources.Registry().Register(new(TList))
}
