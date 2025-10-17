package unit_test

import "github.com/nan-www/convToMap/unit_test/aa"

//go:generate convToMap example.go
type Example struct {
	ID           int               `json:"id,omitempty"`
	Name         string            `json:"name,omitempty"`
	Float        float64           `json:"float64,omitempty"`
	Ignore       map[string]string `json:"-"`
	InlineStruct `json:",inline"`
	aa.Hamabe    `json:",inline"`
}

type InlineStruct struct {
	A string
	B int
}
