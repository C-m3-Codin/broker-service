package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit / post")
	payload := jsonRespose{
		Error:   false,
		Message: "Hit the Broker",
	}

	// out, _ := json.MarshalIndent(payload, "", "\t")
	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit /check Get")
	payload := jsonRespose{
		Error:   false,
		Message: "Hit the Broker",
	}

	// out, _ := json.MarshalIndent(payload, "", "\t")
	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit handle post")
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.errorJSON(w, errors.New("Unknown Error"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// send json  to the auth micro service
	jsonData, _ := json.MarshalIndent(a, " ", "\t")
	// call auth micro service
	request, err := http.NewRequest("POST", "http://localhost.com/authenticate", bytes.NewBuffer(jsonData))
	// get back status code
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid creds"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth"))
		return
	}

	// read response body
	var jsonFromAuth jsonRespose
	// decode
	err = json.NewDecoder(response.Body).Decode(&jsonFromAuth)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromAuth.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonRespose
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromAuth.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}
