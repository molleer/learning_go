package main

import (
	//"crypto/rand"
	"github.com/dgrijalva/jwt-go"
	//"encoding/base64"
	"fmt"
)

func main() {
	secret := []byte("noTMwMrtsxtYfEFt+VaTXG3mEswCOMVwKpAhjRRWy40=")

	type Hello struct{
		Hello string
	}

	var claims struct {
		Secret Hello
		jwt.StandardClaims
	}

	claims.Secret = Hello{"Hello there"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	fmt.Printf("%v %v", ss, err)
}
