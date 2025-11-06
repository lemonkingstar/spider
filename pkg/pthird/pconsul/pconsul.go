package pconsul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ServiceRegistry interface {
	ServiceDiscovery
	Register(ServiceInstance, bool) error
	Deregister() error
	DeregisterID(string) error
}

type registry struct {
	client	*api.Client
	localID 	string
	localName	string
}

func (r *registry) Register(s ServiceInstance, check bool) error {
	registration := &api.AgentServiceRegistration{
		ID: s.GetInstanceID(),
		Name: s.GetServiceName(),
		Address: s.GetAddress(),
		Port: s.GetPort(),
		Meta: s.GetMeta(),
		Tags: s.GetTags(),
	}

	if check {
		asc := &api.AgentServiceCheck{
			HTTP: fmt.Sprintf("http://%s:%d/healthz", s.GetAddress(), s.GetPort()),
			Timeout: "10s",
			Interval: "10s",
			// 自动注销不健康的服务节点
			DeregisterCriticalServiceAfter: "1m",
		}
		registration.Check = asc
	}

	r.localID = s.GetInstanceID()
	r.localName = s.GetServiceName()
	return r.client.Agent().ServiceRegister(registration)
}

func (r *registry) Deregister() error {
	if r.localID == "" { return nil }

	r.localID = ""
	return r.client.Agent().ServiceDeregister(r.localID)
}

func (r *registry) DeregisterID(instanceID string) error {
	if instanceID == "" { return nil }

	return r.client.Agent().ServiceDeregister(instanceID)
}

func NewServiceRegistry(address string) (ServiceRegistry, error) {
	config := api.DefaultConfig()
	config.Address = address
	//config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &registry{client: client}, nil
}
