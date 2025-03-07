package services

import "os"

// Constantes que identifican los servicios.
const (
	McCommon = "mc-common"
	McAuth   = "mc-auth"
)

// Service representa un servicio con su nombre y puerto.
type Service struct {
	Name string
	Port string
}

// Services es una lista de Service.
type Services []Service

// GetServices retorna la lista de servicios configurados.
// Si las variables de entorno no están definidas, se usan valores por defecto.
func GetServices() Services {
	commonPort := os.Getenv("MC_COMMON_SERVICE_SERVICE_PORT")
	if commonPort == "" {
		commonPort = "8082"
	}
	authPort := os.Getenv("MC_AUTH_SERVICE_SERVICE_PORT")
	if authPort == "" {
		authPort = "8081"
	}

	return Services{
		{Name: McCommon, Port: commonPort},
		{Name: McAuth, Port: authPort},
	}
}
