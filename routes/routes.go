package routes

import (
	"github.com/gorilla/mux"
	"github.com/mauriciomartinezc/real-estate-api-gateway/utils"
	commonDiscovery "github.com/mauriciomartinezc/real-estate-mc-common/discovery"
	"log"
	"net/http"
)

func InitRoutes(router *mux.Router, discoveryClient commonDiscovery.DiscoveryClient) {
	// Asegurar que todas las rutas de API tengan el prefijo "/api"
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Definir rutas
	CommonRoutes(apiRouter, discoveryClient)
	AuthRoutes(apiRouter, discoveryClient)

	// Establecer manejador 404 solo después de definir todas las rutas
	router.NotFoundHandler = http.HandlerFunc(custom404Handler)

	// Mostrar rutas registradas en logs
	routerWalk(router)
}

// Manejador para rutas no encontradas
func custom404Handler(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, http.StatusNotFound, false, "Not Found", nil)
}

// Mostrar rutas en logs para depuración
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
