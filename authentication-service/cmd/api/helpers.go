package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	//MaxBytesReader limit the size of incoming request bodies, error of type *MaxBytesError for a Read beyond the limit
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body) //returns a new decoder that reads from r
	err := dec.Decode(data)        //Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by data
	if err != nil {
		return err
	}

	// a trick to check if file has one json value
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data) //Marshal returns the JSON encoding of data
	if err != nil {
		return err
	}

	// write Headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value //(http.ResponseWriter).Header() http.Header returns the header map that will be sent by WriteHeader.
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
