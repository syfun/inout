package models

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Item ...
type Item struct {
	ID            int64  `json:"id" column:"id"`
	Name          string `json:"name" column:"name"`
	Type          string `json:"type" column:"type"`
	Specification string `json:"specification" column:"specification"`
	Unit          string `json:"unit" column:"unit"`
	Push          int64  `json:"push" column:"push"`
	Pop           int64  `json:"pop" column:"pop"`
	Now           int64  `json:"now" column:"now"`
	Desc          string `json:"desc" column:"desc"`
}

// DefaultItem for Get and All.
var DefaultItem = &Item{}

// Get get item by id.
func (*Item) Get(query *DBQuery) (interface{}, error) {
	item := &Item{}
	rows, err := db.Query("select * from item where id=?", query.Get("id"))
	if err != nil {
		return nil, &DBError{err, "cannot get item"}
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	vals := make([]interface{}, len(cols))
	colsMap := make(map[string]int)
	for i, c := range cols {
		colsMap[c] = i
	}
	itemValue := reflect.ValueOf(item).Elem()
	itemType := itemValue.Type()
	for i := 0; i < itemType.NumField(); i++ {
		colName := itemType.Field(i).Tag.Get("column")
		colIndex, ok := colsMap[colName]
		if !ok {
			return nil, fmt.Errorf("item has no column %s", colName)
		}
		vals[colIndex] = itemValue.Field(i).Addr().Interface()
	}
	if !rows.Next() {
		return nil, errors.New("item not found")
	}
	err = rows.Scan(vals...)
	if err != nil {
		return nil, &DBError{err, "cannot get item"}
	}
	return item, nil
}

// GetItem wrap Get.
func GetItem(query *DBQuery) (interface{}, error) {
	return DefaultItem.Get(query)
}

// All query items
func (it *Item) All(query *DBQuery) (interface{}, error) {
	var (
		rows   *sql.Rows
		err    error
		length = query.Length()
	)
	if length == 0 {
		rows, err = db.Query("select * from item;")
	} else {
		vals := make([]interface{}, length)
		cols := make([]string, length)
		// stmt := "select * from item where "
		index := 0
		for key := range query.Values {
			cols[index] = key + "=?"
			vals[index] = query.Get(key)
			index++
		}
		stmt := fmt.Sprintf(
			"select * from item where %s;",
			strings.Join(cols, " and "),
		)
		rows, err = db.Query(stmt, vals...)
	}

	if err != nil {
		return nil, &DBError{err, "cannot query item"}
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	colsMap := make(map[string]int)
	for i, c := range cols {
		colsMap[c] = i
	}
	indexes := make([]int, len(cols))
	itemType := reflect.ValueOf(it).Elem().Type()
	for i := 0; i < itemType.NumField(); i++ {
		colName := itemType.Field(i).Tag.Get("column")
		colIndex, ok := colsMap[colName]
		if !ok {
			return nil, fmt.Errorf("item has no column %s", colName)
		}
		indexes[i] = colIndex
	}

	items := make([]*Item, 0)
	for rows.Next() {
		item := &Item{}
		vals := make([]interface{}, len(cols))
		itemValue := reflect.ValueOf(item).Elem()
		for i, index := range indexes {
			vals[index] = itemValue.Field(i).Addr().Interface()
		}
		if err = rows.Scan(vals...); err != nil {
			return nil, &DBError{err, "scan error"}
		}
		items = append(items, item)
	}
	return items, nil
}

// Insert insert item.
func (it *Item) Insert() error {
	itemValue := reflect.ValueOf(it).Elem()
	itemType := itemValue.Type()
	num := itemType.NumField()
	cols := make([]string, num-1)
	vals := make([]interface{}, num-1)
	index := 0
	for i := 0; i < num; i++ {
		field := itemType.Field(i)
		if field.Name == "ID" {
			continue
		}
		cols[index] = field.Tag.Get("column")
		vals[index] = itemValue.Field(i).Interface()
		index++
	}

	stmt := fmt.Sprintf(
		"insert into item (%s) values (%s);",
		strings.Join(cols, ","),
		strings.TrimRight(strings.Repeat("?, ", len(cols)), ", "),
	)
	res, err := db.Exec(stmt, vals...)
	if err != nil {
		return &DBError{err, "insert item error"}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return &DBError{err, "cann get last insert id"}
	}
	it.ID = id
	return nil
}

// Update update item name.
func (it *Item) Update() error {
	_, err := db.Exec("update label set name=? where id=?", lb.Name, lb.ID)
	if err != nil {
		return &DBError{err, "cannot update label"}
	}
	return nil
}

// Delete delete item.
func (it *Item) Delete(query *DBQuery) error {
	_, err := db.Exec("delete from item where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "cannot delete item"}
	}
	return nil
}
