package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	type errorBody struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorBody{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, content interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = w.Write(data)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func readiness(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, response{Status: "ok"})
}

func err(w http.ResponseWriter, req *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
