package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var (
	h3     = regexp.MustCompile(`(?m)^### (.*)`)
	h2     = regexp.MustCompile(`(?m)^## (.*)`)
	h1     = regexp.MustCompile(`(?m)^# (.*)`)
	pre    = regexp.MustCompile("(?s)```(.*)```")
	quote  = regexp.MustCompile(`(?m)^> (.*)`)
	anchor = regexp.MustCompile(`\[(.*)\]\((.*)\)`)
)

func parse(s string) string {
	s = h3.ReplaceAllString(s, `<h3>$1</h3>`)
	s = h2.ReplaceAllString(s, `<h2>$1</h2>`)
	s = h1.ReplaceAllString(s, `<h1>$1</h1>`)
	s = pre.ReplaceAllString(s, `<pre>$1</pre>`)
	s = quote.ReplaceAllString(s, `<blockquote>$1</blockquote>`)
	s = anchor.ReplaceAllString(s, `<a href="$2">$1</a>`)
	return s
}

func page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Page: %s\n\n", r.URL.Path)
	file, _ := os.ReadFile("page/index.md")
	fmt.Fprint(w, parse(string(file)))
}

func main() {
	fmt.Println("http://localhost:8080")
	http.HandleFunc("/", page)
	http.ListenAndServe(":8080", nil)
}
