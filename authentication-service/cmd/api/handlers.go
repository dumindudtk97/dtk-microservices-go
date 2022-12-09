package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type RequestPayloads struct {
	Email    string
	Password string
}

func (app *Config) personCreate(w http.ResponseWriter, r *http.Request) {

	var requestPayloads RequestPayloads

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&requestPayloads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Fprintf(w, "requestPayloads: %+v", requestPayloads)
	log.Printf("requestPayload is %+v", requestPayloads)

	user, err := app.Models.User.GetByEmail(requestPayloads.Email)
	log.Printf("requestPayload is %s", requestPayloads)
	log.Printf("email in auth handler is %s", requestPayloads.Email)
	if err != nil {
		app.errorJSON(w, errors.New("No user found"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayloads.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Wrong Password"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct { //payload of request is a json with email and password
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	// user is taken from db by email
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	log.Printf("requestPayload is %s", requestPayload)
	log.Printf("email in auth handler is %s", requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("No user found"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Wrong Password"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
