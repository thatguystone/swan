package swan

import "testing"

func TestImagesPyExtractor(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/images/",
		func(t *testing.T, name string, a *Article, r *Result) {
			if r.Expected.TopImage.Src != "" && a.Img == nil {
				t.Fatalf("No image found for %s", name)
			}
		})
}
