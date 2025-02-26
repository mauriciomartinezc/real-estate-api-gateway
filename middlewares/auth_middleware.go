package middlewares

import (
	"fmt"
	"github.com/mauriciomartinezc/real-estate-api-gateway/utils"
	"github.com/mauriciomartinezc/real-estate-mc-common/config"
	"net/http"
	"os"
	"time"
)

var (
	authServiceURL = os.Getenv("AUTH_SERVICE_URL")
	cacheClient    = config.NewCacheClient()
)

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

// ValidateToken envía el token al microservice de autenticación para su validación
func ValidateToken(token string) bool {
	// Remover el prefijo "Bearer" si está presente
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var cachedToken string
	if err := cacheClient.Get(token, &cachedToken); err == nil && cachedToken == "valid" {
		return true
	}

	// Si no está en caché, validar con el servicio de autenticación
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/validate", authServiceURL), nil)
	if err != nil {
		fmt.Println("Error en la solicitud de validación:", err)
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Token inválido:", err)
		return false
	}
	defer resp.Body.Close()

	// Guardar token válido en caché con tiempo de expiración adecuado
	err = cacheClient.Set(token, "valid", 10*time.Minute)
	if err != nil {
		fmt.Println("Error al almacenar el token en cache:", err)
		return false
	}

	return true
}
