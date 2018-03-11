package models

import "strconv"

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

// Create ...
func (s *Stock) Create(t Table) (Table, error) {
	rst, err := Create("insert into stock (item_id, warehouse, push, pop, now, desc) values (:item_id, :warehouse, :push, :pop, :now, :desc);", t)
	if err != nil {
		return nil, err
	}
	id, _ := rst.LastInsertId()
	var stock Stock
	err = s.Get(&stock, NewDBQuery(nil, map[string]string{"id": strconv.FormatInt(id, 10)}))
	return &stock, err
}

var stockCols = stringSliceToMap("stock", []string{
	"id", "item_id", "warehouse", "push", "pop", "now", "desc",
})

// Update ..
func (s *Stock) Update(query *DBQuery, data map[string]interface{}) (Table, error) {
	err := Update("update stock set %s where id=:id", query, data, stockCols)
	if err != nil {
		return nil, err
	}
	var stock Stock
	err = s.Get(&stock, query)
	return &stock, err
}

var stockSelect = `select stock.id, item_id, warehouse, warehouse, stock.push as push, stock.pop as pop, stock.now as now, stock.desc as desc, name, type, specification, unit from stock left outer join item on stock.item_id=item.id`

// Get ...
func (*Stock) Get(dest interface{}, query *DBQuery) error {
	return Get(dest, stockSelect+" where stock.id=?", query)
}

// List ...
func (*Stock) List(dest interface{}, query *DBQuery) (int64, error) {
	cols := MergeMap(itemCols, stockCols)
	err := List(dest, stockSelect, query, cols)
	if err != nil {
		return 0, err
	}
	return Count("select count(*) from stock left outer join item on stock.item_id=item.id", query, cols)
}

// Delete ...
func (*Stock) Delete(query *DBQuery) error {
	return Delete("delete from stock", query)
}
