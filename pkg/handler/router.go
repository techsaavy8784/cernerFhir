package handler

import (
	//"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	//fmt.Printf("ROUTES: %v\n", routes)
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	// Add another route to serve static files for the Swagger-UI
	router.PathPrefix("/docs").Handler(http.StripPrefix("/docs", http.FileServer(http.Dir("./swagger-ui-2.2.5/dist"))))

	return router
}
