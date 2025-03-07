package main

import (
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/loadbalancer"
	"github.com/mauriciomartinezc/real-estate-api-gateway/middlewares"
	"github.com/mauriciomartinezc/real-estate-api-gateway/routes"
	"github.com/mauriciomartinezc/real-estate-api-gateway/services"
	"github.com/mauriciomartinezc/real-estate-api-gateway/utils"
	"github.com/mauriciomartinezc/real-estate-mc-common/config"
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

	loadBalancers := initLoadBalancers()

	// Initialize API routes
	routes.InitRoutes(router, loadBalancers)

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

	return router
}

func initLoadBalancers() map[string]*loadbalancer.DynamicLoadBalancer {
	defaultEndpoints := utils.GetDefaultEndpointLb()

	commonLb := loadbalancer.NewDynamicLoadBalancer("COMMON_SERVICE_DNS", defaultEndpoints[services.McCommon])
	authLB := loadbalancer.NewDynamicLoadBalancer("AUTH_SERVICE_DNS", defaultEndpoints[services.McAuth])

	loadBalancers := map[string]*loadbalancer.DynamicLoadBalancer{
		services.McCommon: commonLb,
		services.McAuth:   authLB,
	}

	return loadBalancers
}
