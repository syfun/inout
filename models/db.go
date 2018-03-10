package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

var db *sqlx.DB

// InitDB init database connection.
func InitDB(dataSourceName string) {
	db = sqlx.MustConnect("sqlite3", dataSourceName)
	initTables()
}

// CloseDB expose close db.
func CloseDB() {
	db.Close()
}

func initTables() {
	_, err := db.Exec(`
		create table if not exists label (
			id integer not null primary key,
			name text, type text
		);
		create table if not exists item (
			id integer not null primary key,
			name text, type text, specification text, unit text,
			push integer, pop integer, now integer, desc text
		);
		create table if not exists push (
			id integer not null primary key,
			item_id integer references item,
			time text, number integer, warehouse text, abstract text, remark text, user text);
		create table if not exists pop (
			id integer not null primary key,
			item_id integer references item,
			time text, number integer, receiver text,
			checker text, warehouse text, abstract text, remark text);
		create table if not exists stock  (
			id integer not null primary key,
			item_id integer references item,
			warehouse text, push integer, pop integer, now integer, desc text);
		create unique index if not exists item_name_unique on item(name);
		create unique index if not exists label_name_type_unique on label(name, type);
		create index if not exists push_item_index on push(item_id);
		create index if not exists pop_item_index on pop(item_id);
		create index if not exists stock_item_index on stock(item_id);`)
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

func (dq *DBQuery) genNamedVars() ([]string, map[string]interface{}) {
	index := 0
	vars := make([]string, dq.Length())
	bindMap := make(map[string]interface{})
	for key := range dq.Values {
		vars[index] = key + "=:" + key
		bindMap[key] = dq.Get(key)
		index++
	}
	return vars, bindMap
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
	insertStmt() string
	getStmt() string
	allStmt() string
	updateStmt() string
	deleteStmt() string
	columns() map[string]interface{}
}

// Model struct.
type Model struct {
	Table Table
}

// Insert insert into table.
func (m *Model) Insert(resource interface{}) (interface{}, error) {
	var (
		result sql.Result
		err    error
	)
	if _, ok := reflect.TypeOf(m.Table).MethodByName("Insert"); !ok {
		result, err = db.NamedExec(m.Table.insertStmt(), resource)
		if err != nil {
			return nil, &DBError{err, "insert error"}
		}
	} else {
		rst := reflect.ValueOf(m.Table).MethodByName("Insert").Call(
			[]reflect.Value{reflect.ValueOf(resource)},
		)
		if !rst[1].IsNil() {
			return nil, &DBError{rst[1].Interface().(error), "insert error"}
		}
		result = rst[0].Interface().(sql.Result)
	}

	id, _ := result.LastInsertId()
	// reflect.ValueOf(resource).Elem().FieldByName("ID").SetInt(id)
	resource, _ = m.Get(NewDBQuery(nil, map[string]string{"id": strconv.FormatInt(id, 10)}))
	return resource, nil
}

// Get query one from table.
func (m *Model) Get(query *DBQuery) (interface{}, error) {
	resource := reflect.New(reflect.TypeOf(m.Table).Elem()).Elem().Addr()
	err := db.Get(resource.Interface(), m.Table.getStmt(), query.Get("id"))
	if err != nil {
		return nil, &DBError{err, "query single error"}
	}
	return resource.Interface(), nil
}

// All query all from table.
func (m *Model) All(query *DBQuery) (interface{}, error) {
	var (
		err    error
		length = query.Length()
	)
	sliceType := reflect.SliceOf(reflect.TypeOf(m.Table))
	resources := reflect.New(sliceType)

	stmt := m.Table.allStmt()
	if length == 0 {
		err = db.Select(resources.Interface(), stmt)
	} else {
		vars, bindMap := query.genNamedVars()
		stmt = fmt.Sprintf("%s where %s", stmt, strings.Join(vars, " and "))
		prepared, _ := db.PrepareNamed(stmt)
		err = prepared.Select(resources.Interface(), bindMap)
	}

	if err != nil {
		return nil, &DBError{err, "query all error"}
	}
	return resources.Interface(), nil
}

// Update update table.
func (m *Model) Update(query *DBQuery, data map[string]interface{}) (interface{}, error) {
	var resource interface{}
	if _, ok := reflect.TypeOf(m.Table).MethodByName("Update"); !ok {
		index := 0
		namedVars := make([]string, 0)
		colsMap := m.Table.columns()
		for key := range data {
			if _, ok := colsMap[key]; !ok {
				continue
			}
			namedVars = append(namedVars, key+"=:"+key)
			index++
		}
		if len(namedVars) != 0 {
			data["id"] = query.Get("id")
			stmt := fmt.Sprintf("%s set %s where id=:id",
				m.Table.updateStmt(), strings.Join(namedVars, ", "))
			_, err := db.NamedExec(stmt, data)
			if err != nil {
				return nil, &DBError{err, "update error"}
			}
		}
		resource, _ = m.Get(query)
	} else {
		rst := reflect.ValueOf(m.Table).MethodByName("Update").Call(
			[]reflect.Value{reflect.ValueOf(query), reflect.ValueOf(data)},
		)
		if !rst[1].IsNil() {
			return nil, &DBError{rst[1].Interface().(error), "update error"}
		}
		resource = rst[0].Interface()
	}

	return resource, nil
}

// Delete delete from table.
func (m *Model) Delete(query *DBQuery) error {
	_, err := db.Exec(m.Table.deleteStmt(), query.Get("id"))
	if err != nil {
		return &DBError{err, "delete error"}
	}
	return nil
}
