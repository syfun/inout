package models

import (
	"log"
	"reflect"
	"testing"
)

func TestInsert(t *testing.T) {
	p := &Push{}
	m, ok := reflect.TypeOf(p).MethodByName("Insert")
	log.Println(m, ok)
	log.Println(
		reflect.TypeOf(p).NumMethod(),
		reflect.TypeOf(p).Elem(),
		reflect.TypeOf(p).Elem().NumMethod(),
		reflect.TypeOf(p),
	)
}
