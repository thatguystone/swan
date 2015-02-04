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

			if a.Img == nil {
				return
			}

			if r.Expected.TopImage.Src != a.Img.Src {
				t.Fatalf("Found wrong image for %s:\n"+
					"	Got:      %s\n"+
					"	Expected: %s",
					name,
					a.Img.Src,
					r.Expected.TopImage.Src)
			}

			if r.Expected.TopImage.Height != a.Img.Height ||
				r.Expected.TopImage.Width != a.Img.Width {

				t.Fatalf("Dimension mismatch for %s: got %dx%d, expected %dx%d",
					name,
					a.Img.Width, a.Img.Height,
					r.Expected.TopImage.Width, r.Expected.TopImage.Height)
			}
		})
}
