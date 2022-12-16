package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"moviehub.com/pkg/discovery"
)

type lServiceName string
type lInstanceID string

// Registry defines an in-memory service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[lServiceName]map[lInstanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service
// registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[lServiceName]map[lInstanceID]*serviceInstance{}}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[lServiceName(serviceName)]; !ok {
		r.serviceAddrs[lServiceName(serviceName)] = map[lInstanceID]*serviceInstance{}
	}

	r.serviceAddrs[lServiceName(serviceName)][lInstanceID(instanceID)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

// Deregister removes a service from the
// registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[lServiceName(serviceName)]; !ok {
		return nil
	}

	delete(r.serviceAddrs[lServiceName(serviceName)], lInstanceID(instanceID))

	return nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Lock()

	if _, ok := r.serviceAddrs[lServiceName(serviceName)]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[lServiceName(serviceName)][lInstanceID(instanceID)]; !ok {
		return errors.New("service instance is not registered yet")
	}
	
	r.serviceAddrs[lServiceName(serviceName)][lInstanceID(instanceID)].lastActive = time.Now()

	return nil
}

// ServiceAddresses returns the list of addresses of
// active instances of a given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[lServiceName(serviceName)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var resp []string
	for _, i := range r.serviceAddrs[lServiceName(serviceName)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		resp = append(resp, i.hostPort)
	}

	return resp, nil
}