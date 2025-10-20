package print_fields

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/nan-www/convToMap/ps"
)

func PrintAllFields(filename, structName string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", filename, err)
		os.Exit(1)
	}
	data := ps.ParseMarkStruct(node)
	if len(data.Structs) == 0 {
		return
	}
	//for _, s := range data.Structs {
	//	if s.Name != structName {
	//		continue
	//	}
	//	s.Fields
	//}
}
