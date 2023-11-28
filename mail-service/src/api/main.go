package main

import (
	"fmt"
	"log"
	model "mail-service/src/models"
	"net/http"
)

type Config struct {
	Mailer model.Mail
}

const webPort = "80"

func main() {
	app := Config{}

	log.Println("Starting mail service on port: ", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
