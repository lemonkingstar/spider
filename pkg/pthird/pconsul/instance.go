package pconsul

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type ServiceInstance interface {
	GetInstanceID() string
	GetServiceName() string
	GetAddress() string
	GetPort() int
	GetMeta() map[string]string
	GetTags() []string
}

type instance struct {
	id 			string
	name 		string
	address		string
	port		int
	meta   		map[string]string
	tags        []string
}

func NewServiceInstance(
	instanceID string,
	serviceName string,
	address string,
	port int,
	meta map[string]string,
	tags []string) (ServiceInstance, error) {
	if serviceName == "" {
		return nil, errors.New("service is none")
	}
	if address == "" || port == 0 {
		return nil, errors.New("address is none")
	}
	if instanceID == "" {
		instanceID = fmt.Sprintf("%s-%d-%d", serviceName, time.Now().Unix(), rand.Intn(9000)+1000)
	}

	return &instance{
		id: instanceID, name: serviceName,
		address: address, port: port,
		meta: meta,
		tags: tags,
	}, nil
}

func (si *instance) GetInstanceID() string {
	return si.id
}

func (si *instance) GetServiceName() string {
	return si.name
}

func (si *instance) GetAddress() string {
	return si.address
}

func (si *instance) GetPort() int {
	return si.port
}

func (si *instance) GetMeta() map[string]string {
	return si.meta
}

func (si *instance) GetTags() []string {
	return si.tags
}
