package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()
	router.Get("/gg", func(w http.ResponseWriter, r *http.Request) error {
		return httpError{nil, "tttt", 200}
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}
