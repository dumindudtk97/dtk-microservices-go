package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// get send request with body containing msg, encode, send with mailer.SendSMTPMessage
func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	// maps to request json
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	// decode request payload (json) to mailMessage format
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Sprintf("requestPayload: %+v", requestPayload)
	log.Printf("requestPayload in sendMail is %+v", requestPayload)

	// msg to send to mailer to send email
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
