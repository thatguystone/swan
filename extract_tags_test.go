package swan

import (
	"reflect"
	"testing"
)

func TestTagsPyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/tags/",
		func(t *testing.T, name string, a *Article, r *Result) {
			et := make(map[string]interface{})
			gt := make(map[string]interface{})

			for _, t := range r.Expected.Tags {
				et[t] = nil
			}

			for _, t := range a.Meta.Tags {
				gt[t] = nil
			}

			if !reflect.DeepEqual(et, gt) {
				t.Fatalf(
					"%s: Tags do not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name,
					gt,
					et)
			}
		})
}
