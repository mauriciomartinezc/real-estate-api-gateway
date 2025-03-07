package middlewares

import (
	"fmt"
	"github.com/mauriciomartinezc/real-estate-api-gateway/services"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mauriciomartinezc/real-estate-api-gateway/loadbalancer"
	"github.com/mauriciomartinezc/real-estate-api-gateway/utils"
	"github.com/mauriciomartinezc/real-estate-mc-common/config"
)

var (
	cacheClient = config.NewCacheClient()
	authLB      *loadbalancer.DynamicLoadBalancer
	authLBOnce  sync.Once
)

// getAuthLoadBalancer crea (una única vez) el load balancer para el servicio de autenticación.
// Se basa en el DNS de Kubernetes, usando la variable de entorno "AUTH_SERVICE_DNS" o un valor por defecto.
func getAuthLoadBalancer() *loadbalancer.DynamicLoadBalancer {
	authLBOnce.Do(func() {
		authDNS := os.Getenv("AUTH_SERVICE_DNS")
		if authDNS == "" {
			defaultEndpoints := utils.GetDefaultEndpointLb()
			// Valor por defecto basado en Kubernetes
			authDNS = defaultEndpoints[services.McAuth]
		}
		authLB = loadbalancer.NewDynamicLoadBalancer("AUTH_SERVICE_DNS", authDNS)
	})
	return authLB
}

// AuthMiddleware utiliza internamente el load balancer para validar el token.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !ValidateToken(token) {
			utils.WriteResponse(w, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateToken realiza la validación del token haciendo una llamada al endpoint "/auth/validate"
// del servicio de autenticación obtenido vía load balancer.
func ValidateToken(token string) bool {
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var cachedToken string
	if err := cacheClient.Get(token, &cachedToken); err == nil && cachedToken == "valid" {
		return true
	}

	lb := getAuthLoadBalancer()
	authEndpoint := lb.GetNextInstance()
	if authEndpoint == "" {
		fmt.Println("No hay endpoint disponible para la autenticación")
		return false
	}

	validateURL := fmt.Sprintf("%s/api/auth/validate", authEndpoint)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", validateURL, nil)
	if err != nil {
		fmt.Println("Error creando la solicitud de validación:", err)
		return false
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Token inválido o error en la comunicación:", err)
		return false
	}
	defer resp.Body.Close()

	if err := cacheClient.Set(token, "valid", 10*time.Minute); err != nil {
		fmt.Println("Error al almacenar el token en cache:", err)
		return false
	}

	return true
}
