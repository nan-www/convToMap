package print_fields

import "testing"

func TestPrintAllFields(t *testing.T) {
	PrintAllFields("../unit_test/example.go", "Example")
}
