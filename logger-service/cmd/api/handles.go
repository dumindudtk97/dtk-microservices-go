package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	// read json into requestPayload
	var requestPayload JSONPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("inbound payload is %+v", requestPayload)

	// log entry has 2 fields, data model add dates
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	error := app.Models.LogEntry.Insert(event)
	if error != nil {
		log.Printf("mongo insert error")
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *Config) WriteLogTest(w http.ResponseWriter, r *http.Request) {

	// read json into requestPayload
	var requestPayload JSONPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log entry has 2 fields, data model add dates
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Printf("mongo insert error")
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("capable of getting name out of request payload: %s %s but get this error %s", requestPayload.Name, requestPayload.Data, err),
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
