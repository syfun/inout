package main

import (
	"testing"
)

type Person struct {
	name string
}

func TestResponseWrite(t *testing.T) {
	resp := Response{Person{"sdfsd"}, 200}
	resp.Write()
}
