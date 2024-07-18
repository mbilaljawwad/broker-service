package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	toolkit "github.com/mbilaljawwad/go-web-toolkit"
)

const AUTH_URL = "http://authentication-service/authenticate"

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var tools toolkit.Tools

	payload := toolkit.JSONResponse{
		Error:   false,
		Message: "Hit the broker endpoint",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var tools toolkit.Tools
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		tools.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	var tools toolkit.Tools
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", AUTH_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		tools.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonRes toolkit.JSONResponse
	err = json.NewDecoder(response.Body).Decode(&jsonRes)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	if jsonRes.Error {
		tools.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	acceptedPayload := toolkit.JSONResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonRes.Data,
	}

	tools.WriteJSON(w, http.StatusAccepted, acceptedPayload)
}
