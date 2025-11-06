package pconsul

import "unsafe"

type ServiceDiscovery interface {
	GetInstances() ([]ServiceInstance, error)
	GetInstancesWithName(string) ([]ServiceInstance, error)
	GetServices() ([]string, error)
}

func (r *registry) GetInstances() ([]ServiceInstance, error) {
	return r.GetInstancesWithName(r.localName)
}

func (r *registry) GetInstancesWithName(serviceName string) ([]ServiceInstance, error) {
	catalogService, _, err := r.client.Catalog().Service(serviceName, "", nil)
	if err != nil { return nil, err }

	result := make([]ServiceInstance, 0, len(catalogService))
	if len(catalogService) > 0 {
		for _, s := range catalogService {
			result = append(result, &instance{
				id: s.ServiceID,
				name: s.ServiceName,
				address: s.Address,
				port: s.ServicePort,
				meta: s.ServiceMeta,
				tags: s.ServiceTags,
			})
		}
	}
	return result, nil
}

func (r *registry) GetServices() ([]string, error) {
	services, _, err := r.client.Catalog().Services(nil)
	if err != nil { return nil, err }

	result := make([]string, unsafe.Sizeof(services))
	index := 0
	for serviceName, _ := range services {
		result[index] = serviceName
		index++
	}
	return result, nil
}
