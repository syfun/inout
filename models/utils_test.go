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
