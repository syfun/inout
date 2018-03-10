package models

import (
	"database/sql"
	"fmt"
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

func (*Push) columns() map[string]interface{} {
	return stringSliceToMap([]string{"id", "item_id", "time", "number", "warehouse", "abstract", "remark", "user"})
}

// Insert ...
func (*Push) Insert(p *Push) (sql.Result, error) {
	tx := db.MustBegin()
	var (
		item  Item
		stock Stock
	)
	err := tx.Get(&item, "select * from item where id=?", p.ItemID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, ""}
	}
	err = tx.Get(&stock, "select * from stock where item_id=? and warehouse=?", p.ItemID, p.Warehouse)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, ""}
	}
	stock.Push += p.Number
	item.Push += p.Number
	_, err = tx.Exec("update item set push=? where id=?", item.Push, item.ID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, ""}
	}
	_, err = tx.Exec("update stock set push=? where id=?", stock.Push, stock.ID)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, ""}
	}
	rst, err := tx.NamedExec("insert into push (item_id, time, number, warehouse, abstract, remark, user) values (:item_id, :time, :number, :warehouse, :abstract, :remark, :user)", p)
	if err != nil {
		tx.Rollback()
		return nil, &DBError{err, ""}
	}
	tx.Commit()

	return rst, err
}

// Update ...
func (p *Push) Update(query *DBQuery, data map[string]interface{}) (interface{}, error) {
	index := 0
	namedVars := make([]string, 0)
	colsMap := p.columns()
	for key := range data {
		if _, ok := colsMap[key]; !ok {
			continue
		}
		namedVars = append(namedVars, key+"=:"+key)
		index++
	}
	tx := db.MustBegin()
	if len(namedVars) != 0 {
		data["id"] = query.Get("id")
		stmt := fmt.Sprintf("update push set %s where id=:id", strings.Join(namedVars, ", "))
		_, err := db.NamedExec(stmt, data)
		if err != nil {
			tx.Rollback()
			return nil, &DBError{err, "update error"}
		}
	}
	tx.Commit()
}

func (*Push) getStmt() string {
	return `select push.id as id, item_id, time, number, warehouse, abstract, remark, user, name, type, specification, unit from push left outer join item on push.item_id=item.id where push.id=?`
}

func (*Push) allStmt() string {
	return `select push.id as id, item_id, time, number, warehouse, abstract, remark, user, name, type, specification, unit from push left outer join item on push.item_id=item.id`
}

func (*Push) updateStmt() string {
	return "update push"
}

func (*Push) deleteStmt() string {
	return "delete from push where id=?"
}
