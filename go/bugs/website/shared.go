package website

import (
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8bus/go/overlay/vnic"
	"github.com/saichler/l8types/go/ifs"
	"strconv"
)

func CreateVnic(vnet uint32) ifs.IVNic {
	resources := l8common.CreateResources("web-"+strconv.Itoa(int(vnet)), "/data/logs/l8bugs", vnet)

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
	l8common.RegisterType(resources, &l8bugs.BugsProject{}, &l8bugs.BugsProjectList{}, "ProjectId")
	l8common.RegisterType(resources, &l8bugs.BugsAssignee{}, &l8bugs.BugsAssigneeList{}, "AssigneeId")
	l8common.RegisterType(resources, &l8bugs.Bug{}, &l8bugs.BugList{}, "BugId")
	l8common.RegisterType(resources, &l8bugs.Feature{}, &l8bugs.FeatureList{}, "FeatureId")
	l8common.RegisterType(resources, &l8bugs.BugsSprint{}, &l8bugs.BugsSprintList{}, "SprintId")
	l8common.RegisterType(resources, &l8bugs.BugsDigest{}, &l8bugs.BugsDigestList{}, "DigestId")
}
