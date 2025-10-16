package unit_test

//go:generate convToMap
type Example struct {
	ID     int               `json:"id,omitempty"`
	Name   string            `json:"name,omitempty"`
	Float  float64           `json:"float64,omitempty"`
	Ignore map[string]string `json:"-"`
}
