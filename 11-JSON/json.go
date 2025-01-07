package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
}

func main() {
	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		fmt.Fprintf(w, "%s %s is %d years old!", user.FirstName, user.LastName, user.Age)
	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		peter := User{
			FirstName: "john",
			LastName:  "doe",
			Age:       25,
		}
		json.NewEncoder(w).Encode(peter)
	})

	http.ListenAndServe(":8080", nil)
}
