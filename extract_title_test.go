package swan

import (
	"testing"
)

func TestTitlePyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/title/",
		func(t *testing.T, name string, a *Article, r *Result) {
			if a.Meta.Title != r.Expected.Title {
				t.Fatalf(
					"%s: Title does not match:\n"+
						"	Got:      %s\n"+
						"	Expected: %s",
					name, a.Meta.Title, r.Expected.Title)
			}
		})
}
