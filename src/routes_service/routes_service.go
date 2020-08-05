package routes_service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/palantir/stacktrace"
	"handlers"
	"log"
	"net/http"
	"services"
	"sync"
)

type RoutesService struct {
	ConfSrv *services.ConfService
}

var routesServiceInstance *RoutesService
var routesServiceOnce sync.Once

func NewRoutesService() *RoutesService {
	routesServiceOnce.Do(func() {
		routesServiceInstance = &RoutesService{
			ConfSrv: services.NewConfService(),
		}
	})
	return routesServiceInstance
}

func (rs *RoutesService) InitRoutes() error {
	log.Println("Configuring routes")
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1/").Subrouter()
	subRouter.HandleFunc("/restaurants", handlers.GetAllRestaurants).Methods("GET")
	subRouter.HandleFunc("/restaurant/{name}", handlers.EditRestaurant).Methods("PUT")
	subRouter.HandleFunc("/restaurant/{name}", handlers.DeleteRestaurant).Methods("DELETE")
	http.Handle("/", router)

	log.Println(fmt.Sprintf("Server is listening in the port %d...", rs.ConfSrv.Config.ServicePort))
	err := http.ListenAndServe(fmt.Sprintf(":%d", rs.ConfSrv.Config.ServicePort), nil)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to init the routes")
	}
	return nil
}
