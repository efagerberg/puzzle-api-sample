package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const tableCreationQuery = `CREATE table IF NOT EXISTS puzzles
(
id serial PRIMARY KEY,
name text NOT NULL,
level int NOT NULL CHECK(level > 0)
)`

// App hp;ds our router and db connection
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize creates mux router and database connection
func (a *App) Initialize(user, dbname, host, port string) {
	connectionString := fmt.Sprintf(
		"user=%s dbname=%s host=%s port=%s sslmode=disable",
		user, dbname, host, port,
	)
	print(connectionString)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()

	a.SeedDatabase()
	a.initializeRoutes()
}

// Run starts the app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// SeedDatabase initialize database table
func (a *App) SeedDatabase() {
	_, err := a.DB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.DB.Exec("INSERT INTO puzzles(name, level) VALUES ('Cake', 3)")
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/puzzle", a.GetPuzzles).Methods("GET")
	a.Router.HandleFunc("/puzzle/{id}", a.GetPuzzle).Methods("GET")
	a.Router.HandleFunc("/puzzle/{id}", a.CreatePuzzle).Methods("POST")
	a.Router.HandleFunc("/puzzle/{id}", a.DeletePuzzle).Methods("DELETE")
}

// GetPuzzles returns all puzzle objects as a response
func (a *App) GetPuzzles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello World")
}

// GetPuzzle returns a puzzle object as a response
func (a *App) GetPuzzle(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	json.NewEncoder(w).Encode("Hello World")

}

// CreatePuzzle creates a puzzle object in our database
func (a *App) CreatePuzzle(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	json.NewEncoder(w).Encode("Hello World")

}

// DeletePuzzle deletes puzzle from database if exists
func (a *App) DeletePuzzle(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	json.NewEncoder(w).Encode("Hello World")
}
