package main

import (
	"broker/cmd/api/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// define a json request payload (map)
type RequestPayload struct {
	Action          string          `json:"action"`
	Auth            AuthPayload     `json:"auth,omitempty"`
	Log             LogPayload      `json:"log,omitempty"`
	Mail            MailPayload     `json:"mail,omitempty"`
	RabbitmqPayload RabbitmqPayload `json:"rabbitmqPayload,omitempty"`
}

type RabbitmqPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

// handle all requests
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Sprintf("requestPayload: %+v", requestPayload)
	log.Printf("requestPayload in HS is %+v", requestPayload)

	switch requestPayload.Action {
	case "auth":
		log.Println("auth selected")
		app.authenticate(w, requestPayload.Auth) //handle authentication
	case "log":
		log.Println("log selected")
		app.logItem(w, requestPayload.Log) //handle log
	case "mail":
		log.Println("mail selected")
		app.sendMail(w, requestPayload.Mail) //handle mail
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log" //name used in docker compose file

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	// create json to send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t") // _ this is ignored error returning from fn
	log.Printf("jsonData: %#v", jsonData)
	// build the request
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) // prepare json
	if err != nil {
		log.Printf("preparing request failed")
		app.errorJSON(w, err)
		return
	}

	// call authentication service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("calling auth sevice failed")
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close() //defer delay exec of a function (close) until nearby function returns

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error response from calling auth service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// auth service set this true in errorJson when invalid credentials or bad request
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// build jsonResponse to send to client
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	//send response to client
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {

	jsonData, _ := json.MarshalIndent(m, "", "\t")
	log.Printf("jsonData: %#v", jsonData)

	mailServiceURL := "http://mail-service/send"

	// build the request
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData)) // prepare json
	if err != nil {
		log.Printf("preparing mail request failed")
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	// call mail service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("calling mail sevice failed")
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close() //defer delay exec of a function (close) until nearby function returns

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error response from calling mail service"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "email sent to " + m.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// logEventwithRabbit: logs event using the logger-service.
// done by pushing to rabbitmq
// listner-service receives event and call logging-service to do the logging
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	publisher, err := event.NewEventPublisher(app.Rabbitmq)
	if err != nil {
		return err
	}

	rabbitpayload := RabbitmqPayload{
		Name: name,
		Data: msg,
	}

	payload := RequestPayload{
		Action:          "rabbiLog",
		RabbitmqPayload: rabbitpayload,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = publisher.Push(string(j), "rabbiLog.INFO")
	if err != nil {
		return err
	}
	return nil
}
