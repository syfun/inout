package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
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

var schema = `
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
create index if not exists stock_item_index on stock(item_id);
`

func initTables() {
	_, err := db.Exec(schema)
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

func (dq *DBQuery) genNamedVars(cols map[string]string) ([]string, map[string]interface{}, string) {
	vars := make([]string, 0)
	bindMap := make(map[string]interface{})
	var (
		page     = 0
		pageSize = 10
		limit    = ""
	)
	for key := range dq.Values {
		if key == "page" {
			page, _ = strconv.Atoi(dq.Get("page"))
			continue
		}
		if key == "page_size" {
			pageSize, _ = strconv.Atoi(dq.Get("page_size"))
			continue
		}
		name, ok := cols[key]
		if !ok {
			continue
		}
		vars = append(vars, fmt.Sprintf("%s like :%s", name, key))
		bindMap[key] = "%" + dq.Get(key) + "%"
	}
	if page != 0 {
		limit = fmt.Sprintf(" limit %d offset %d", pageSize, (page-1)*pageSize)
	}
	return vars, bindMap, limit
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
	Create(Table) (Table, error)
	Update(*DBQuery, map[string]interface{}) (Table, error)
	Get(interface{}, *DBQuery) error
	List(interface{}, *DBQuery) (int64, error)
	Delete(*DBQuery) error
}

// Create create one row.
// var format := `insert into item (name, type, specification, unit, push, pop, now, desc)
//               values (:name, :type, :specification, :unit, :push, :pop, :now, :desc)`
// rst, err := Create(format, t)
func Create(format string, t Table) (sql.Result, error) {
	rst, err := db.NamedExec(format, t)
	if err != nil {
		return nil, &DBError{err, "insert error"}
	}
	return rst, nil
}

func getNamedVars(data map[string]interface{}, cols map[string]string) []string {
	index := 0
	namedVars := make([]string, 0)
	for key := range data {
		name, ok := cols[key]
		if !ok {
			continue
		}
		namedVars = append(namedVars, fmt.Sprintf("%s:=%s", name, key))
		index++
	}
	return namedVars
}

// Update rows.
// err := Update("update item set %s where id=:id", query, data, colsMap)
func Update(format string, query *DBQuery, data map[string]interface{}, colsMap map[string]string) error {
	namedVars := getNamedVars(data, colsMap)
	if len(namedVars) != 0 {
		data["id"] = query.Get("id")
		_, err := db.NamedExec(
			fmt.Sprintf(format, strings.Join(namedVars, ", ")), data)
		if err != nil {
			return &DBError{err, "update error"}
		}
	}
	return nil
}

// Get select one row.
// err := Get(dest, "select * from item", query)
func Get(dest interface{}, format string, query *DBQuery) error {
	err := db.Get(dest, format, query.Get("id"))
	if err != nil {
		return &DBError{err, "get error"}
	}
	return nil
}

// List select multi rows
// err := List(dest, "select * from item", query)
func List(dest interface{}, format string, query *DBQuery, cols map[string]string) error {
	var err error
	vars, bindMap, limit := query.genNamedVars(cols)
	if len(vars) == 0 {
		format += limit
		err = db.Select(dest, format)
	} else {
		stmt, _ := db.PrepareNamed(
			format + " where " + strings.Join(vars, " and ") + limit)
		err = stmt.Select(dest, bindMap)
	}
	if err != nil {
		return &DBError{err, "list error"}
	}
	return nil
}

// Count get table count
func Count(format string, query *DBQuery, cols map[string]string) (int64, error) {
	var (
		count int64
		err   error
	)
	vars, bindMap, _ := query.genNamedVars(cols)
	if len(vars) == 0 {
		err = db.Get(&count, format)
	} else {
		format += " where " + strings.Join(vars, " and ")
		prepared, _ := db.PrepareNamed(format)
		err = prepared.Get(&count, bindMap)
	}
	if err != nil {
		return 0, &DBError{err, "cannot get table count"}
	}
	return count, nil
}

// Delete delete row.
// err := Delete("delete from item", query)
func Delete(format string, query *DBQuery) error {
	_, err := db.Exec(format+" where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "delete error"}
	}
	return nil
}
