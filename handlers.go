package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/syfun/inout/models"
)

func getBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

func createLabel(w http.ResponseWriter, r *http.Request) (Response, error) {
	body, err := getBody(r)
	if err != nil {
		return nil, &httpError{err, "cannot get http body", 500}
	}
	label := &models.Label{}
	if err = json.Unmarshal(body, label); err != nil {
		return nil, &httpError{err, "cannot parse body to Label", 400}
	}
	label, err = models.CreateLabel(label.Name, label.Type)
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &JSONResponse{label, 201}, nil
}

func getLabels(w http.ResponseWriter, r *http.Request) (Response, error) {
	labels, err := models.GetLabels(r.URL.Query().Get("type"))
	if err != nil {
		return nil, &httpError{err, "", 500}
	}
	return &JSONResponse{labels, 200}, nil
}

func updateLabel(w http.ResponseWriter, r *http.Request) (Response, error) {
	labelID, _ := strconv.ParseInt(getParams(r).ByName("labelID"), 10, 64)
	body, err := getBody(r)
	if err != nil {
		return nil, &httpError{err, "cannot get http body", 500}
	}

	label := &models.Label{}
	if err = json.Unmarshal(body, label); err != nil {
		return nil, &httpError{err, "cannot parse body to Label", 400}
	}

	label, err = models.UpdateLabel(labelID, label.Name)
	if err != nil {
		return nil, &httpError{err, "", 500}
	}

	return &JSONResponse{label, 200}, nil
}
