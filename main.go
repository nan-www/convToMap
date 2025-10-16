package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nan-www/convToMap/struct_to_map"
)

func main() {
	// 期望 go generate 传递要处理的文件名
	if len(os.Args) < 2 {
		// 常见错误：没有传递文件名
		fmt.Fprintf(os.Stderr, "Usage: %s <filename.go>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	filename := os.Args[1]
	struct_to_map.GenStruct2MapFile(filename)
}
