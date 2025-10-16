package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nan-www/convToMap/struct_to_map"
)

func main() {
	fmt.Println("Begin generating!")
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename.go>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	filename := os.Args[1]
	fmt.Printf("Args: %s", filename)
	struct_to_map.GenStruct2MapFile(filename)
}
