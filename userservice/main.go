package main

/**
* Pit falls:
* 	PostgreSQL query has placeholders $1, $2, $3.. not ?, ?, ?..
*	json.Marshal requires the names of the variables in a struct begins with a capital letter
 */

import (
	"github.com/dgrijalva/jwt-go"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"errors"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "go"
	password = "abc123"
	dbname   = "godb"
)

type JWTClaims struct {
	UUID string
	jwt.StandardClaims
}

var db *sql.DB
var SERVICE_SECRET = []byte(os.Getenv("GO_USER_SERVICE_SECRET"))

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

	http.HandleFunc("/api/new", method(http.MethodPost, handleNewUser))
	http.HandleFunc("/api/login", method(http.MethodPost, handleLogin))
	http.HandleFunc("/api/user", method(http.MethodGet, auth(handleGetUser)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func auth(handler func(w http.ResponseWriter, r *http.Request)) (func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request){
		token := r.Header.Get("Authorization")
		UUID, err := validateJWT(token)
		if err != nil {
			log.Println("Unauthorized request to %s", r.URL)
			http.Error(w, "Invaild JWT token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("UUID", UUID)
		handler(w,r)
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	UUID := r.Header.Get("UUID")
	user, err := getUser(UUID)
	if err != nil {
		log.Printf("Unable to find user with UUID: %s", UUID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	jsonUser, _ := json.Marshal(user)
	fmt.Fprint(w, string(jsonUser))
}

func handleNewUser(w http.ResponseWriter, r *http.Request){
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
	var body struct {
		Email string
		Password string
	}

	var respons struct {
		Jwt string
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

	jwt, err := createJWT(user.UUID)
	if err != nil {
		log.Println(err)
	}

	respons.Jwt = jwt

	jsonResp, er := json.Marshal(respons)
	if er != nil {
		http.Error(w, "Unable to parse user to json",http.StatusInternalServerError)
	}

	fmt.Fprint(w,string(jsonResp))
}

func createJWT(UUID string) (string, error){

	var claims JWTClaims

	claims.UUID = UUID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SERVICE_SECRET)
}

func validateJWT(token string) (string, error){
	
	keyFunc := func(t *jwt.Token) (interface{}, error) {return SERVICE_SECRET, nil}
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, keyFunc)
	if err != nil {
		return "", err
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok || !parsedToken.Valid {
		return "", errors.New("Invalid JWT token")
	}

	return claims.UUID, nil
}

func method(method string, handler func(w http.ResponseWriter, r *http.Request)) (func(w http.ResponseWriter, r *http.Request)){
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method != method {
			http.Error(w, 
				fmt.Sprintf("Method not supported: %s",r.Method), 
				http.StatusBadRequest)
			return
		}

		handler(w, r)
	}
}