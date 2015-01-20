package swan

import (
	"testing"
)

func TestPublishDatePyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/publishdate/",
		func(t *testing.T, name string, a *Article, r *Result) {
			if a.Meta.PublishDate != r.Expected.PublishDate {
				t.Fatalf(
					"%s: PublishDate does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.PublishDate, r.Expected.PublishDate)
			}
		})
}
