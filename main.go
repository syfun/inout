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

	// ctrl := &RestController{&models.Label{}, []string{"POST"}}
	// router.Post("/labels", ctrl.Create)
	// router.Get("/labels/:id", ctrl.Get)
	// router.Get("/labels", ctrl.All)
	// router.Patch("/labels/:id", ctrl.Update)
	// router.Delete("/labels/:id", ctrl.Delete)

	ctrl := &RestController{&models.Item{}, []string{"POST"}}
	router.Get("/items/:id", ctrl.Get)
	router.Post("/items", ctrl.Create)
	router.Get("/items", ctrl.All)
	// router.Get("/labels", getLabels)
	// router.Post("/labels/:labelID", updateLabel)

	log.Fatal(http.ListenAndServe(":8000", router))
}
