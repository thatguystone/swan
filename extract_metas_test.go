package swan

import (
	"reflect"
	"testing"
)

func TestOpenGraphPyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/opengraph/",
		func(t *testing.T, name string, a *Article, r *Result) {
			if !reflect.DeepEqual(a.Meta.OpenGraph, r.Expected.OpenGraph) {
				t.Fatalf(
					"%s: Authors do not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name,
					a.Meta.OpenGraph,
					r.Expected.OpenGraph)
			}
		})
}
