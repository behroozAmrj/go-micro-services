package main

import (
	controller "mail-service/src/api/Controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	//"mail-service/src/api/controllers"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//brokerController := brokerCTRL.BrokerController{}

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/send", controller.SendMail)
	//mux.Post("/",brokerCTRL.BrokersHandler)

	return mux
}
