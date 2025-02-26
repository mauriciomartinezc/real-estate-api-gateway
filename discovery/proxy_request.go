package discovery

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mauriciomartinezc/real-estate-mc-common/discovery"
)

// HandleProxyRequest redirects requests to the specified service via Consul
func HandleProxyRequest(w http.ResponseWriter, r *http.Request, discoveryClient discovery.DiscoveryClient, serviceName string, path string) {
	// Get the service address from Consul
	address, err := discoveryClient.GetServiceAddress(serviceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving service address %s: %v", serviceName, err), http.StatusInternalServerError)
		return
	}

	// Construct the target URL for the microservice endpoint
	url := fmt.Sprintf("%s%s", address, path)

	// Create a new request for the microservice
	proxyReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request to service %s: %v", serviceName, err), http.StatusInternalServerError)
		return
	}

	// Copy relevant headers from the original request
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Perform the request to the microservice
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error communicating with %s: %v", serviceName, err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy relevant headers from the microservice response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Propagate the HTTP status code to the client
	w.WriteHeader(resp.StatusCode)

	// Write the response body to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying the response body", http.StatusInternalServerError)
		return
	}
}
