package main

import (
	"log"
	"net/http"

	"github.com/syfun/inout/models"
)

func main() {
	models.InitDB("./test.db")
	defer models.CloseDB()

	router := NewRouter()

	router.Post("/labels", createLabel)
	router.Get("/labels", getLabels)
	router.Post("/labels/:labelID", updateLabel)

	log.Fatal(http.ListenAndServe(":8000", router))
}
