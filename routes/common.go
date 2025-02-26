package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/discovery"
	commonDiscovery "github.com/mauriciomartinezc/real-estate-mc-common/discovery"
	"net/http"
)

func CommonRoutes(router *mux.Router, discoveryClient commonDiscovery.DiscoveryClient) {
	router.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		discovery.HandleProxyRequest(w, r, discoveryClient, McCommon, "/api/countries")
	}).Methods("GET")

	router.HandleFunc("/states/{countryUuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		discovery.HandleProxyRequest(w, r, discoveryClient, McCommon, fmt.Sprintf("/api/states/%s", vars["countryUuid"]))
	}).Methods("GET")

	router.HandleFunc("/cities/{stateUuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		discovery.HandleProxyRequest(w, r, discoveryClient, McCommon, fmt.Sprintf("/api/cities/%s", vars["stateUuid"]))
	}).Methods("GET")
}
