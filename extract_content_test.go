package swan

import "testing"

func TestPyContentExtractors(t *testing.T) {
	t.Parallel()

	runPyTests(t,
		"test_data/python-goose/content/",
		func(t *testing.T, name string, a *Article, r *Result) {
			e := r.Expected

			if e.MetaDescription != "" && e.MetaDescription != a.Meta.Description {
				t.Fatalf(
					"%s: MetaDescription does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Description, e.MetaDescription)
			}

			if e.MetaKeywords != "" && e.MetaKeywords != a.Meta.Keywords {
				t.Fatalf(
					"%s: MetaKeywords does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Keywords, e.MetaKeywords)
			}

			if e.Title != "" && e.Title != a.Meta.Title {
				t.Fatalf(
					"%s: Title does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Title, e.Title)
			}

			if e.MetaLang != "" && e.MetaLang != a.Meta.Lang {
				t.Fatalf(
					"%s: Lang does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, a.Meta.Lang, e.MetaLang)
			}

			cleaned := a.CleanedText
			if len(r.Expected.CleanedText) < len(cleaned) {
				cleaned = cleaned[:len(r.Expected.CleanedText)]
			}

			if cleaned != r.Expected.CleanedText {
				t.Fatalf(
					"%s: CleanedText does not match:\n"+
						"	Got: %s\n"+
						"	Expected: %s",
					name, cleaned, r.Expected.CleanedText)
			}
		})
}
