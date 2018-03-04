package models

// Label can be user, warehouse, item type
type Label struct {
	ID   int64  `json:"id" column:"id"`
	Name string `json:"name" column:"name"`

	// Type can be user, warehouse, item type
	Type string `json:"type" column:"type"`
}

// TName ...
func (lb *Label) TName() string {
	return "label"
}
