package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nsp5488/blog_aggregator/internal/database"
)

func addHandlers(mux *http.ServeMux, apiConf apiConfig) {
	mux.HandleFunc("GET /v1/readiness", readiness)
	mux.HandleFunc("GET /v1/err", err)

	mux.HandleFunc("POST /v1/users", apiConf.createUsers)
	mux.HandleFunc("GET /v1/users", apiConf.middlewareAuth(apiConf.getUserAuthed))

	mux.HandleFunc("POST /v1/feeds", apiConf.middlewareAuth(apiConf.createFeed))
	mux.HandleFunc("GET /v1/feeds", apiConf.getFeeds)

	mux.HandleFunc("POST /v1/feed_follows", apiConf.middlewareAuth(apiConf.followFeed))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConf.middlewareAuth(apiConf.unfollowFeed))
	mux.HandleFunc("GET /v1/feed_follows", apiConf.middlewareAuth(apiConf.getFollowedFeeds))
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("PSQL_CONN")

	if port == "" {
		log.Println("Could not load port number from ENV file, defaulting to port 8080")
		port = "8080"
	}
	if dbUrl == "" {
		log.Fatal("Could not retreive connection string from ENV file")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error while establishing DB connection")
	}
	dbQueries := database.New(db)
	apiConf := apiConfig{dbQueries}

	mux := http.NewServeMux()
	addHandlers(mux, apiConf)

	server := http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
