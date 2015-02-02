// Internal tool: update the stopwords list from python-goose. Use `make` to run it.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	out := bytes.NewBufferString("package swan\n\nvar(stopwords = map[string]map[string]bool{\n")

	files, err := filepath.Glob("python-goose/goose/resources/text/stopwords*")
	if err != nil {
		log.Fatalf("Could not glob for stopwords: %s", err)
	}

	for _, f := range files {
		lang := f
		lang = lang[:len(lang)-4]
		lang = lang[len(lang)-2:]
		seen := make(map[string]bool)

		c, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("Could not read file: %s", err)
		}

		// Some stopword files from python-goose start with an invalid UTF8 char
		ws := string(c)
		ws = strings.Trim(ws, "\ufeff")

		out.WriteString(fmt.Sprintf("`%s`: map[string]bool{\n", lang))
		for _, w := range strings.Split(ws, "\n") {
			w = strings.TrimSpace(w)
			if !strings.HasPrefix(w, "#") && len(w) > 0 {
				if _, ok := seen[w]; !ok {
					fmt.Fprintf(out, "`%s`: true,\n", w)
					seen[w] = true
				}
			}
		}
		out.WriteString("},\n")
	}

	out.WriteString("})\n")

	ioutil.WriteFile("../stopwords_list.go", out.Bytes(), 0644)
}
