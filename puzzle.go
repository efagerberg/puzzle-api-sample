package main

import (
	"database/sql"
)

// Puzzle is a representation of a database row in pizzles
type puzzle struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level uint   `json:"level"`
}

func (p *puzzle) getPuzzle(db *sql.DB) error {
	return db.QueryRow("SELECT name, level FROM puzzles WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Level)
}

func (p *puzzle) updatePuzzle(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE puzzles SET name=$1, level=$2 WHERE id=$3",
			p.Name, p.Level, p.ID)

	return err
}

func (p *puzzle) deletePuzzle(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM puzzles WHERE id=$1", p.ID)

	return err
}

func (p *puzzle) createPuzzle(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO puzzles(name, level) VALUES($1, $2) RETURNING id",
		p.Name, p.Level).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getPuzzles(db *sql.DB, start, count int) ([]puzzle, error) {
	rows, err := db.Query(
		"SELECT id, name, level FROM puzzles LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	puzzles := []puzzle{}

	for rows.Next() {
		var p puzzle
		if err := rows.Scan(&p.ID, &p.Name, &p.Level); err != nil {
			return nil, err
		}
		puzzles = append(puzzles, p)
	}

	return puzzles, nil
}
