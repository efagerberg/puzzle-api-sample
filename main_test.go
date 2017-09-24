// main_test.go

package main_test

import (
	"log"
	"os"
	"testing"

	"."
)

var a main.App

const tableCreationQuery = `CREATE table IF NOT EXISTS puzzles
(
id serial PRIMARY KEY,
name text NOT NULL,
evel int NOT NULL CHECK(level > 0)
)`

func TestMain(m *testing.M) {
	a = main.App{}
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
