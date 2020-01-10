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

	createUserTable()

	http.HandleFunc("/api/new", handleNewUser)
	http.HandleFunc("/api/login", handleLogin)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleNewUser(w http.ResponseWriter, r *http.Request){
	if !requireMethod(w,r,http.MethodPost) {return}

	var body struct {
		User User
		Password string
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := newUser(body.User, body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleLogin(w http.ResponseWriter, r *http.Request){
	if !requireMethod(w,r,http.MethodPost) {return}

	var body struct {
		Email string
		Password string
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := login(body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if jsonUser, er := json.Marshal(user); er != nil {
		http.Error(w, "Unable to parse user to json",http.StatusInternalServerError)
	} else {
		fmt.Fprint(w,string(jsonUser))
	}
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool{
	if r.Method != method {
		http.Error(w, 
			fmt.Sprintf("Method not supported: %s",r.Method), 
			http.StatusBadRequest)
		return false
	}

	return true
}