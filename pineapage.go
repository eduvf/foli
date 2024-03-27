package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	h3     = regexp.MustCompile(`^### (.*)`)
	h2     = regexp.MustCompile(`^## (.*)`)
	h1     = regexp.MustCompile(`^# (.*)`)
	quote  = regexp.MustCompile(`^> (.*)`)
	anchor = regexp.MustCompile(`^\[(.*)\]\((.*)\)`)
)

func parse(s string) string {
	result := ""
	isPre := false
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
			result += anchor.ReplaceAllString(line, `<a href="$2">$1</a>`)
		case line == "```":
			result += `<pre>`
			isPre = true
		default:
			result += line
		}

		result += "\n"
	}

	return result
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
