package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// JSON is a map, which key is string and value is interface{}.
type JSON map[string]interface{}

// WriteJSON write json data to response.
func WriteJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	if code != 0 {
		w.WriteHeader(code)
	}
	json.NewEncoder(w).Encode(v)
}

// Router wrap for httprouter.Router
type Router struct {
	*httprouter.Router
}

type key string

// NewRouter create a new Router
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

type httpHandler func(http.ResponseWriter, *http.Request) error

type httpError struct {
	error   error
	Message string
	Code    int
}

func (he httpError) Error() string {
	return he.Message
}

func wrapHandler(h httpHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), key("params"), ps)
		if err := h(w, r.WithContext(ctx)); err != nil {
			code := 500
			if err, ok := err.(httpError); ok {
				code = err.Code
			}
			log.Printf("Handle error, %s", err.Error())
			WriteJSON(w, JSON{"error": err.Error()}, code)
		}
	}
}

// Get is wrap with httprouter.Router.GET.
func (r *Router) Get(path string, handler httpHandler) {
	r.GET(path, wrapHandler(handler))
}

// Post is wrap with httprouter.Router.POST
func (r *Router) Post(path string, handler httpHandler) {
	r.POST(path, wrapHandler(handler))
}

// Put is wrap with httprouter.Router.PUT
func (r *Router) Put(path string, handler httpHandler) {
	r.PUT(path, wrapHandler(handler))
}

// Delete is wrap with httprouter.Router.DELETE
func (r *Router) Delete(path string, handler httpHandler) {
	r.DELETE(path, wrapHandler(handler))
}
