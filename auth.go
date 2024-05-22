package main

import (
	"net/http"

	"github.com/nsp5488/blog_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConf *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authString, err := extractAuthHeader(req, "ApiKey")

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header.")
			return
		}

		user, err := apiConf.DB.GetUserByApiKey(req.Context(), authString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not retreive user from Database")
		}
		handler(w, req, user)
	})

}
