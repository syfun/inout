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

	router.Register(&RestController{&models.Item{}, "item", AllOptions})
	router.Register(&RestController{&models.Label{}, "label", AllOptions})
	router.Register(&RestController{&models.Push{}, "push", AllOptions})
	router.Register(&RestController{&models.Pop{}, "pop", AllOptions})
	router.Register(&RestController{&models.Stock{}, "stock", AllOptions})
	router.NotFound = http.FileServer(
		http.Dir("/Users/sunyu/workspace/goprojects/src/github.com/syfun/inout/build")).ServeHTTP
	log.Fatal(http.ListenAndServe(":8000", router))
}
