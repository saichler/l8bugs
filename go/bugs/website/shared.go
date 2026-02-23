package website

import (
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8bus/go/overlay/vnic"
	"github.com/saichler/l8types/go/ifs"
	"strconv"
)

func CreateVnic(vnet uint32) ifs.IVNic {
	resources := common.CreateResources("web-" + strconv.Itoa(int(vnet)))
	resources.SysConfig().VnetPort = vnet

	RegisterTypes(resources)

	nic := vnic.NewVirtualNetworkInterface(resources, nil)
	nic.Resources().SysConfig().KeepAliveIntervalSeconds = 60
	nic.Start()
	nic.WaitForConnection()

	return nic
}

func RegisterTypes(resources ifs.IResources) {
	registerBugsTypes(resources)
}

func registerBugsTypes(resources ifs.IResources) {
	common.RegisterType[l8bugs.BugsProject, l8bugs.BugsProjectList](resources, "ProjectId")
	common.RegisterType[l8bugs.BugsAssignee, l8bugs.BugsAssigneeList](resources, "AssigneeId")
	common.RegisterType[l8bugs.Bug, l8bugs.BugList](resources, "BugId")
	common.RegisterType[l8bugs.Feature, l8bugs.FeatureList](resources, "FeatureId")
}
