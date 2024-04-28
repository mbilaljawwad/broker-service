package main

import (
	"log"
	"net/http"
)

const webPort = "2001"

type Config struct {}

func main () {
	app := Config{}
	log.Printf("Starting Broker service on Port %s\n", webPort)

	// define http server

	srv := &http.Server{
		Addr: ":" + webPort,
		Handler: app.routes(),
	}


	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}