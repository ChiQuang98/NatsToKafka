package routers

import (
	"NatsToKafka/controllers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func SetMcuRouter(router *mux.Router) *mux.Router {
	router.Handle("/anomalyconfig",
		negroni.New(
			negroni.HandlerFunc(controllers.AnomalyConfig),
		)).Methods("POST")
	return router
}

