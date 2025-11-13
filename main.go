package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nan-www/convToMap/generator"
	"github.com/nan-www/convToMap/map_to_struct"
	"github.com/nan-www/convToMap/struct_to_map"
)

const tag = "//go:generate convToMap"

func main() {
	fmt.Println("Begin generating!")
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename.go>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	filename := os.Args[1]
	fmt.Printf("Args: %s", filename)
	generator.Gen(filename, tag, struct_to_map.GenTemplate, map_to_struct.GenTemplate)
}
