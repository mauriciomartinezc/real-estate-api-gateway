package middlewares

import (
	"log"
	"net/http"
	"os"
)

// Definir valores por defecto en caso de que no estén en él .env
var corsOrigin = getEnv("CORS_ORIGIN", "*")
var corsMethods = getEnv("CORS_METHODS", "GET, POST, PUT, DELETE, OPTIONS")
var corsHeaders = getEnv("CORS_HEADERS", "Content-Type, Authorization, Accept-Language, X-Company-Id")

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Encabezados de seguridad
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Configuración de CORS
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", corsMethods)
		w.Header().Set("Access-Control-Allow-Headers", corsHeaders)

		// Manejo de preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("SecurityHeadersMiddleware: Adding security headers for %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// getEnv obtiene una variable de entorno con un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
