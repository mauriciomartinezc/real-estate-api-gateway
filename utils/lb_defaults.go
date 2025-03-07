package utils

import (
	"fmt"
	"sync"

	"github.com/mauriciomartinezc/real-estate-api-gateway/services"
)

var (
	defaultEndpointsCache map[string]string
	endpointsOnce         sync.Once
)

// GetDefaultEndpointLb construye un mapa de endpoints por servicio basado en la lista definida en services.GetServices()
// y lo cachea para evitar recalcularlo en cada llamada.
func GetDefaultEndpointLb() map[string]string {
	endpointsOnce.Do(func() {
		defaultEndpointsCache = make(map[string]string)
		for _, service := range services.GetServices() {
			// Se asume que el DNS de Kubernetes sigue el formato: http://<service>-service.default.svc.cluster.local:<port>
			defaultEndpointsCache[service.Name] = fmt.Sprintf("http://%s-service.default.svc.cluster.local:%s", service.Name, service.Port)
		}
	})
	return defaultEndpointsCache
}
