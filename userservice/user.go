package main

import (
	"database/sql"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
	"errors"
	"encoding/base64"
	"encoding/base32"
)

const (
	PW_SALT_BYTES = 32
    PW_HASH_BYTES = 64
)

type User struct {
	UUID string
	Email string
	First_name string 
	Last_name string
	Admin bool
}

func createUserTable() (sql.Result, error){
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		UUID UUID NOT NULL,
		Email TEXT NOT NULL UNIQUE,
		First_name TEXT NOT NULL,
		Last_name TEXT NOT NULL,
		Salt CHARACTER(56) NOT NULL,
		Password CHARACTER(88) NOT NULL,
		Admin BOOLEAN,
		PRIMARY KEY (UUID)
		)`)
}

func newUser(u User, password string) (error){
	salt := make([]byte, PW_SALT_BYTES)
    if _, err := rand.Read(salt); err != nil {
		return err
	}

	hashedPass, err := scrypt.Key([]byte(password), salt, 1 << 14, 8, 1, PW_HASH_BYTES)
	if err != nil { return err }

	_, err = db.Exec(`
	INSERT INTO users 
	(UUID, Email, First_name, Last_name, Salt, Password, Admin) 
	VALUES(uuid_generate_v4(), $1, $2, $3, $4, $5, $6)`,
		u.Email,
		u.First_name,
		u.Last_name,
		base32.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hashedPass),
		u.Admin)
	
	return err
}

func login(email string, pass string) (User, error){
	var u User
	var hashedPass string
	var salt string
	if err := db.QueryRow("SELECT * FROM users WHERE Email=$1", email).Scan(&u.UUID,
		&u.Email,
		&u.First_name,
		&u.Last_name,
		&salt,
		&hashedPass,
		&u.Admin); err != nil{
			return User{},err
		}

	saltBytes, _ := base32.StdEncoding.DecodeString(salt)

	hash, err := scrypt.Key([]byte(pass), saltBytes, 1 << 14, 8, 1, PW_HASH_BYTES)
	if err == nil && base64.StdEncoding.EncodeToString(hash) == hashedPass {
		return u, nil
	}

	return User{}, errors.New("The Email and password does not match")
}