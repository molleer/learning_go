package main

/**
* Pit falls:
* 	PostgreSQL query has placeholders $1, $2, $3.. not ?, ?, ?..
*	json.Marshal requires the names of the variables in a struct begins with a capital letter
 */

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "go"
	password = "abc123"
	dbname   = "godb"
)

var db *sql.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	time.Sleep(4 * time.Second)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//createPersonTable()
	createSuggestionsTable()

	http.HandleFunc("/api/", handleRoot)
	http.HandleFunc("/api/delete", handleDeleteSuggestion)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request){
	defer recoverBadRequest(w)

	switch r.Method {
	case http.MethodPost:
		handleInsert(w,r)
	case http.MethodGet:
		handleGetSuggestions(w,r)
	default:
		err := fmt.Sprintf("Request method not supported: %s", r.Method)
		panic(err)
	}
}

func handleInsert(w http.ResponseWriter, r *http.Request){
	defer recoverBadRequest(w)
	var s Suggestion
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &s); err != nil {
		panic("Could not parse body")
	}
	if _, err := insertSuggestion(s); err != nil{
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}

func handleGetSuggestions(w http.ResponseWriter, r *http.Request){
	var res []Suggestion
	var err error

	if res, err = getSuggestions(); err != nil {
		log.Println("Unable to get suggestions")
		http.Error(w,"Somthing went wrong", http.StatusInternalServerError)
	}

	ss,_ := json.Marshal(res)
	fmt.Fprint(w, string(ss))
}

func handleDeleteSuggestion(w http.ResponseWriter, r *http.Request){
	defer recoverBadRequest(w)
	if r.Method != http.MethodDelete {
		panic(fmt.Sprintf("Request method not supported: %s", r.Method))
	}
	if err := deleteSuggestion(r.URL.Query().Get("UUID")); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func recoverBadRequest(w http.ResponseWriter){
	if rec := recover(); rec != nil {
		log.Println(rec)
		http.Error(w, fmt.Sprintf("%s", rec), http.StatusBadRequest)
	}
}