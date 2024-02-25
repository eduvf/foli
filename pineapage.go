package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("d", ".", "directory to work on")
	flag.Parse()

	err := filepath.WalkDir(*dir, scan)
	if err != nil {
		fmt.Println(err)
	}
}

func scan(path string, file fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}
	if strings.HasPrefix(path, ".") {
		return nil
	}
	fmt.Printf("%v %s\n", file.IsDir(), path)
	return nil
}
