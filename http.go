package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/syfun/inout/models"
)

// JSON is a map, which key is string and value is interface{}.
type JSON map[string]interface{}

// Response interface
type Response interface {
	Write(http.ResponseWriter)
}

// JSONResponse is json format response
type JSONResponse struct {
	Data interface{}
	Code int
}

func (jsonRes *JSONResponse) Write(w http.ResponseWriter) {
	WriteJSON(w, jsonRes.Data, jsonRes.Code)
}

// TextResponse is text format response
type TextResponse struct {
	Data string
	Code int
}

func (txtRes *TextResponse) Write(w http.ResponseWriter) {
	if txtRes.Code != 0 {
		w.WriteHeader(txtRes.Code)
	}
	fmt.Fprintln(w, txtRes.Data)
}

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

// // CtrlOptionMap option to method map.
// var CtrlOptionMap = map[string]string{
// 	"get":    "Get",
// 	"all":    "All",
// 	"update": "Update",
// 	"delete": "Delete",
// 	"create": "Create",
// }

// Register register controller.
func (r *Router) Register(ctrl *RestController) {
	single := fmt.Sprintf("/%ss/:id", ctrl.Name)
	multi := fmt.Sprintf("/%ss", ctrl.Name)
	for _, opt := range ctrl.Options {
		switch opt {
		case "get":
			r.Get(single, ctrl.Get)
		case "all":
			r.Get(multi, ctrl.All)
		case "update":
			r.Patch(single, ctrl.Update)
		case "delete":
			r.Delete(single, ctrl.Delete)
		case "create":
			r.Post(multi, ctrl.Create)
		}
	}
}

type httpHandler func(http.ResponseWriter, *http.Request) (Response, error)

type httpError struct {
	error   error
	Message string
	Code    int
}

func (he httpError) Error() string {
	return fmt.Sprintf("%s: %s", he.Message, he.error.Error())
}

func wrapHandler(h httpHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), key("params"), ps)
		resp, err := h(w, r.WithContext(ctx))
		if err != nil {
			code := 500
			if err, ok := err.(*httpError); ok {
				code = err.Code
			}
			log.Printf("Handle error, %s", err.Error())
			WriteJSON(w, JSON{"error": err.Error()}, code)
			return
		}
		resp.Write(w)
	}
}

func getParams(r *http.Request) httprouter.Params {
	ctx := r.Context()
	params := ctx.Value(key("params"))
	if ps, ok := params.(httprouter.Params); ok {
		return ps
	}
	var ps httprouter.Params
	return ps
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

// Patch is wrap with httprouter.Router.PATCH
func (r *Router) Patch(path string, handler httpHandler) {
	r.PATCH(path, wrapHandler(handler))
}

// RestController is restful api controller.
type RestController struct {
	Model   *models.Model
	Name    string
	Options []string
}

// AllOptions ...
var AllOptions = []string{"all", "create", "get", "update", "delete"}

// GetBody get request real bytes.
func GetBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

var noneArgs = []reflect.Value{}

// Create create resource.
func (rest *RestController) Create(w http.ResponseWriter, r *http.Request) (Response, error) {
	body, err := GetBody(r)
	if err != nil {
		return nil, &httpError{err, "cannot get http body", 500}
	}
	resource := reflect.New(reflect.TypeOf(rest.Model.Table).Elem()).Elem().Addr()
	if err = json.Unmarshal(body, resource.Interface()); err != nil {
		return nil, &httpError{err, "cannot parse body", 400}
	}
	res, err := rest.Model.Insert(resource.Interface())
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &JSONResponse{res, 201}, nil
}

// Get get one resource.
func (rest *RestController) Get(w http.ResponseWriter, r *http.Request) (Response, error) {
	resource, err := rest.Model.Get(models.NewDBQuery(nil, map[string]string{"id": getParams(r).ByName("id")}))
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &JSONResponse{resource, 200}, nil
}

// All get all resources.
func (rest *RestController) All(w http.ResponseWriter, r *http.Request) (Response, error) {
	resources, err := rest.Model.All(models.NewDBQuery(r.URL.Query(), nil))
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &JSONResponse{resources, 200}, nil
}

// Update update one resource.
func (rest *RestController) Update(w http.ResponseWriter, r *http.Request) (Response, error) {
	body, err := GetBody(r)
	if err != nil {
		return nil, &httpError{err, "cannot get http body", 500}
	}
	data := make(map[string]interface{})
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, &httpError{err, "cannot parse body", 400}
	}
	if len(data) == 0 {
		return nil, &httpError{errors.New("body has no content"), "", 400}
	}
	query := models.NewDBQuery(nil, map[string]string{"id": getParams(r).ByName("id")})
	resource, err := rest.Model.Update(query, data)
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &JSONResponse{resource, 200}, nil
}

// Delete delete one resource.
func (rest *RestController) Delete(w http.ResponseWriter, r *http.Request) (Response, error) {
	err := rest.Model.Delete(models.NewDBQuery(nil, map[string]string{"id": getParams(r).ByName("id")}))
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &TextResponse{"", 204}, nil
}
