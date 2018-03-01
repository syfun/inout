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

// CreateLabel insert label
func CreateLabel(name, typ string) (*Label, error) {
	res, err := db.Exec("insert into label (name, type) values (?, ?);", name, typ)
	if err != nil {
		return nil, &DBError{err, "insert error"}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, &DBError{err, "cann get last insert id"}
	}
	return &Label{id, name, typ}, nil
}

// GetLabels query labels
func GetLabels(typ string) ([]Label, error) {
	var (
		rows *sql.Rows
		err  error
	)
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

// UpdateLabel update label name.
func UpdateLabel(id int64, name string) (*Label, error) {
	_, err := db.Exec("update label set name=? where id=?", name, id)
	if err != nil {
		return nil, &DBError{err, "cannot update label"}
	}
	var typ string
	db.QueryRow("select type from label where id=?", id).Scan(&typ)
	return &Label{id, name, typ}, err
}
