package main

import (
	"database/sql"
	"errors"
)

// Puzzle is a representation of a database row in pizzles
type puzzle struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriLevelce uint   `json:"level"`
}

func (p *puzzle) getPuzzle(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *puzzle) updatePuzzle(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *puzzle) deletePuzzle(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *puzzle) createPuzzle(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getPuzzles(db *sql.DB, start, count int) ([]puzzle, error) {
	return nil, errors.New("Not implemented")
}
