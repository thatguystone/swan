package swan

import (
	"reflect"
	"testing"
)

func TestAuthorPyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/authors/",
		func(t *testing.T, name string, a *Article, r *Result) {
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
		})
}
