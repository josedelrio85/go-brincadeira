package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bysidecar/go_components/mysqllibtest/implementeddb"
	"github.com/bysidecar/go_components/mysqllibtest/models"
	"github.com/bysidecar/go_components/readparams"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Env struct {
	db *sql.DB
}

var env *Env

func main() {
	connString := readparams.GetConnString(1)
	db, err := implementeddb.OpenConnection(connString)
	if err != nil {
		log.Panic(err)
		return
	}
	env = &Env{db: db}

	router := mux.NewRouter()
	router.HandleFunc("/sources", getSources).Methods("GET")
	router.HandleFunc("/sources/{id}", getSource).Methods("GET")
	router.HandleFunc("/sources", createSource).Methods("POST")
	router.HandleFunc("/sources/{id}", updateSource).Methods("PUT")
	// router.HandleFunc("/sources/{id}", deleteSource).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getSources(w http.ResponseWriter, r *http.Request) {
	elements, err := models.Index(env.db)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(elements)
}

func getSource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	el, err := models.Get(env.db, id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(el)
}

func createSource(w http.ResponseWriter, r *http.Request) {

	var source models.Sources

	err := json.NewDecoder(r.Body).Decode(&source)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// fmt.Println(source)

	el, err := models.Insert(env.db, source)
	if err != nil {
		log.Panic(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(el)
}

func updateSource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	var source models.Sources
	err := json.NewDecoder(r.Body).Decode(&source)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	source.Souid = id
	el, err := models.Update(env.db, source)
	if err != nil {
		log.Panic(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(el)
}
