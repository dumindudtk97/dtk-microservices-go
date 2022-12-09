package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type RequestPayloads struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayloads RequestPayloads

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&requestPayloads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log to check request payload
	//fmt.Fprintf(w, "requestPayloads: %+v", requestPayloads)
	log.Printf("requestPayload is %+v", requestPayloads)

	user, err := app.Models.User.GetByEmail(requestPayloads.Email)
	log.Printf("email in auth handler is %s", requestPayloads.Email)
	if err != nil {
		app.errorJSON(w, errors.New("No user found"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayloads.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Wrong Password"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	log.Printf("outbound payload (prints only if logged in user) is %+v", payload)

	app.writeJSON(w, http.StatusAccepted, payload)
}
