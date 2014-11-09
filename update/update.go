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
	out := bytes.NewBufferString("package swan\n\nvar(stopwords = map[string][]string{\n")

	files, err := filepath.Glob("python-goose/goose/resources/text/stopwords*")
	if err != nil {
		log.Fatalf("Could not glob for stopwords: %s", err)
	}

	for _, f := range files {
		lang := f
		lang = lang[:len(lang)-4]
		lang = lang[len(lang)-2:]

		c, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("Could not read file: %s", err)
		}

		// Some stopword files from python-goose start with an invalid UTF8 char
		ws := string(c)
		ws = strings.Trim(ws, "\ufeff")

		out.WriteString(fmt.Sprintf("`%s`: []string{\n", lang))
		for _, w := range strings.Split(ws, "\n") {
			if !strings.HasPrefix(w, "#") && len(w) > 0 {
				out.WriteString(fmt.Sprintf("`%s`,", strings.TrimSpace(w)))
			}
		}
		out.WriteString("},\n")
	}

	out.WriteString("})\n")

	ioutil.WriteFile("../stopwords.go", out.Bytes(), 0644)
}
