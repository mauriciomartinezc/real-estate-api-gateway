package main

import (
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/middlewares"
	"github.com/mauriciomartinezc/real-estate-api-gateway/routes"
	"github.com/mauriciomartinezc/real-estate-api-gateway/utils"
	"github.com/mauriciomartinezc/real-estate-mc-common/config"
	"github.com/mauriciomartinezc/real-estate-mc-common/discovery"
	"github.com/mauriciomartinezc/real-estate-mc-common/discovery/consul"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	port := getServerPort()

	router := initializeRouter()

	discoveryClient := initializeDiscoveryClient()

	// Initialize API routes
	routes.InitRoutes(router, discoveryClient)

	log.Printf("API Gateway running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Println("SERVER_PORT not set, defaulting to 8080")
		port = "8080"
	}
	return port
}

func initializeRouter() *mux.Router {
	router := mux.NewRouter()

	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil || rateLimit <= 0 {
		log.Println("Invalid RATE_LIMIT, defaulting to 100")
		rateLimit = 100
	}

	// Attach middlewares
	router.Use(middlewares.SecurityHeadersMiddleware)
	router.Use(middlewares.RateLimiterMiddleware(rateLimit, 100))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteResponse(w, http.StatusOK, true, "Success", nil)
	}).Methods("GET")

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteResponse(w, http.StatusOK, true, "Success", nil)
	}).Methods("GET")

	return router
}

func initializeDiscoveryClient() discovery.DiscoveryClient {
	return consul.NewConsultApi()
}
