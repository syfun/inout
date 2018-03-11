package models

import (
	"strconv"
)

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

// Create ...
func (it *Item) Create(t Table) (Table, error) {
	rst, err := Create("insert into item (name, type, specification, unit, push, pop, now, desc) values (:name, :type, :specification, :unit, :push, :pop, :now, :desc)", t)
	if err != nil {
		return nil, err
	}
	id, _ := rst.LastInsertId()
	var item Item
	err = it.Get(&item, NewDBQuery(nil, map[string]string{"id": strconv.FormatInt(id, 10)}))
	return &item, err
}

var itemCols = stringSliceToMap("item", []string{
	"id", "name", "type", "specification", "unit", "push", "pop", "now", "desc",
})

// Update ..
func (it *Item) Update(query *DBQuery, data map[string]interface{}) (Table, error) {
	err := Update("update item set %s where id=:id", query, data, itemCols)
	if err != nil {
		return nil, err
	}
	var item Item
	err = it.Get(&item, query)
	return &item, err
}

// Get ...
func (*Item) Get(dest interface{}, query *DBQuery) error {
	return Get(dest, "select * from item where id=?", query)
}

// List ...
func (*Item) List(dest interface{}, query *DBQuery) (int64, error) {
	err := List(dest, "select * from item", query, itemCols)
	if err != nil {
		return 0, err
	}
	return Count("select count(*) from item", query, itemCols)
}

// Delete ...
func (*Item) Delete(query *DBQuery) error {
	return Delete("delete from item", query)
}
