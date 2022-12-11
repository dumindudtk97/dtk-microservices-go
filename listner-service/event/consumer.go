package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// consumer for recieving events from rabbitmq
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// rabbitmq payloads
type RabbitmqPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// define a json request payload (map)
type RequestPayload struct {
	Action          string          `json:"action"`
	Auth            AuthPayload     `json:"auth,omitempty"`
	Log             LogPayload      `json:"log,omitempty"`
	Mail            MailPayload     `json:"mail,omitempty"`
	RabbitmqPayload RabbitmqPayload `json:"rabbitmqPayload,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// create a consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// this should declare exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// listen function listen to "topics" exchange // see rabbitmq docs
func (consumer *Consumer) Listen(topics []string) error {
	// get the channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// rabbitmq conventions see docs
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"topics",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	//Consume immediately starts delivering queued messages.
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// channel : Channels provide a way for two goroutines to communicate with one another and synchronize their execution.
	forever := make(chan bool)
	// go : run func in background , so the fn that launched it can  do other things //managed by golang run-time.
	go func() {
		for d := range messages {
			var payload RequestPayload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [topics, %s]\n", q.Name)

	<-forever // make main wait  (this is a trick to make main wait indefinetly)

	return nil
}

// handle the payload from the rabbitmq queue
func handlePayload(payload RequestPayload) {

	switch payload.Action {
	case "rabbiLog", "event":
		err := logEvent(payload.RabbitmqPayload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		err := authEvent(payload.Auth)
		if err != nil {
			log.Println(err)
		}

	default:
		err := logEvent(payload.RabbitmqPayload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry RabbitmqPayload) error {

	// call logger-service and log // same thing as in broker service

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log" //name used in docker compose file

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	// create http client
	client := &http.Client{}

	// call the logger service with entry payload
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	log.Println("entry logged :", entry)
	return nil
}

// same thing as broker
func authEvent(entry AuthPayload) error {
	// call authentication-service and authenticate

	// create json to send to the auth service from payload
	jsonData, _ := json.MarshalIndent(entry, "", "\t") // _ this is ignored error returning from fn
	log.Printf("jsonData: %#v", jsonData)

	// build the request
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) // prepare json
	if err != nil {
		log.Printf("preparing request failed")
		return err
	}

	// call authentication service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("calling auth sevice failed")
		return err
	}
	defer response.Body.Close() //defer delay exec of a function (close) until nearby function returns

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("invalid credentials")
	} else if response.StatusCode != http.StatusAccepted {
		return errors.New("error response from calling auth service")
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		return err
	}

	// auth service set this true in errorJson when invalid credentials or bad request
	if jsonFromService.Error {
		return err
	}

	println("authenticated :", entry.Email, entry.Password)
	return nil
}
