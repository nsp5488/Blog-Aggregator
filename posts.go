package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

type respPost struct {
	ID          uuid.UUID  `json:"ID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func DBPosttoResp(dbPost database.GetPostsByUserRow) respPost {
	return respPost{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: &dbPost.Description.String,
		PublishedAt: &dbPost.PublishedAt.Time,
		FeedID:      dbPost.FeedID,
	}
}

func (apiConf *apiConfig) getPosts(w http.ResponseWriter, req *http.Request, user database.User) {
	limit := req.URL.Query().Get("limit")
	dbLimit := 10
	var err error

	if limit != "" {
		dbLimit, err = strconv.Atoi(limit)
		if err != nil {
			dbLimit = 10
		}
	}
	posts, err := apiConf.DB.GetPostsByUser(req.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(dbLimit),
	})
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch posts")
		return
	}
	respBody := make([]respPost, len(posts))

	for i, post := range posts {
		respBody[i] = DBPosttoResp(post)
	}

	respondWithJSON(w, http.StatusOK, respBody)

}
