package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

func (apiConf *apiConfig) createFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Unable to decode parameters in createFeed")
		respondWithError(w, http.StatusBadRequest, "Unable to decode request")
		return
	}
	uuid, err := uuid.NewV6()
	if err != nil {
		log.Fatal("Unable to allocate new UUID")
	}
	dbParams := database.CreateFeedParams{
		ID:        uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}
	feed, err := apiConf.DB.CreateFeed(req.Context(), dbParams)
	if err != nil {
		log.Println("Unable to create feed")
		return
	}
	respondWithJSON(w, http.StatusCreated, responseBody{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	})
}
