package models

// Label can be user, warehouse, item type
type Label struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	// Type can be user, warehouse, item type
	Type string `json:"type"`
}

func (*Label) columns() map[string]interface{} {
	return stringSliceToMap([]string{"id", "name", "type"})
}

func (*Label) insertStmt() string {
	return `insert into label (name, type) values (:name, :type)`
}

func (*Label) getStmt() string {
	return `select * from label where id=?`
}

func (*Label) allStmt() string {
	return `select * from label`
}

func (*Label) updateStmt() string {
	return "update label"
}

func (*Label) deleteStmt() string {
	return "delete from label where id=?"
}
