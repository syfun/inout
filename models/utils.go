package models

import "reflect"

func stringSliceToMap(value string, slice []string) map[string]string {
	m := map[string]string{}
	for _, s := range slice {
		m[s] = value + "." + s
	}
	return m
}

// MergeMap merge maps.
func MergeMap(m1, m2 map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range m1 {
		m[k] = v
	}

	for k, v := range m2 {
		m[k] = v
	}
	return m
}

// MapToStruct convert map m to struct s.
func MapToStruct(m map[string]interface{}, structs ...interface{}) {
	var (
		val reflect.Value
		typ reflect.Type
	)
	for _, s := range structs {
		val = reflect.ValueOf(s).Elem()
		typ = val.Type()
		for i := 0; i < typ.NumField(); i++ {
			f := val.Field(i)
			name := typ.Field(i).Name
			v, ok := m[name]
			if ok && f.CanSet() {
				f.Set(reflect.ValueOf(v))
			}
		}
	}
}
