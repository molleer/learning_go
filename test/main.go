package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	salt := make([]byte, 64)
    if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	a := base64.StdEncoding.EncodeToString(salt)
	fmt.Printf("%s \n%d\n",a,  len(a))
}
