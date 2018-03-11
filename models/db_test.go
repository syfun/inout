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
