package models

import (
	"database/sql"
)

// Label can be user, warehouse, item type
type Label struct {
	ID   int64
	Name string

	// Type can be user, warehouse, item type
	Type string
}

// CreateLabel insert label
func CreateLabel(name, typ string) (*Label, *DBError) {
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
func GetLabels(typ string) ([]Label, *DBError) {
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
	var labels []Label
	for rows.Next() {
		var label Label
		rows.Scan(&label.ID, &label.Name, &label.Type)
		labels = append(labels, label)
	}
	return labels, nil
}
