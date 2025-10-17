package struct_to_map

import "testing"

func TestGen(t *testing.T) {
	GenStruct2MapFile("../unit_test/example.go")
}
