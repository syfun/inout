package models

import (
	"time"
)

// Push ...
type Push struct {
	ID        int64     `json:"id" column:"id"`
	ItemID    int64     `json:"name" column:"item_id"`
	Time      time.Time `json:"time" column:"time"`
	Number    int64     `json:"number" column:"number"`
	Warehouse string    `json:"warehouse" column:"warehouse"`
	Abstract  string    `json:"abstract" column:"abstract"`
	Remark    string    `json:"remark" column:"remark"`
	User      string    `json:"user" column:"user"`
}

// TName get Push table name.
func (p *Push) TName() string {
	return "push"
}
