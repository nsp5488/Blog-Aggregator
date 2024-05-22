package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func addHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/readiness", readiness)
	mux.HandleFunc("GET /v1/err", err)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Could not load port number from ENV file, defaulting to port 8080")
		port = "8080"
	}

	mux := http.NewServeMux()
	addHandlers(mux)

	server := http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
