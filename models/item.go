package models

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

// TName get item table name.
func (it *Item) TName() string {
	return "item"
}
