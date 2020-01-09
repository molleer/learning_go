package main

import (
	"encoding/json"
	"fmt"
)

type Test struct {
	a int
	b int
	c int
}

func main() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	in := `{"ID":100, "Name":"Reds", "Colors":["Blue"],"Hello":"There"}`
	var res ColorGroup
	if err := json.Unmarshal([]byte(in), &res); err != nil {
		fmt.Println("Could not convert to struct")
		fmt.Println(err)
	}

	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
	fmt.Println(res)
}
