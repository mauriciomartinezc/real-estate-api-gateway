package routes

import (
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/discovery"
	"github.com/mauriciomartinezc/real-estate-api-gateway/loadbalancer"
	"github.com/mauriciomartinezc/real-estate-api-gateway/middlewares"
)

// AuthRoutes configura las rutas públicas y protegidas.
func AuthRoutes(router *mux.Router, lb *loadbalancer.DynamicLoadBalancer) {
	// Rutas públicas de autenticación
	authPublic := router.PathPrefix("/auth").Subrouter()
	authPublic.HandleFunc("/register", discovery.ProxyHandler(lb)).Methods("POST")
	authPublic.HandleFunc("/login", discovery.ProxyHandler(lb)).Methods("POST")

	// Rutas protegidas (se aplica el middleware de autenticación)
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middlewares.AuthMiddleware)

	protected.HandleFunc("/auth/resetPassword", discovery.ProxyHandler(lb)).Methods("POST")

	protected.HandleFunc("/profiles", discovery.ProxyHandler(lb)).Methods("POST")
	protected.HandleFunc("/profiles", discovery.ProxyHandler(lb)).Methods("GET")
	protected.HandleFunc("/{uuid}", discovery.ProxyHandler(lb)).Methods("PUT")

	protected.HandleFunc("/companies", discovery.ProxyHandler(lb)).Methods("POST")
	protected.HandleFunc("/companies/{uuid}", discovery.ProxyHandler(lb)).Methods("GET")
	protected.HandleFunc("/companies/{uuid}", discovery.ProxyHandler(lb)).Methods("PUT")
	protected.HandleFunc("/companies/me", discovery.ProxyHandler(lb)).Methods("GET")
}
