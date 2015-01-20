package swan

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	PyContentDir = "test_data/python-goose/content/"
)

type Result struct {
	URL      string
	Expected Expected
}

type Expected struct {
	Authors         []string
	CleanedText     string `json:"cleaned_text"`
	MetaDescription string `json:"meta_description"`
	MetaFavicon     string `json:"meta_favicon"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaLang        string `json:"meta_lang"`
	OpenGraph       map[string]string
	Tags            []string
	Title           string
}

func readPyTest(
	t *testing.T,
	path string,
	info os.FileInfo) (string, Result, string, bool) {
	if info.IsDir() || !strings.HasSuffix(path, ".json") {
		return "", Result{}, "", false
	}

	r := Result{}
	name := strings.Replace(path, ".json", "", -1)

	jsonf, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	if err = json.NewDecoder(jsonf).Decode(&r); err != nil {
		t.Fatal(err)
	}

	h := strings.Replace(path, ".json", ".html", -1)
	htmlf, err := os.Open(h)
	if err != nil {
		t.Fatal(err)
	}

	html, err := ioutil.ReadAll(htmlf)
	if err != nil {
		t.Fatal(err)
	}

	return name, r, string(html), true
}

func TestPyExtractors(t *testing.T) {
	t.Parallel()

	filepath.Walk(PyContentDir,
		func(path string, info os.FileInfo, err error) error {
			name, r, html, ok := readPyTest(t, path, info)
			if !ok {
				return nil
			}

			a, err := FromHTML(html, Config{})
			if err != nil {
				t.Fatal(err)
			}

			e := r.Expected

			if e.MetaDescription != "" && e.MetaDescription != a.Meta.Description {
				t.Fatalf(
					"%s: MetaDescription does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Description, e.MetaDescription)
			}

			if e.MetaKeywords != "" && e.MetaKeywords != a.Meta.Keywords {
				t.Fatalf(
					"%s: MetaKeywords does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Keywords, e.MetaKeywords)
			}

			if e.Title != "" && e.Title != a.Meta.Title {
				t.Fatalf(
					"%s: Title does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Title, e.Title)
			}

			if e.MetaLang != "" && e.MetaLang != a.Meta.Lang {
				t.Fatalf(
					"%s: Lang does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Lang, e.MetaLang)
			}

			return nil
		})
}
