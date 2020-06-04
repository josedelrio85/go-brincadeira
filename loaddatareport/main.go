package main

import (
	"log"
	"net/http"
	"time"

	client "github.com/bysidecar/go_components/loaddatareport/pkg"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	log.Printf("Import api started at %s", time.Now().Format("2006-01-02 15-04-05"))

	h := client.Handler{}

	router := mux.NewRouter()

	router.PathPrefix("/import/ws/").Handler(h.HandleFunction()).Methods(http.MethodPost)
	router.PathPrefix("/import/leontel/").Handler(h.HandleFunction()).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":4000", cors.Default().Handler(router)))
}
