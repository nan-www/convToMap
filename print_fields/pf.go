package print_fields

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/bytedance/sonic"
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
	ret := make([]string, 0)
	for _, s := range data.Structs {
		if s.Name != structName {
			continue
		}
		for _, field := range s.Fields {
			ret = append(ret, field.TagName)
		}
	}
	marshalString, _ := sonic.MarshalString(ret)
	fmt.Println(marshalString)
}
