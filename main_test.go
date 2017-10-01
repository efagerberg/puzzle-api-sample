// main_test.go

package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/efagerberg/puzzle-api-sample/app"
)

var a app.App

const tableCreationQuery = `CREATE table IF NOT EXISTS puzzles
(
id serial PRIMARY KEY,
name text NOT NULL,
level int NOT NULL CHECK(level > 0)
)`

func TestMain(m *testing.M) {
	a = app.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
	)
	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM puzzles")
	a.DB.Exec("ALTER SEQUENCE puzzles_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addPuzzles(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO puzzles(name, level) VALUES($1, $2)", "Puzzle "+strconv.Itoa(i), (i+1)*10)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/puzzles", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentPuzzle(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/puzzle/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Puzzle not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Puzzle not found'. Got '%s'", m["error"])
	}
}

func TestCreatePuzzle(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test puzzle","level":1}`)

	req, _ := http.NewRequest("POST", "/puzzle", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test puzzle" {
		t.Errorf("Expected puzzle name to be 'test puzzle'. Got '%v'", m["name"])
	}

	// the level and id are compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["level"] != 1.0 {
		t.Errorf("Expected puzzle level to be '1.0'. Got '%v'", m["puzzle"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected puzzle ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetPuzzle(t *testing.T) {
	clearTable()
	addPuzzles(1)

	req, _ := http.NewRequest("GET", "/puzzle/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdatePuzzle(t *testing.T) {
	clearTable()
	addPuzzles(1)

	req, _ := http.NewRequest("GET", "/puzzle/1", nil)
	response := executeRequest(req)
	var originalPuzzle map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPuzzle)

	payload := []byte(`{"name":"test puzzle - updated name","level":1}`)

	req, _ = http.NewRequest("PUT", "/puzzle/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalPuzzle["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalPuzzle["id"], m["id"])
	}

	if m["name"] == originalPuzzle["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalPuzzle["name"], m["name"], m["name"])
	}

	if m["level"] == originalPuzzle["level"] {
		t.Errorf("Expected the level to change from '%v' to '%v'. Got '%v'", originalPuzzle["level"], m["level"], m["level"])
	}
}

func TestDeletePuzzle(t *testing.T) {
	clearTable()
	addPuzzles(1)

	req, _ := http.NewRequest("GET", "/puzzle/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/puzzle/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/puzzle/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
