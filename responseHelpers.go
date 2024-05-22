package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
		log.Println(err.Error())
		return
	}

	_, err = w.Write(data)

	if err != nil {
		log.Println(err.Error())
	}
}
func extractAuthHeader(req *http.Request, authType string) (string, error) {
	s := req.Header.Get("Authorization")

	if s == "" {
		return "", errors.New("no authorization header provided")
	}

	authString, found := strings.CutPrefix(s, authType+" ")

	if !found {
		return "", errors.New("authorization header type mismatch")
	}
	return authString, nil

}

func generateUUID() uuid.UUID {
	uuid, err := uuid.NewV6()
	if err != nil {
		log.Fatal("Could not generate a new UUID!")
	}
	return uuid
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
