package common

import (
	_ "github.com/lib/pq"
	"github.com/saichler/l8orm/go/orm/persist"
	"github.com/saichler/l8orm/go/orm/plugins/postgres"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8types/go/types/l8web"
	"github.com/saichler/l8utils/go/utils/web"
	"google.golang.org/protobuf/proto"
)

type ProtoMessage[T any] interface {
	*T
	proto.Message
}

type ServiceConfig struct {
	ServiceName   string
	ServiceArea   byte
	PrimaryKey    string
	Callback      ifs.IServiceCallback
	Transactional bool
}

func ActivateService[T any, TList any, PT ProtoMessage[T], PTL ProtoMessage[TList]](cfg ServiceConfig, creds, dbname string, vnic ifs.IVNic) {
	_, user, pass, _, err := vnic.Resources().Security().Credential(creds, dbname, vnic.Resources())
	if err != nil {
		panic(err)
	}
	db := OpenDBConection(dbname, user, pass)
	p := postgres.NewPostgres(db, vnic.Resources())

	sla := ifs.NewServiceLevelAgreement(&persist.OrmService{}, cfg.ServiceName, cfg.ServiceArea, true, cfg.Callback)
	sla.SetServiceItem(PT(new(T)))
	sla.SetServiceItemList(PTL(new(TList)))
	sla.SetPrimaryKeys(cfg.PrimaryKey)
	sla.SetArgs(p)

	if cfg.Transactional {
		sla.SetTransactional(true)
		sla.SetReplication(true)
		sla.SetReplicationCount(3)
	}

	ws := web.New(cfg.ServiceName, cfg.ServiceArea, 0)
	ws.AddEndpoint(PT(new(T)), ifs.POST, &l8web.L8Empty{})
	ws.AddEndpoint(PTL(new(TList)), ifs.POST, &l8web.L8Empty{})
	ws.AddEndpoint(PT(new(T)), ifs.PUT, &l8web.L8Empty{})
	ws.AddEndpoint(PT(new(T)), ifs.PATCH, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Query{}, ifs.DELETE, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Query{}, ifs.GET, PTL(new(TList)))
	sla.SetWebService(ws)

	vnic.Resources().Services().Activate(sla, vnic)
}

func ServiceHandler(serviceName string, serviceArea byte, vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return vnic.Resources().Services().ServiceHandler(serviceName, serviceArea)
}

func GetEntity[T any](serviceName string, serviceArea byte, filter *T, vnic ifs.IVNic) (*T, error) {
	handler, ok := ServiceHandler(serviceName, serviceArea, vnic)
	if ok {
		resp := handler.Get(object.New(nil, filter), vnic)
		if resp.Error() != nil {
			return nil, resp.Error()
		}
		if resp.Element() != nil {
			return resp.Element().(*T), nil
		}
		return nil, nil
	}
	resp := vnic.Request("", serviceName, serviceArea, ifs.GET, filter, 30)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	if resp.Element() != nil {
		return resp.Element().(*T), nil
	}
	return nil, nil
}

func GetEntities[T any](serviceName string, serviceArea byte, filter *T, vnic ifs.IVNic) ([]*T, error) {
	handler, ok := ServiceHandler(serviceName, serviceArea, vnic)
	if ok {
		resp := handler.Get(object.New(nil, filter), vnic)
		if resp.Error() != nil {
			return nil, resp.Error()
		}
		return extractElements[T](resp.Elements()), nil
	}
	resp := vnic.Request("", serviceName, serviceArea, ifs.GET, filter, 30)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return extractElements[T](resp.Elements()), nil
}

func extractElements[T any](elems []interface{}) []*T {
	result := make([]*T, 0, len(elems))
	for _, e := range elems {
		if t, ok := e.(*T); ok {
			result = append(result, t)
		}
	}
	return result
}

func PutEntity[T any](serviceName string, serviceArea byte, entity *T, vnic ifs.IVNic) error {
	handler, ok := ServiceHandler(serviceName, serviceArea, vnic)
	if ok {
		resp := handler.Put(object.New(nil, entity), vnic)
		if resp.Error() != nil {
			return resp.Error()
		}
		return nil
	}
	resp := vnic.Request("", serviceName, serviceArea, ifs.PUT, entity, 30)
	if resp.Error() != nil {
		return resp.Error()
	}
	return nil
}

func PostEntity[T any](serviceName string, serviceArea byte, entity *T, vnic ifs.IVNic) (*T, error) {
	handler, ok := ServiceHandler(serviceName, serviceArea, vnic)
	if ok {
		resp := handler.Post(object.New(nil, entity), vnic)
		if resp.Error() != nil {
			return nil, resp.Error()
		}
		if resp.Element() != nil {
			return resp.Element().(*T), nil
		}
		return entity, nil
	}
	resp := vnic.Request("", serviceName, serviceArea, ifs.POST, entity, 30)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	if resp.Element() != nil {
		return resp.Element().(*T), nil
	}
	return entity, nil
}

func QueryEntities[T any](serviceName string, serviceArea byte, queryText string, vnic ifs.IVNic) ([]*T, error) {
	query := &l8api.L8Query{Text: queryText}
	resp := vnic.Request("", serviceName, serviceArea, ifs.GET, query, 30)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return extractElements[T](resp.Elements()), nil
}

func EntityExists[T any](serviceName string, serviceArea byte, filter *T, vnic ifs.IVNic) (bool, error) {
	existing, err := GetEntities[T](serviceName, serviceArea, filter, vnic)
	if err != nil {
		return false, err
	}
	return len(existing) > 0, nil
}
