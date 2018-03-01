package models

import (
	"database/sql"
	"fmt"
	"log"

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
			name text, type text, specifaction text, unit text,
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
