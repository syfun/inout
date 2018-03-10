package models

// Stock ...
type Stock struct {
	ID            int64  `json:"id"`
	ItemID        int64  `json:"item_id" db:"item_id"`
	Warehouse     string `json:"warehouse"`
	Push          int64  `json:"push"`
	Pop           int64  `json:"pop"`
	Now           int64  `json:"now"`
	Desc          string `json:"desc"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Specification string `json:"specification"`
	Unit          string `json:"unit"`
}

func (*Stock) columns() map[string]interface{} {
	return stringSliceToMap([]string{"id", "item_id", "warehouse", "push", "pop", "now", "desc"})
}

func (*Stock) insertStmt() string {
	return `insert into stock (item_id, warehouse, push, pop, now, desc) values (:item_id, :warehouse, :push, :pop, :now, :desc);`
}

func (*Stock) getStmt() string {
	return `select stock.id as id, item_id, warehouse, warehouse, stock.push as push, stock.pop as pop, stock.now as now, stock.desc as desc, name, type, specification, unit from stock left outer join item on stock.item_id=item.id where stock.id=?`
}

func (*Stock) allStmt() string {
	return `select stock.id as id, item_id, warehouse, warehouse, stock.push as push, stock.pop as pop, stock.now as now, stock.desc as desc, name, type, specification, unit from stock left outer join item on stock.item_id=item.id`
}

func (*Stock) updateStmt() string {
	return "update stock"
}

func (*Stock) deleteStmt() string {
	return "delete from stock where id=?"
}
