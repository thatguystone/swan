package swan

// Extractor gets useful content from a document
type Extractor interface {
	// Score returns a score, 0-100, for how well it matches a document. If 0,
	// it's considered not to match. If 100, it matches immediately and no more
	// scoring is done.
	Score(a *Article) int

	// Extract performs article extraction, setting both CleanedText and
	// CleanedHTML on the article, optionally overwriting any pre-determined
	// meta values.
	Extract(a *Article) error
}

var extractors = []Extractor{}

// AddExtractor adds an extractor to the extractors used for document retrieval.
func AddExtractor(e Extractor) {
	extractors = append(extractors, e)
}

func (a *Article) extractContent() error {
	var ext Extractor
	score := -1

	for _, e := range extractors {
		s := e.Score(a)
		if s > score {
			score = s
			ext = e
		}

		if score == 100 {
			break
		}
	}

	return ext.Extract(a)
}
