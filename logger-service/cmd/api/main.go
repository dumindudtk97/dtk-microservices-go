package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in orde to disconnect (handling disconnect is a mongo requirement)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}() // () syntax to use as creating fnc

	app := Config{
		Models: data.New(client),
	}

	// start web server
	log.Printf("logger server starting on %s", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}

// see doc of driver to see how to connect to mongo in detail
func connectToMongo() (*mongo.Client, error) {

	// create connection options // better to give them in other ways like env, cmdline options etc. than specify pass,username here
	// can do that in docker as secets and pass them
	clientOptins := options.Client().ApplyURI(mongoURL)
	clientOptins.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptins)
	if err != nil {
		log.Println("Error connecting to MongoDB", err)
		return nil, err
	}

	return c, nil
}
