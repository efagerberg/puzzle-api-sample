package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/efagerberg/puzzle-api-sample/utils"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const tableCreationQuery = `CREATE table IF NOT EXISTS puzzles
(
id serial PRIMARY KEY,
name text NOT NULL,
level int NOT NULL CHECK(level > 0)
)`

// App handles our router and db connection
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
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	utils.CheckError(err)
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
	utils.CheckError(err)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/puzzles", a.getPuzzles).Methods("GET")
	a.Router.HandleFunc("/puzzle", a.createPuzzle).Methods("POST")
	a.Router.HandleFunc("/puzzle/{id:[0-9]+}", a.getPuzzle).Methods("GET")
	a.Router.HandleFunc("/puzzle/{id:[0-9]+}", a.updatePuzzle).Methods("PUT")
	a.Router.HandleFunc("/puzzle/{id:[0-9]+}", a.deletePuzzle).Methods("DELETE")
}

func (a *App) getPuzzles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	puzzles, err := getPuzzles(a.DB, start, count)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, puzzles)
}

func (a *App) getPuzzle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid puzzle ID")
		return
	}

	p := puzzle{ID: id}
	if err := p.getPuzzle(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "Puzzle not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, p)

}

func (a *App) createPuzzle(w http.ResponseWriter, r *http.Request) {
	var p puzzle
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createPuzzle(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, p)
}

func (a *App) deletePuzzle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Puzzle ID")
		return
	}

	p := puzzle{ID: id}
	if err := p.deletePuzzle(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updatePuzzle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid puzzle ID")
		return
	}

	var p puzzle
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updatePuzzle(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, p)
}
