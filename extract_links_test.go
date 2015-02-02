package swan

import (
	"testing"
)

func TestLinksPyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/links/",
		func(t *testing.T, name string, a *Article, r *Result) {
			if len(a.Meta.Links) != r.Expected.Links {
				t.Fatalf(
					"%s: Incorrect link count:\n"+
						"	Got: %d\n"+
						"	Expected: %d",
					name, len(a.Meta.Links), r.Expected.Links)
			}
		})
}
