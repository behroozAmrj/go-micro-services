package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Application struct{}

func main() {
	app := &Application{}
	log.Printf("Start broker service on port %s" , webPort)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}