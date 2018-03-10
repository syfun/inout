package models

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Remove("./for_test.db")
	InitDB("./for_test.db")
	res := m.Run()
	CloseDB()
	os.Remove("./for_test.db")
	os.Exit(res)
}

func TestModelInsert(t *testing.T) {
	m := &Model{&Item{}}
	it := &Item{Name: "玻璃", Type: "有机玻璃"}
	rst, err := m.Insert(it)
	if err != nil {
		t.Error(err)
	}
	t.Log(rst)
}
