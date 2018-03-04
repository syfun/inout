package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

var db *sql.DB

// InitDB init database connection.
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// initTables()
}

// CloseDB expose close db.
func CloseDB() {
	db.Close()
}

func initTables() {
	_, err := db.Exec(`
		create table if not exists label (
			id integer not null primary key,
			name text, type text);
		create table if not exists item (
			id integer not null primary key,
			name text, type text, specification text, unit text,
			push integer, pop integer, now integer, desc text);
		create table if not exists push (
			id integer not null primary key,
			item_id integer references item,
			time datetime, number integer, warehouse text, abstract text, remark text, user text);
		create table if not exists pop (
			id integer not null primary key,
			item_id integer references item,
			time datetime, number integer, receiver text,
			checker text, warehouse text, abstract text, remark text);
		create table if not exists stock  (
			id integer not null primary key,
			item_id integer references item,
			warehouse text, push integer, pop integer, now integer, desc text);
		create index push_item_index on push(item_id);
		create index pop_item_index on pop(item_id);
		create index stock_item_index on stock(item_id);`)
	if err != nil {
		log.Fatal(err)
	}
}

// DBError for db error when exec.
type DBError struct {
	error   error
	Message string
}

func (dr DBError) Error() string {
	return fmt.Sprintf("%s: %s", dr.Message, dr.error.Error())
}

// DBQuery for db query.
type DBQuery struct {
	url.Values
}

// Length return DBQuery length.
func (dq *DBQuery) Length() int {
	return len(dq.Values)
}

// NewDBQuery create new DBQuery.
func NewDBQuery(values url.Values, kv map[string]string) *DBQuery {
	if values == nil {
		values = make(url.Values)
	}
	query := &DBQuery{values}
	if kv != nil {
		for k, v := range kv {
			query.Set(k, v)
		}
	}
	return query
}

// Table interface.
type Table interface {
	// TName for table name
	TName() string
}

// Model struct.
type Model struct {
	Table Table
}

// Insert insert into table.
func (m *Model) Insert(resource interface{}) (interface{}, error) {
	val, _ := resource.(reflect.Value)
	typ := val.Type()
	num := typ.NumField()
	cols := make([]string, num-1)
	vals := make([]interface{}, num-1)
	index := 0

	for i := 0; i < num; i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("column")
		if field.Name == "ID" || tag == "" {
			continue
		}

		cols[index] = tag
		vals[index] = val.Field(i).Interface()
		index++
	}

	stmt := fmt.Sprintf(
		"insert into %s (%s) values (%s);",
		m.Table.TName(),
		strings.Join(cols, ","),
		strings.TrimRight(strings.Repeat("?, ", len(cols)), ", "),
	)
	res, err := db.Exec(stmt, vals...)

	if err != nil {
		return nil, &DBError{err, fmt.Sprintf("insert %s error", m.Table)}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, &DBError{err, "cann get last insert id"}
	}
	val.FieldByName("ID").SetInt(id)
	return val.Interface(), nil
}

// Get query one from table.
func (m *Model) Get(query *DBQuery) (interface{}, error) {
	table := m.Table.TName()
	stmt := fmt.Sprintf("select * from %s where id=?", table)
	rows, err := db.Query(stmt, query.Get("id"))
	if err != nil {
		return nil, &DBError{err, "query error from " + table}
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	dest := make([]interface{}, len(cols))
	colsMap := make(map[string]int)
	for i, c := range cols {
		colsMap[c] = i
	}
	resource := reflect.New(reflect.TypeOf(m.Table).Elem()).Elem()
	typ := resource.Type()
	for i := 0; i < typ.NumField(); i++ {
		colName := typ.Field(i).Tag.Get("column")
		colIndex, ok := colsMap[colName]
		if !ok {
			return nil, fmt.Errorf("%s has no column %s", table, colName)
		}
		dest[colIndex] = resource.Field(i).Addr().Interface()
	}
	if !rows.Next() {
		return nil, fmt.Errorf("%s not found", table)
	}
	err = rows.Scan(dest...)
	if err != nil {
		return nil, &DBError{err, "scan error"}
	}
	return resource.Interface(), nil
}

// All query all from table.
func (m *Model) All(query *DBQuery) (interface{}, error) {
	var (
		rows   *sql.Rows
		err    error
		length = query.Length()
		table  = m.Table.TName()
	)
	if length == 0 {
		rows, err = db.Query(fmt.Sprintf("select * from %s;", table))
	} else {
		vals := make([]interface{}, length)
		cols := make([]string, length)
		index := 0
		for key := range query.Values {
			cols[index] = key + "=?"
			vals[index] = query.Get(key)
			index++
		}
		stmt := fmt.Sprintf(
			"select * from %s where %s;",
			table, strings.Join(cols, " and "),
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
	typ := reflect.TypeOf(m.Table).Elem()
	for i := 0; i < typ.NumField(); i++ {
		colName := typ.Field(i).Tag.Get("column")
		colIndex, ok := colsMap[colName]
		if !ok {
			return nil, fmt.Errorf("%s has no column %s", table, colName)
		}
		indexes[i] = colIndex
	}

	resources := make([]interface{}, 0)
	for rows.Next() {
		resource := reflect.New(reflect.TypeOf(m.Table).Elem()).Elem()
		dest := make([]interface{}, len(cols))
		for i, index := range indexes {
			dest[index] = resource.Field(i).Addr().Interface()
		}
		if err = rows.Scan(dest...); err != nil {
			return nil, &DBError{err, "scan error"}
		}
		resources = append(resources, resource.Interface())
	}
	return resources, nil
}

// Update update table.
func (m *Model) Update(query *DBQuery, data map[string]interface{}) (interface{}, error) {
	colsMap := make(map[string]bool)
	typ := reflect.TypeOf(m.Table).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("column")
		if field.Name == "ID" || tag == "" {
			continue
		}
		colsMap[tag] = true
	}
	cols := make([]string, len(data))
	vals := make([]interface{}, len(data)+1)
	index := 0
	for k, v := range data {
		cols[index] = k + "=?"
		vals[index] = v
		index++
	}
	vals[index] = query.Get("id")
	table := m.Table.TName()
	stmt := fmt.Sprintf(
		"update %s set %s where id=?",
		table, strings.Join(cols, ", "),
	)
	_, err := db.Exec(stmt, vals...)
	if err != nil {
		return nil, &DBError{err, "cannot update " + table}
	}
	resource, _ := m.Get(query)
	return resource, nil
}

// Delete delete from table.
func (m *Model) Delete(query *DBQuery) error {
	table := m.Table.TName()
	_, err := db.Exec(fmt.Sprintf("delete from %s where id=?", table), query.Get("id"))
	if err != nil {
		return &DBError{err, "cannot delete from " + table}
	}
	return nil
}
