package main

import (
	"database/sql"
)

type Suggestion struct {
	UUID string
	Timestamp string //not null TIMESTAMP
	Title string 	 //not null
	Text string		 //not null
	Author string	 //not null
}

func createSuggestionsTable() (sql.Result, error){
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS suggestions (
		UUID UUID NOT NULL,
		Timestamp TIMESTAMP NOT NULL,
		Title TEXT NOT NULL,
		Text TEXT NOT NULL,
		Author TEXT,
		PRIMARY KEY (UUID)
		)`)
}

func insertSuggestion(s Suggestion) (sql.Result, error){
	return db.Exec(`
	INSERT INTO suggestions 
	(UUID, Timestamp, Title, Text, Author) 
	VALUES(uuid_generate_v4(), now(), $1,$2,$3)`,
		s.Title, s.Text, s.Author)
}

func getSuggestions() ([]Suggestion, error) {
	rows, err := db.Query("SELECT * FROM suggestions")
	if err != nil {
		return nil, err
	}

	var suggestions []Suggestion
	for rows.Next() {
		var s Suggestion
		err = rows.Scan(
			&s.UUID,
			&s.Timestamp,
			&s.Title,
			&s.Text,
			&s.Author)

		if err == nil {
			suggestions = append(suggestions, s)
		}
	} 

	return suggestions, nil
}

func deleteSuggestion(UUID string) (error){
	_, err := db.Exec("DELETE FROM suggestions WHERE UUID=$1", UUID)
	return err
}