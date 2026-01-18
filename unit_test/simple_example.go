package unit_test

//go:generate convToMap simple_example.go
type SimpleExample struct {
	Str   string `json:"str"`
	Point Point  `json:"point"`
	NMIXX `json:"inline"`
}

//go:generate convToMap simple_example.go
type Point struct {
	X int     `json:"x"`
	Y float64 `json:"y"`
}

//go:generate convToMap simple_example.go
type NMIXX struct {
	K *string `json:"k"`
}
