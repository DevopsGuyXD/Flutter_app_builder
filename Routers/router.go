package router

import (
	controller "github.com/DevopsGuyXD/Bizapp/Controllers"
	"github.com/gorilla/mux"
)

// ====================== Routes ======================
func RouterCollection() *mux.Router{

	router := mux.NewRouter()

	router.HandleFunc("/", controller.HomeHandler).Methods("GET")
	router.HandleFunc("/ping", controller.HealthCheck).Methods("GET")
	router.HandleFunc("/api/bizapp/create", controller.CreateApp).Methods("POST")

	return router
}