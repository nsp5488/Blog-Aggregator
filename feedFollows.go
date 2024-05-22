package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

type responseFeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func convDBFeedFollowToResp(feedFollow database.FeedFollow) responseFeedFollow {
	return responseFeedFollow{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	}
}

func (apiConf *apiConfig) followFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	params := &parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	feedFollow, err := apiConf.DB.FollowFeed(req.Context(), database.FollowFeedParams{
		ID:        generateUUID(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to follow feed at this time")
		return
	}

	respondWithJSON(w, http.StatusCreated, convDBFeedFollowToResp(feedFollow))
}

func (apiConf *apiConfig) unfollowFeed(w http.ResponseWriter, req *http.Request, _ database.User) {
	id, err := uuid.Parse(req.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to unfollow that feed")
		return
	}
	err = apiConf.DB.UnfollowFeed(req.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete at this time.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (apiConf *apiConfig) getFollowedFeeds(w http.ResponseWriter, req *http.Request, user database.User) {
	feedFollows, err := apiConf.DB.GetFollowedFeeds(req.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve feeds at this time")
		return
	}

	respBody := make([]responseFeedFollow, len(feedFollows))

	for i, feedFollow := range feedFollows {
		respBody[i] = convDBFeedFollowToResp(feedFollow)
	}
	respondWithJSON(w, http.StatusOK, respBody)
}
