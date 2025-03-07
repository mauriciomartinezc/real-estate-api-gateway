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

	profiles := protected.PathPrefix("/profiles").Subrouter()
	profiles.HandleFunc("/", discovery.ProxyHandler(lb)).Methods("POST")
	profiles.HandleFunc("/", discovery.ProxyHandler(lb)).Methods("GET")
	profiles.HandleFunc("/{uuid}", discovery.ProxyHandler(lb)).Methods("PUT")

	companies := protected.PathPrefix("/companies").Subrouter()
	companies.HandleFunc("/", discovery.ProxyHandler(lb)).Methods("POST")
	companies.HandleFunc("/{uuid}", discovery.ProxyHandler(lb)).Methods("GET")
	companies.HandleFunc("/{uuid}", discovery.ProxyHandler(lb)).Methods("PUT")
	companies.HandleFunc("/me", discovery.ProxyHandler(lb)).Methods("GET")
}
