package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

func createLabel(w http.ResponseWriter, r *http.Request) error {
	body, err := getBody(r)
	if err != nil {
		return &httpError{err, "cannot get http body", 400}
	}
	label := &models.Label{}
	if err = json.Unmarshal(body, label); err != nil {
		return &httpError{err, "cannot parse body to Label", 400}
	}
	label, err = models.CreateLabel(label.Name, label.Type)
	if err != nil {
		return &httpError{err, "", 500}
	}
}
