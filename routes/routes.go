package routes

import (
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/loadbalancer"
	"github.com/mauriciomartinezc/real-estate-api-gateway/services"
	"github.com/mauriciomartinezc/real-estate-api-gateway/utils"
	"log"
	"net/http"
)

func InitRoutes(router *mux.Router, loadBalancers map[string]*loadbalancer.DynamicLoadBalancer) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	CommonRoutes(apiRouter, loadBalancers[services.McCommon])
	AuthRoutes(apiRouter, loadBalancers[services.McAuth])

	router.NotFoundHandler = http.HandlerFunc(custom404Handler)

	routerWalk(router)
}

func custom404Handler(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, http.StatusNotFound, false, "Not Found", nil)
}

func routerWalk(router *mux.Router) {
	log.Println("Registered routes:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Printf("Route: %s", path)
		}
		return nil
	})
}
