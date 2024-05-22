package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

type responseFeed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func (apiConf *apiConfig) createFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type responseBody struct {
		Feed       responseFeed       `json:"feed"`
		FeedFollow responseFeedFollow `json:"feed_follow"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Unable to decode parameters in createFeed")
		respondWithError(w, http.StatusBadRequest, "Unable to decode request")
		return
	}

	dbParams := database.CreateFeedParams{
		ID:        generateUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}
	feed, err := apiConf.DB.CreateFeed(req.Context(), dbParams)
	if err != nil {
		log.Println("Unable to create feed")
		respondWithError(w, http.StatusBadRequest, "Unable to create feed")
		return
	}
	feedFollow, err := apiConf.DB.FollowFeed(req.Context(), database.FollowFeedParams{
		ID:        generateUUID(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Fatal("Error while creating follow relation between user and feed")
	}
	respBody := responseBody{
		Feed: responseFeed{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		},
		FeedFollow: responseFeedFollow{
			ID:        feedFollow.ID,
			UserID:    feedFollow.UserID,
			FeedID:    feedFollow.FeedID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
		},
	}
	respondWithJSON(w, http.StatusCreated, respBody)
}

func (apiConf *apiConfig) getFeeds(w http.ResponseWriter, req *http.Request) {
	feeds, err := apiConf.DB.GetFeeds(req.Context())

	if err != nil {
		log.Println("Unable to retrieve feeds!")
		respondWithError(w, http.StatusInternalServerError, "Could not retreive feeds at this time.")
	}

	respBody := make([]responseFeed, len(feeds))

	for i, feed := range feeds {
		respBody[i] = responseFeed{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, respBody)
}
