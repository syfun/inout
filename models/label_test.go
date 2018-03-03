package models

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Remove("./for_test.db")
	InitDB("./for_test.db")
	initTables()
	res := m.Run()
	os.Remove("./for_test.db")
	os.Exit(res)
}
func TestCreateLabel(t *testing.T) {
	label, err := Insert("test", "user")
	if err != nil {
		t.Error(err)
	}
	t.Log(label)
}

func TestGetLabels(t *testing.T) {
	_, err := Insert("test1", "user")
	if err != nil {
		t.Error(err)
	}
	_, err = Insert("test2", "warehouse")
	if err != nil {
		t.Error(err)
	}
	labels, err := GetLabels("user")
	if err != nil {
		t.Error(err)
	}
	t.Log("user type", labels)

	labels, err = GetLabels("")
	if err != nil {
		t.Error(err)
	}
	t.Log("all", labels)
}

func TestUpdateLabel(t *testing.T) {
	label, err := Insert("test1", "user")
	if err != nil {
		t.Error(err)
	}
	_, err = Update(label.ID, "updated")
	if err != nil {
		t.Error(err)
	}
}
