package routes

import (
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/discovery"
	"github.com/mauriciomartinezc/real-estate-api-gateway/loadbalancer"
)

func CommonRoutes(router *mux.Router, lb *loadbalancer.DynamicLoadBalancer) {
	router.HandleFunc("/countries", discovery.ProxyHandler(lb)).Methods("GET")
	router.HandleFunc("/states/{countryUuid}", discovery.ProxyHandler(lb)).Methods("GET")
	router.HandleFunc("/cities/{stateUuid}", discovery.ProxyHandler(lb)).Methods("GET")
}
