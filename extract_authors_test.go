package swan

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	PyAuthorsDir = "test_data/python-goose/authors/"
)

func TestAuthorPyExtractor(t *testing.T) {
	t.Parallel()

	filepath.Walk(PyAuthorsDir,
		func(path string, info os.FileInfo, err error) error {
			name, r, html, ok := readPyTest(t, path, info)
			if !ok {
				return nil
			}

			a, err := FromHTML(html, Config{})
			if err != nil {
				t.Fatal(err)
			}

			ea := make(map[string]interface{})
			ga := make(map[string]interface{})

			for _, a := range r.Expected.Authors {
				ea[a] = nil
			}

			for _, a := range a.Meta.Authors {
				ga[a] = nil
			}

			if !reflect.DeepEqual(ea, ga) {
				t.Fatalf(
					"%s: Authors do not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name,
					ga,
					ea)
			}

			return nil
		})
}
