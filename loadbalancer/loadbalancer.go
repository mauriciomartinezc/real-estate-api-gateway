package loadbalancer

import (
	"os"
	"strings"
	"sync"
)

// DynamicLoadBalancer administra una lista de endpoints para un servicio.
// Si se configura la variable de entorno (por ejemplo, AUTH_SERVICE_DNS) se usará una lista
// de endpoints separados por comas; de lo contrario se usará el endpoint por defecto.
type DynamicLoadBalancer struct {
	endpoints    []string
	currentIndex int
	mutex        sync.Mutex
}

// NewDynamicLoadBalancer crea un load balancer para el servicio.
// serviceEnvVar es el nombre de la variable de entorno que contiene los endpoints (opcional).
// defaultEndpoint es el valor por defecto (p.ej. "http://auth-service.default.svc.cluster.local").
func NewDynamicLoadBalancer(serviceEnvVar string, defaultEndpoint string) *DynamicLoadBalancer {
	endpointsStr := os.Getenv(serviceEnvVar)
	var endpoints []string
	if endpointsStr != "" {
		// Permite configurar varios endpoints separados por comas
		split := strings.Split(endpointsStr, ",")
		for _, e := range split {
			trimmed := strings.TrimSpace(e)
			if trimmed != "" {
				endpoints = append(endpoints, trimmed)
			}
		}
	} else {
		endpoints = []string{defaultEndpoint}
	}

	return &DynamicLoadBalancer{
		endpoints: endpoints,
	}
}

// GetNextInstance retorna el siguiente endpoint en round-robin.
func (lb *DynamicLoadBalancer) GetNextInstance() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.endpoints) == 0 {
		return ""
	}

	endpoint := lb.endpoints[lb.currentIndex]
	lb.currentIndex = (lb.currentIndex + 1) % len(lb.endpoints)
	return endpoint
}
