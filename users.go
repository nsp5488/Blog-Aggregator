package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

func (apiConf *apiConfig) createUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	type responseShape struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error while parsing request parameters in createUsers\n" + err.Error())
		return
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Println("Error whil generating UUID in createUsers")
		return
	}
	ctx := context.Background()
	dbParams := database.CreateUserParams{
		ID:        uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	}

	user, err := apiConf.DB.CreateUser(ctx, dbParams)
	if err != nil {
		log.Println("Error while creating user in database\n" + err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, responseShape{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	})
}

func (apiConf *apiConfig) getUserByKey(w http.ResponseWriter, req *http.Request) {
	authString, err := extractAuthHeader(req, "ApiKey")
	type responseShape struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header.")
		return
	}

	user, err := apiConf.DB.GetUserByApiKey(context.Background(), authString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retreive user from Database")
	}

	respondWithJSON(w, http.StatusOK, responseShape{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	})

}
