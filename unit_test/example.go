package unit_test

import (
	"sync"

	"github.com/nan-www/convToMap/unit_test/aa"
)

//go:generate convToMap example.go
type Example struct {
	// 目前不支持不在同文件内的 inline 结构体
	aa.Hamabe    `json:",inline"`
	FooPtr       *Foo              `json:"fooPtr"`
	Foo          Foo               `json:"foo,omitempty"`
	ID           int               `json:"id,omitempty"`
	Name         string            `json:"name,omitempty"`
	Float        float64           `json:"float64,omitempty"`
	Ignore       map[string]string `json:"-"`
	PtrInt       *int64            `json:"ptrInt,omitempty"`
	InlineStruct `json:",inline"`
	Str2Foo      map[string]Foo `json:"str2Foo,omitempty"`
	Foos         []Foo          `json:"foos,omitempty"`
	Mu           sync.Mutex
}

//go:generate convToMap example.go
type InlineStruct struct {
	A string
	B int
}

//go:generate convToMap example.go
type Foo struct {
	Bar string `json:"bar"`
}
