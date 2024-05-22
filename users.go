package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

type responseUser struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func convDBUserToResp(user database.User) responseUser {
	return responseUser{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

func (apiConf *apiConfig) createUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error while parsing request parameters in createUsers\n" + err.Error())
		return
	}

	dbParams := database.CreateUserParams{
		ID:        generateUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	}

	user, err := apiConf.DB.CreateUser(req.Context(), dbParams)
	if err != nil {
		log.Println("Error while creating user in database\n" + err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, convDBUserToResp(user))
}

func (apiConf *apiConfig) getUserAuthed(w http.ResponseWriter, req *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, convDBUserToResp)
}
