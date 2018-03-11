package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

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
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	router.NotFound = http.FileServer(http.Dir(path.Join(dir, "build"))).ServeHTTP
	log.Fatal(http.ListenAndServe(":8000", router))
}
