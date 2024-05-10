package main

import (
	"bufio"
	_ "embed"
	"encoding/csv"
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
)

func parse(s string) string {
	result := ""

	isPre := false
	isList := false

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()

		if isPre {
			if line == "```" {
				isPre = false
				result += "</pre>\n"
			} else {
				result += line + "\n"
			}
			continue
		}

		if isList && !list.MatchString(line) {
			result += "</ul>\n"
			isList = false
		}

		switch {
		case h3.MatchString(line):
			result += `<h3>` + line[4:] + `</h3>`
		case h2.MatchString(line):
			result += `<h2>` + line[3:] + `</h2>`
		case h1.MatchString(line):
			result += `<h1>` + line[2:] + `</h1>`
		case quote.MatchString(line):
			result += `<blockquote>` + line[2:] + `</blockquote>`
		case anchor.MatchString(line):
			link := anchor.FindStringSubmatch(line)[2]
			ext := link[len(link)-4:]
			if ext == ".csv" {
				file, _ := os.ReadFile("page/" + link)
				result += table(string(file), ',')
			} else if ext == ".tsv" {
				file, _ := os.ReadFile("page/" + link)
				result += table(string(file), '\t')
			} else {
				result += anchor.ReplaceAllString(line, `<a href="$2">$1</a><br>`)
			}
		case line == "```":
			result += `<pre>`
			isPre = true
		case list.MatchString(line):
			if !isList {
				result += "<ul>\n"
			}
			result += `  <li>` + line[2:] + `</li>`
			isList = true
		default:
			result += line
		}

		result += "\n"
	}

	if isPre {
		result += "</pre>\n"
	}

	if isList {
		result += "</ul>\n"
	}

	return result
}

func table(s string, delim rune) string {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = delim
	rows, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	t := "<table>\n"
	for i, row := range rows {
		t += "  <tr>\n"
		for _, cell := range row {
			if i == 0 {
				t += "    <th>" + cell + "</th>\n"
			} else {
				t += "    <td>" + cell + "</td>\n"
			}
		}
		t += "  </tr>\n"
	}
	t += "</table>"
	return t
}

func toc(w http.ResponseWriter, path string) {
	files, err := os.ReadDir("page/" + path)
	if err != nil {
		return
	}

	for _, file := range files {
		link := path + "/" + file.Name()
		fmt.Fprintf(w, `<br><a href="%s">%s</a>`, link, file.Name())
	}
	fmt.Fprint(w, "<hr>")
}

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<style>%s</style>", style)
	fmt.Fprintf(w, "Page: %s\n\n", r.URL.Path)

	path := "index.md"
	if r.URL.Path != "/" {
		path = r.URL.Path
	}
	info, _ := os.Stat("page/" + path)
	if info.IsDir() {
		toc(w, r.URL.Path)
	}
	file, _ := os.ReadFile("page/" + path)
	fmt.Fprint(w, parse(string(file)))
}

func main() {
	fmt.Println("http://localhost:8080")
	http.HandleFunc("/", page)
	http.HandleFunc("/favicon.ico",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write(favicon)
		})
	http.ListenAndServe(":8080", nil)
}
