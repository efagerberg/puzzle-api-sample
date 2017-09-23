package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// GetPuzzles returns all puzzle objects as a response
func GetPuzzles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello World")
}

// GetPuzzle returns a puzzle object as a response
func GetPuzzle(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	json.NewEncoder(w).Encode("Hello World")

}

// CreatePuzzle creates a puzzle object in our database
func CreatePuzzle(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	json.NewEncoder(w).Encode("Hello World")

}

// DeletePuzzle deletes puzzle from database if exists
func DeletePuzzle(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	json.NewEncoder(w).Encode("Hello World")

}

func main() {

	db, err := sql.Open("postgres", "user=postgres dbname=postgres host=database port=5432 sslmode=disable")
	CheckError(err)
	defer db.Close()

	age := 24
	_, err = db.Exec("CREATE table IF NOT EXISTS users (id serial PRIMARY KEY, name text, age int)")
	CheckError(err)
	_, err = db.Exec("INSERT INTO users(name, age) VALUES ('Evan Fagerberg', 24)")
	CheckError(err)
	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	CheckError(err)
	log.Print(rows)

	router := mux.NewRouter()
	router.HandleFunc("/puzzle", GetPuzzles).Methods("GET")
	router.HandleFunc("/puzzle/{id}", GetPuzzle).Methods("GET")
	router.HandleFunc("/puzzle/{id}", CreatePuzzle).Methods("POST")
	router.HandleFunc("/puzzle/{id}", DeletePuzzle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// CheckError A function for handling errors
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
