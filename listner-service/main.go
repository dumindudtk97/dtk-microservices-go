package main

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	// try to connect rabbitmq running on docker
	conn, err := connect()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	// start listning to msgs

	// create a consumer

	// watch the queue and consume events from topics

}

// rabbitmq might start slow so need a backoff routine
func connect() (*amqp.Connection, error) {
	var counts int8
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		amqpURL := "amqp://guest:guest@localhost"
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
	}
}
