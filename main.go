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

	router.Register(&RestController{
		&models.Model{Table: &models.Item{}},
		"item", AllOptions,
	})
	router.Register(&RestController{
		&models.Model{Table: &models.Label{}},
		"label", AllOptions,
	})

	router.Register(&RestController{
		&models.Model{Table: &models.Push{}},
		"push", AllOptions,
	})

	router.Register(&RestController{
		&models.Model{Table: &models.Pop{}},
		"pop", AllOptions,
	})

	router.Register(&RestController{
		&models.Model{Table: &models.Stock{}},
		"stock", AllOptions,
	})

	log.Fatal(http.ListenAndServe(":8000", router))
}
