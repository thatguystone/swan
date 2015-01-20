package swan

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	PyOpenGraphDir = "test_data/python-goose/opengraph/"
)

func TestOpenGraphPyExtractor(t *testing.T) {
	t.Parallel()

	filepath.Walk(PyOpenGraphDir,
		func(path string, info os.FileInfo, err error) error {
			name, r, html, ok := readPyTest(t, path, info)
			if !ok {
				return nil
			}

			a, err := FromHTML(html, Config{})
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(a.Meta.OpenGraph, r.Expected.OpenGraph) {
				t.Fatalf(
					"%s: Authors do not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name,
					a.Meta.OpenGraph,
					r.Expected.OpenGraph)
			}

			return nil
		})
}
