package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

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

func (*Pop) relatedCreate(p *Pop) (sql.Result, error) {
	var (
		item  Item
		stock Stock
	)
	err := db.Get(&item, "select id, pop, now from item where id=?", p.ItemID)
	if err != nil {
		return nil, &DBError{err, "item not found"}
	}
	err = db.Get(&stock, "select id, pop, now from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		return nil, &DBError{err, "stock not found"}
	}
	stock.Pop += p.Number
	item.Pop += p.Number
	stock.Now -= p.Number
	item.Now -= p.Number
	tx := db.MustBegin()
	_, err = tx.Exec("update item set pop=?, now=? where id=?", item.Pop, item.Now, item.ID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, "update item pop error"}
	}
	_, err = tx.Exec("update stock set pop=?, now=? where id=?", stock.Pop, stock.Now, stock.ID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, "update stock pop error"}
	}
	rst, err := tx.NamedExec("insert into pop (item_id, time, number, warehouse, abstract, remark, receiver, checker) values (:item_id, :time, :number, :warehouse, :abstract, :remark, :receiver, :checker)", p)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, "insert error"}
	}
	tx.Commit()
	return rst, nil
}

func (*Pop) relatedUpdate(number int64, format string, query *DBQuery, data map[string]interface{}, colsMap map[string]string) error {
	var (
		item  Item
		stock Stock
		p     Pop
	)
	err := db.Get(&p, "select item_id, warehouse, number from pop where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "pop not found"}
	}
	err = db.Get(&item, "select id, pop, now from item where id=?", p.ItemID)
	if err != nil {
		return &DBError{err, "item not found"}
	}
	err = db.Get(&stock, "select id, pop, now from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		return &DBError{err, "stock not found"}
	}
	delta := number - p.Number
	item.Pop += delta
	stock.Pop += delta
	item.Now -= delta
	stock.Now -= delta

	tx := db.MustBegin()
	_, err = tx.Exec("update item set pop=?, now=? where id=?", item.Pop, item.Now, item.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update item pop error"}
	}
	_, err = tx.Exec("update stock set pop=?, now=? where id=?", stock.Pop, stock.Now, stock.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update stock pop error"}
	}
	namedVars := getNamedVars(data, colsMap)
	data["id"] = query.Get("id")
	_, err = tx.NamedExec(
		fmt.Sprintf(format, strings.Join(namedVars, ", ")), data)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update pop error"}
	}
	tx.Commit()
	return nil
}

// Create ...
func (p *Pop) Create(t Table) (Table, error) {
	var (
		rst sql.Result
		err error
	)
	tt := t.(*Pop)

	rst, err = p.relatedCreate(tt)
	if err != nil {
		return nil, err
	}

	id, _ := rst.LastInsertId()
	var pop Pop
	err = p.Get(&pop, NewDBQuery(nil, map[string]string{"id": strconv.FormatInt(id, 10)}))
	return &pop, err
}

var popCols = stringSliceToMap("pop", []string{
	"id", "item_id", "time", "number", "warehouse", "abstract", "remark", "receiver", "checker",
})

// Update ..
func (p *Pop) Update(query *DBQuery, data map[string]interface{}) (Table, error) {
	number, ok := data["number"]
	var err error
	if !ok {
		err = Update("update pop set %s where id=:id", query, data, popCols)
	} else {
		n := int64(number.(float64))
		err = p.relatedUpdate(n, "update pop set %s where id=:id", query, data, popCols)
	}

	if err != nil {
		return nil, err
	}
	var pop Pop
	err = p.Get(&pop, query)
	return &pop, err
}

var popSelect = `select pop.id as id, item_id, time, number, warehouse, abstract, remark, receiver, checker, name, type, specification, unit from pop left outer join item on pop.item_id=item.id`

// Get ...
func (*Pop) Get(dest interface{}, query *DBQuery) error {
	return Get(dest, popSelect+" where pop.id=?", query)
}

// List ...
func (*Pop) List(dest interface{}, query *DBQuery) (int64, error) {
	cols := MergeMap(itemCols, popCols)
	err := List(dest, popSelect, query, cols)
	if err != nil {
		return 0, err
	}
	return Count("select count(*) from pop left outer join item on pop.item_id=item.id", query, cols)
}

// Delete ...
func (*Pop) Delete(query *DBQuery) error {
	var (
		item  Item
		stock Stock
		p     Pop
	)
	err := db.Get(&p, "select id, item_id, warehouse, number from pop where id=?", query.Get("id"))
	if err != nil {
		return &DBError{err, "pop not found"}
	}
	err = db.Get(&item, "select id, pop, now from item where id=?", p.ItemID)
	if err != nil {
		return &DBError{err, "item not found"}
	}
	err = db.Get(&stock, "select id, pop, now from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		return &DBError{err, "stock not found"}
	}
	item.Pop -= p.Number
	stock.Pop -= p.Number
	item.Now += p.Number
	stock.Now += p.Number

	tx := db.MustBegin()
	_, err = tx.Exec("update item set pop=?, now=? where id=?", item.Pop, item.Now, item.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update item pop error"}
	}
	_, err = tx.Exec("update stock set pop=?, now=? where id=?", stock.Pop, stock.Now, stock.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "update stock pop error"}
	}
	_, err = tx.Exec("delete from pop where id=?", p.ID)
	if err != nil {
		tx.Rollback()
		return &DBError{err, "delete pop error"}
	}
	tx.Commit()
	return nil
}
