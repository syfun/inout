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
	label, err := CreateLabel("test", "user")
	if err != nil {
		t.Error(err)
	}
	t.Log(label)
}

func TestGetLabels(t *testing.T) {
	_, err := CreateLabel("test1", "user")
	if err != nil {
		t.Error(err)
	}
	_, err = CreateLabel("test2", "warehouse")
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
	label, err := CreateLabel("test1", "user")
	if err != nil {
		t.Error(err)
	}
	_, err = UpdateLabel(label.ID, "updated")
	if err != nil {
		t.Error(err)
	}
}
