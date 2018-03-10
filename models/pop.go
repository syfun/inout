package models

// Pop ...
type Pop struct {
	ID            int64  `json:"id"`
	ItemID        int64  `json:"item_id" db:"item_id"`
	Time          string `json:"time"`
	Number        int64  `json:"number"`
	Warehouse     string `json:"warehouse"`
	Abstract      string `json:"abstract"`
	Remark        string `json:"remark"`
	Receiver      string `json:"receiver"`
	Checker       string `json:"checker"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Specification string `json:"specification"`
	Unit          string `json:"unit"`
}

func (*Pop) columns() map[string]interface{} {
	return stringSliceToMap([]string{"id", "item_id", "time", "number", "warehouse", "abstract", "remark", "receiver", "checker"})
}

func (*Pop) insertStmt() string {
	return `insert into pop (item_id, time, number, warehouse, abstract, remark, receiver, checker) values (:item_id, :time, :number, :warehouse, :abstract, :remark, :receiver, :checker);`
}

func (*Pop) getStmt() string {
	return `select pop.id as id, item_id, time, number, warehouse, abstract, remark, receiver, checker, name, type, specification, unit from pop left outer join item on pop.item_id=item.id where pop.id=?`
}

func (*Pop) allStmt() string {
	return `select pop.id as id, item_id, time, number, warehouse, abstract, remark, receiver, checker, name, type, specification, unit from pop left outer join item on pop.item_id=item.id`
}

func (*Pop) updateStmt() string {
	return "update pop"
}

func (*Pop) deleteStmt() string {
	return "delete from pop where id=?"
}
