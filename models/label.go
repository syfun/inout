package models

import (
	"database/sql"
)

// Label can be user, warehouse, item type
type Label struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	// Type can be user, warehouse, item type
	Type string `json:"type"`
}

// DefaultLabel for Get and All.
var DefaultLabel = &Label{}

// Get get label by id.
func (lb *Label) Get(query *DBQuery) (interface{}, error) {
	label := &Label{}
	err := db.QueryRow("select * from label where id=?", query.Get("id")).Scan(
		&label.ID, &label.Name, &label.Type,
	)
	if err != nil {
		return nil, &DBError{err, "label not found"}
	}
	return label, nil
}

// GetLabel wrap Get.
func GetLabel(query *DBQuery) (interface{}, error) {
	return DefaultLabel.Get(query)
}

// All query labels
func (lb *Label) All(query *DBQuery) (interface{}, error) {
	var (
		rows *sql.Rows
		err  error
	)
	typ := query.Get("type")
	if typ == "" {
		rows, err = db.Query("select * from label;")
	} else {
		rows, err = db.Query("select * from label where type=?", typ)
	}
	if err != nil {
		return nil, &DBError{err, "cannot select from label"}
	}
	labels := make([]Label, 0)
	for rows.Next() {
		var label Label
		rows.Scan(&label.ID, &label.Name, &label.Type)
		labels = append(labels, label)
	}
	return labels, nil
}

// AllLabels wrap All.
func AllLabels(query *DBQuery) (interface{}, error) {
	return DefaultLabel.All(query)
}

// Insert insert label
func (lb *Label) Insert() error {
	res, err := db.Exec("insert into label (name, type) values (?, ?);", lb.Name, lb.Type)
	if err != nil {
		return &DBError{err, "insert error"}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return &DBError{err, "cann get last insert id"}
	}
	lb.ID = id
	return nil
}

// Update update label name.
func (lb *Label) Update() error {
	_, err := db.Exec("update label set name=? where id=?", lb.Name, lb.ID)
	if err != nil {
		return &DBError{err, "cannot update label"}
	}
	return nil
}

// Delete delete label.
func (lb *Label) Delete(query *DBQuery) error {
	_, err := db.Exec("delete from label where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "cannot delete label"}
	}
	return nil
}
