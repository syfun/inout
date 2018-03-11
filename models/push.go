package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// Push ...
type Push struct {
	ID            int64  `json:"id"`
	ItemID        int64  `json:"item_id" db:"item_id"`
	Time          string `json:"time"`
	Number        int64  `json:"number"`
	Warehouse     string `json:"warehouse"`
	Abstract      string `json:"abstract"`
	Remark        string `json:"remark"`
	User          string `json:"user"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Specification string `json:"specification"`
	Unit          string `json:"unit"`
}

func (*Push) relatedCreate(p *Push) (sql.Result, error) {
	var (
		item  Item
		stock Stock
	)
	err := db.Get(&item, "select id, push, now from item where id=?", p.ItemID)
	if err != nil {
		return nil, &DBError{err, "item not found"}
	}
	err = db.Get(&stock, "select id, push, now from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		return nil, &DBError{err, "stock not found"}
	}
	stock.Push += p.Number
	item.Push += p.Number
	stock.Now += p.Number
	item.Now += p.Number

	tx := db.MustBegin()
	_, err = tx.Exec("update item set push=?, now=? where id=?", item.Push, item.Now, item.ID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, "update item push error"}
	}
	_, err = tx.Exec("update stock set push=?, now=? where id=?", stock.Push, stock.Now, stock.ID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, "update stock push error"}
	}
	rst, err := tx.NamedExec("insert into push (item_id, time, number, warehouse, abstract, remark, user) values (:item_id, :time, :number, :warehouse, :abstract, :remark, :user)", p)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, "insert error"}
	}
	tx.Commit()
	return rst, nil
}

func (*Push) relatedUpdate(number int64, format string, query *DBQuery, data map[string]interface{}, colsMap map[string]string) error {
	var (
		item  Item
		stock Stock
		p     Push
	)
	err := db.Get(&p, "select item_id, warehouse, number from push where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "push not found"}
	}
	err = db.Get(&item, "select id, push, now from item where id=?", p.ItemID)
	if err != nil {
		return &DBError{err, "item not found"}
	}
	err = db.Get(&stock, "select id, push, now from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		return &DBError{err, "stock not found"}
	}
	delta := number - p.Number
	item.Push += delta
	stock.Push += delta
	item.Now += delta
	stock.Now += delta

	tx := db.MustBegin()
	_, err = tx.Exec("update item set push=?, now=? where id=?", item.Push, item.Now, item.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update item push error"}
	}
	_, err = tx.Exec("update stock set push=?, now=? where id=?", stock.Push, stock.Now, stock.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update stock push error"}
	}
	namedVars := getNamedVars(data, colsMap)
	data["id"] = query.Get("id")
	_, err = tx.NamedExec(
		fmt.Sprintf(format, strings.Join(namedVars, ", ")), data)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update push error"}
	}
	tx.Commit()
	return nil
}

// Create ...
func (p *Push) Create(t Table) (Table, error) {
	var (
		rst sql.Result
		err error
	)
	tt := t.(*Push)
	rst, err = p.relatedCreate(tt)
	if err != nil {
		return nil, err
	}

	id, _ := rst.LastInsertId()
	var push Push
	err = p.Get(&push, NewDBQuery(nil, map[string]string{"id": strconv.FormatInt(id, 10)}))
	return &push, err
}

var pushCols = stringSliceToMap("push", []string{
	"id", "item_id", "time", "number", "warehouse", "abstract", "remark", "user",
})

// Update ..
func (p *Push) Update(query *DBQuery, data map[string]interface{}) (Table, error) {
	number, ok := data["number"]
	var err error
	if !ok {
		err = Update("update push set %s where id=:id", query, data, pushCols)
	} else {
		n := int64(number.(float64))
		err = p.relatedUpdate(n, "update push set %s where id=:id", query, data, pushCols)
	}

	if err != nil {
		return nil, err
	}
	var push Push
	err = p.Get(&push, query)
	return &push, err
}

var selectStmt = `select push.id as id, item_id, time, number, warehouse, abstract, remark, user, name, type, specification, unit from push left outer join item on push.item_id=item.id`

// Get ...
func (*Push) Get(dest interface{}, query *DBQuery) error {
	return Get(dest, selectStmt+" where push.id=?", query)
}

// List ...
func (*Push) List(dest interface{}, query *DBQuery) (int64, error) {
	cols := MergeMap(itemCols, pushCols)
	err := List(dest, selectStmt, query, cols)
	if err != nil {
		return 0, err
	}
	return Count("select count(*) from push left outer join item on push.item_id=item.id", query, cols)
}

// Delete ...
func (*Push) Delete(query *DBQuery) error {
	var (
		item  Item
		stock Stock
		p     Push
	)
	err := db.Get(&p, "select id, item_id, warehouse, number from push where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "push not found"}
	}
	err = db.Get(&item, "select id, push, now from item where id=?", p.ItemID)
	if err != nil {
		return &DBError{err, "item not found"}
	}
	err = db.Get(&stock, "select id, push, now from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		return &DBError{err, "stock not found"}
	}
	item.Push -= p.Number
	stock.Push -= p.Number
	item.Now -= p.Number
	stock.Now -= p.Number

	tx := db.MustBegin()
	_, err = tx.Exec("update item set push=?, now=? where id=?", item.Push, item.Now, item.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update item push error"}
	}
	_, err = tx.Exec("update stock set push=?, now=? where id=?", stock.Push, stock.Now, stock.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update stock push error"}
	}
	_, err = tx.Exec("delete from push where id=?", p.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "delete push error"}
	}
	tx.Commit()
	return nil
}
