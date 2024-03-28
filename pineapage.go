package main

import (
	"bufio"
	"encoding/csv"
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
			if link[len(link)-4:] == ".csv" {
				file, _ := os.ReadFile("page/" + link)
				result += table(string(file))
			} else {
				result += anchor.ReplaceAllString(line, `<a href="$2">$1</a>`)
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

func table(s string) string {
	r := csv.NewReader(strings.NewReader(s))
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
