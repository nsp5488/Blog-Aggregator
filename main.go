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
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("PSQL_CONN")

	if port == "" {
		log.Fatal("Could not load port number from ENV file, defaulting to port 8080")
		port = "8080"
	}
	if dbUrl == "" {
		log.Fatal("Could not retreive connection string from ENV file")
		os.Exit(-1)
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error while establishing DB connection")
		os.Exit(-1)
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