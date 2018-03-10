package models

// Item ...
type Item struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Specification string `json:"specification"`
	Unit          string `json:"unit"`
	Push          int64  `json:"push"`
	Pop           int64  `json:"pop"`
	Now           int64  `json:"now"`
	Desc          string `json:"desc"`
}

func (*Item) columns() map[string]interface{} {
	return stringSliceToMap([]string{"id", "name", "type", "specification", "unit", "push", "pop", "now", "desc"})
}

func (*Item) insertStmt() string {
	return `insert into item (name, type, specification, unit, push, pop, now, desc) values (:name, :type, :specification, :unit, :push, :pop, :now, :desc);`
}

func (*Item) getStmt() string {
	return `select * from item where id=?`
}

func (*Item) allStmt() string {
	return `select * from item`
}

func (*Item) updateStmt() string {
	return "update item"
}

func (*Item) deleteStmt() string {
	return "delete from item where id=?"
}
