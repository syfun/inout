package models

import (
	"testing"
)

func TestStringSliceToMap(t *testing.T) {
	m := stringSliceToMap([]string{"1", "2"})
	if len(m) != 2 {
		t.Error("slice to map error")
	}
	t.Log(m)
}

func TestMapToStruct(t *testing.T) {
	var s struct {
		Name string
		Age  int64
	}
	m := map[string]interface{}{
		"Name": "Jack",
		"age":  29,
	}
	MapToStruct(m, &s)
	t.Log(s)
}
