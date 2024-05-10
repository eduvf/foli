package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//go:embed style.css
var style string

//go:embed favicon.ico
var favicon []byte

var (
	h3     = regexp.MustCompile(`^### (.*)`)
	h2     = regexp.MustCompile(`^## (.*)`)
	h1     = regexp.MustCompile(`^# (.*)`)
	list   = regexp.MustCompile(`^\* (.*)`)
	quote  = regexp.MustCompile(`^> (.*)`)
	anchor = regexp.MustCompile(`^\[(.*)\]\((.*)\)`)
	img    = regexp.MustCompile(`^!\[(.*)\]\((.*)\)`)
)

const GREEN = "\x1b[32m"
const YELLOW = "\x1b[33m"

func parse(s string) string {
	res := ""
	scan := bufio.NewScanner(strings.NewReader(s))

	for scan.Scan() {
		line := scan.Text()

		switch {
		case h3.MatchString(line):
			res += `<h3>` + line[4:] + `</h3>`
		case h2.MatchString(line):
			res += `<h2>` + line[3:] + `</h2>`
		case h1.MatchString(line):
			res += `<h1>` + line[2:] + `</h1>`
		case quote.MatchString(line):
			res += `<blockquote>` + line[2:] + `</blockquote>`
		case anchor.MatchString(line):
			res += anchor.ReplaceAllString(line, `<p><a href="$2">$1</a></p>`)
		case img.MatchString(line):
			res += img.ReplaceAllString(line, `<figure><img alt="$1" src="$2"><figcaption>$1</figcaption></figure>`)
		case line == "```":
			res += `<pre>`
			for scan.Scan() && scan.Text() != "```" {
				res += scan.Text() + "\n"
			}
			res += `</pre>`
		case list.MatchString(line):
			res += "<ul>"
			res += "<li>" + line[2:] + "</li>"
			for scan.Scan() && list.MatchString(scan.Text()) {
				res += "<li>" + scan.Text()[2:] + "</li>"
			}
			res += "</ul>"
		default:
			res += line + "\n"
		}
	}

	return res
}

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<style>%s</style>", style)
	fmt.Fprintf(w, "Page: %s\n\n", r.URL.Path)
	fmt.Fprint(w, "<hr>")

	path := r.URL.Path
	if path == "/" {
		path = "index.md"
	}

	file, err := os.ReadFile("page/" + path)
	if err != nil {
		fmt.Print(YELLOW + err.Error())
	}
	fmt.Fprint(w, parse(string(file)))
}

func main() {
	fmt.Println(GREEN + "http://localhost:8080")
	http.HandleFunc("/", page)
	http.HandleFunc("/favicon.ico",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write(favicon)
		})
	http.ListenAndServe(":8080", nil)
}
