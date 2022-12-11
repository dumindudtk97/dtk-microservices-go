package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbitmq *amqp.Connection
}

func main() {

	// try to connect rabbitmq running on docker
	conn, err := connect()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	app := Config{
		Rabbitmq: conn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:           ":" + webPort,
		Handler:        app.routes(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start http server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

// rabbitmq might start slow so need a backoff routine
func connect() (*amqp.Connection, error) {
	var counts int8
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		amqpURL := "amqp://guest:guest@rabbitmq" //docker service
		c, err := amqp.Dial(amqpURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	log.Println("connected")

	return connection, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
		os.Exit(1)
	}
}
