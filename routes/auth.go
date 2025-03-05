package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/discovery"
	"github.com/mauriciomartinezc/real-estate-api-gateway/middlewares"
	commonDiscovery "github.com/mauriciomartinezc/real-estate-mc-common/discovery"
	"net/http"
)

func AuthRoutes(router *mux.Router, discoveryClient commonDiscovery.DiscoveryClient) {
	auth(router, discoveryClient)
	profile(router, discoveryClient)
	company(router, discoveryClient)
}

func auth(router *mux.Router, discoveryClient commonDiscovery.DiscoveryClient) {
	authRoutes := router.PathPrefix("/auth").Subrouter()

	authRoutes.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/auth/register")
	}).Methods("POST")

	authRoutes.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/auth/login")
	}).Methods("POST")

	authRoutes.Use(middlewares.AuthMiddleware)
	authRoutes.HandleFunc("/resetPassword", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/resetPassword")
	}).Methods("POST")
}

func profile(router *mux.Router, discoveryClient commonDiscovery.DiscoveryClient) {
	authRoutes := router.PathPrefix("/auth").Subrouter()
	authRoutes.Use(middlewares.AuthMiddleware)

	authRoutes.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/profiles")
	}).Methods("POST")

	authRoutes.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/profiles")
	}).Methods("GET")

	authRoutes.HandleFunc("/profiles/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, fmt.Sprintf("/api/profiles/%s", vars["uuid"]))
	}).Methods("PUT")

}

func company(router *mux.Router, discoveryClient commonDiscovery.DiscoveryClient) {
	companyRoutes := router.PathPrefix("/companies").Subrouter()
	companyRoutes.Use(middlewares.AuthMiddleware)

	companyRoutes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/companies")
	}).Methods("POST")

	companyRoutes.HandleFunc("/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, fmt.Sprintf("/api/companies/%s", vars["uuid"]))
	}).Methods("GET")

	companyRoutes.HandleFunc("/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, fmt.Sprintf("/api/companies/%s", vars["uuid"]))
	}).Methods("PUT")

	companyRoutes.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McAuth, "/api/companies/me")
	}).Methods("GET")
}
