package unit_test

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/nan-www/convToMap/unit_test/aa"
)

func TestMap2Struct(t *testing.T) {
	m := map[string]any{
		"float64": 2.9,
		"id":      99,
		"name":    "name",
		"A":       "a1",
		"B":       2,
	}
	ex := &Example{}
	ex._2Struct(m)
	marshalString, _ := sonic.MarshalString(ex)
	t.Logf("ex: %v", marshalString)

	example := &Example{
		Hamabe: aa.Hamabe{
			Minami: "minami",
		},
		FooPtr: &Foo{
			Bar: "*bar",
		},
		Foo: Foo{
			Bar: "bar",
		},
		ID:    1,
		Name:  "name",
		Float: 1.0,
	}
	m = example._2Map()
	t.Logf("m: %v", m)
	ex1 := &Example{}
	ex1._2Struct(m)
	marshalString, _ = sonic.MarshalString(ex1)
	t.Logf("ex1: %v", marshalString)
}
