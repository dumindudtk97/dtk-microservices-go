package main

import (
	"log"
	"net/http"
	"time"
)

const webPort = "80"

type Config struct{}

func main() {

	app := Config{}

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
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
