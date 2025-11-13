package generator

import (
	"bytes"
	"fmt"
	"github.com/nan-www/convToMap/ps"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

const suffix = "_generated_%d.go"

// Gen. command: e.g. //go:generate convToMap
func Gen(filename, command string, genTemplate ...string) {

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", filename, err)
		os.Exit(1)
	}
	data := ps.ParseMarkStruct(node, command)
	if len(data.Structs) == 0 {
		fmt.Fprintf(os.Stderr, "No struct found with %s", command)
		return
	}
	for idx, tp := range genTemplate {
		var buf bytes.Buffer
		t := template.Must(template.New("generator").Parse(tp))
		if err := t.Execute(&buf, data); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
			os.Exit(1)
		}

		formattedCode, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting generated code: %v\nCode:\n%s\n", err, buf.String())
			os.Exit(1)
		}
		suffixName := fmt.Sprintf(suffix, idx)
		// Output generate file.
		outputFilename := strings.TrimSuffix(filename, ".go") + suffixName
		if err := os.WriteFile(outputFilename, formattedCode, fs.FileMode(0644)); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully generated code for %d struct(s) in %s\n", len(data.Structs), outputFilename)
	}
}
