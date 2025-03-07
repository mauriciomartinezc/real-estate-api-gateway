package discovery

import (
	"fmt"
	"github.com/mauriciomartinezc/real-estate-api-gateway/loadbalancer"
	"io"
	"log"
	"net/http"
)

func ProxyHandler(lb *loadbalancer.DynamicLoadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleProxyRequest(w, r, lb)
	}
}

// HandleProxyRequest redirects requests to the specified service via Consul
func handleProxyRequest(w http.ResponseWriter, r *http.Request, lb *loadbalancer.DynamicLoadBalancer) {
	target := lb.GetNextInstance()
	if target == "" {
		http.Error(w, "No hay un endpoint disponible para el servicio", http.StatusServiceUnavailable)
		return
	}

	// Se asume que se usa la misma RequestURI para mantener la ruta y parámetros.
	url := target + r.RequestURI

	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creando la solicitud proxy: %v", err), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error en la comunicación con el servicio destino: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copiar headers de la respuesta
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copiando el cuerpo de la respuesta: %v", err)
	}
}
