package main

import (
	"net/http"
)


func (app *Config) Broker (w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Message: "Hit the broker endpoint",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}