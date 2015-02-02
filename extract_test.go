package swan

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Result struct {
	URL      string
	Expected Expected
}

type Expected struct {
	Authors         []string
	CleanedText     string `json:"cleaned_text"`
	Links           int
	MetaDescription string `json:"meta_description"`
	MetaFavicon     string `json:"meta_favicon"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaLang        string `json:"meta_lang"`
	OpenGraph       map[string]string
	PublishDate     string `json:"publish_date"`
	Tags            []string
	Title           string
}

func runPyTests(
	t *testing.T,
	dir string,
	fn func(t *testing.T, name string, a *Article, r *Result)) {

	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				t.Fatalf("error while walking directory %s: %s", dir, err)
			}

			if info.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}

			r := Result{}
			name := strings.Replace(path, ".json", "", -1)

			jsonf, err := os.Open(path)
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			if err = json.NewDecoder(jsonf).Decode(&r); err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			h := strings.Replace(path, ".json", ".html", -1)
			htmlf, err := os.Open(h)
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			html, err := ioutil.ReadAll(htmlf)
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			a, err := FromHTML(string(html))
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			fn(t, name, &a, &r)
			return nil
		})
}
