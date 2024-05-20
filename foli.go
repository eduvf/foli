package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//go:embed style.css
var style string

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

func warn(err error) bool {
	if err != nil {
		fmt.Println(YELLOW + err.Error())
		return true
	}
	return false
}

func parse(w http.ResponseWriter, md []byte) {
	scan := bufio.NewScanner(bytes.NewReader(md))
	write := func(s string) { fmt.Fprint(w, s) }

	for scan.Scan() {
		line := scan.Text()

		switch {
		case h3.MatchString(line):
			write(`<h3>` + line[4:] + `</h3>`)
		case h2.MatchString(line):
			write(`<h2>` + line[3:] + `</h2>`)
		case h1.MatchString(line):
			write(`<h1>` + line[2:] + `</h1>`)
		case quote.MatchString(line):
			write(`<blockquote>` + line[2:] + `</blockquote>`)
		case anchor.MatchString(line):
			write(anchor.ReplaceAllString(line, `<p><a href="$2">$1</a></p>`))
		case img.MatchString(line):
			write(img.ReplaceAllString(line, `<figure><img alt="$1" src="$2"><figcaption>$1</figcaption></figure>`))
		case line == "```":
			write(`<pre>`)
			for scan.Scan() && scan.Text() != "```" {
				write(scan.Text() + "\n")
			}
			write(`</pre>`)
		case list.MatchString(line):
			write("<ul>")
			write("<li>" + line[2:] + "</li>")
			for scan.Scan() && list.MatchString(scan.Text()) {
				write("<li>" + scan.Text()[2:] + "</li>")
			}
			write("</ul>")
		default:
			write(line + "\n")
		}
	}
}

func nav(w http.ResponseWriter) {
	fmt.Fprint(w, "<nav>")
	home, err := os.ReadDir("page")
	if warn(err) {
		return
	}
	for _, entry := range home {
		name := strings.TrimSuffix(entry.Name(), ".md")
		fmt.Fprintf(w, ` <a href="/%s">%s</a>`, entry.Name(), name)
	}
	fmt.Fprint(w, "</nav>")
}

func ls(w http.ResponseWriter, dir []fs.DirEntry, path string) {
	for _, entry := range dir {
		name := entry.Name()
		link := path + "/" + name
		fmt.Fprintf(w, `<p><a href="%s">%s</a></p>`, link, name)
	}
}

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<style>%s</style>", style)
	nav(w)

	path := r.URL.Path
	if path == "/" {
		path = "home.md"
	}

	info, err := os.Stat("page/" + path)
	if warn(err) {
		return
	}
	if info.IsDir() {
		dir, err := os.ReadDir("page/" + path)
		if warn(err) {
			return
		}
		ls(w, dir, path)
	} else {
		file, err := os.ReadFile("page/" + path)
		if warn(err) {
			return
		}
		parse(w, file)
	}
}

func main() {
	fmt.Println(GREEN + "http://localhost:8080")
	http.HandleFunc("/", page)
	http.ListenAndServe(":8080", nil)
}
