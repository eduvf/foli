package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Page struct {
	Title string
	Body  string
}

func main() {
	dir := flag.String("d", ".", "directory to work on")
	flag.Parse()

	err := filepath.WalkDir(*dir, scan)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", view)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func view(writer http.ResponseWriter, request *http.Request) {
	p := &Page{
		Title: "Test",
		Body:  "Hi!",
	}
	fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
