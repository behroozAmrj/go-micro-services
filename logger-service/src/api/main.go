package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/src/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort    = "80"
	rpcPort    = "5001"
	mongodbURL = "mongodb://mongo:27017"
	gRpcPort   = "50001"
)

var client *mongo.Client

type Config struct {
	Model data.Model
}

func main() {
	// connect to Mongo
	mongoClient, err := connectoToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient
	
	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close connection

	defer func()  {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()	

	app := Config{
		Model: data.New(client),
	}
	log.Printf("starting service on port" , webPort)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		
	}
}



func connectoToMongo() (*mongo.Client, error) {
	// create connectio options
	clientOptions := options.Client().ApplyURI(mongodbURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err :=  mongo.Connect(context.TODO() , clientOptions)
	if err != nil {
		log.Println("Error connecting to Mongo" , err)
		return nil, err
	}

	log.Println("Connecting to Mongo!")
	return c , nil

}