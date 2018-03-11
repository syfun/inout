package models

import "strconv"

// Label can be user, warehouse, item type
type Label struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	// Type can be user, warehouse, item type
	Type string `json:"type"`
}

// Create ...
func (lb *Label) Create(t Table) (Table, error) {
	rst, err := Create("insert into label (name, type) values (:name, :type)", t)
	if err != nil {
		return nil, err
	}
	id, _ := rst.LastInsertId()
	var label Label
	err = lb.Get(&label, NewDBQuery(nil, map[string]string{"id": strconv.FormatInt(id, 10)}))
	return &label, err
}

var labelCols = stringSliceToMap("label", []string{"id", "name", "type"})

// Update ..
func (lb *Label) Update(query *DBQuery, data map[string]interface{}) (Table, error) {
	err := Update("update label set %s where id=:id", query, data, labelCols)
	if err != nil {
		return nil, err
	}
	var label Label
	err = lb.Get(&label, query)
	return &label, err
}

// Get ...
func (*Label) Get(dest interface{}, query *DBQuery) error {
	return Get(dest, "select * from label where id=?", query)
}

// List ...
func (*Label) List(dest interface{}, query *DBQuery) (int64, error) {
	err := List(dest, "select * from label", query, labelCols)
	if err != nil {
		return 0, err
	}
	return Count("select count(*) from label", query, labelCols)
}

// Delete ...
func (*Label) Delete(query *DBQuery) error {
	return Delete("delete from label", query)
}
