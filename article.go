package swan

import "github.com/PuerkitoBio/goquery"

// Article is a fully extracted and cleaned document.
type Article struct {
	// Newline-separated and cleaned content
	CleanedText string

	// HTML-formatted content with inline images, videos, and whatever else was
	// found relevant to the original article
	CleanedHTML string

	// All metadata associated with the original document
	Meta struct {
		Authors     []string
		Canonical   string
		Description string
		Domain      string
		Favicon     string
		Keywords    string
		Lang        string
		OpenGraph   map[string]string
		PublishDate string
		Tags        []string
		Title       string
	}

	// Document backing this article
	Doc *goquery.Document

	cfg Config
}

var (
	extractors = []func(a *Article) error{
		extractMetas,

		extractAuthors,
		extractLinks,
		extractPublishDate,
		extractTags,
		extractTitle,

		extractContent,
		extractImages,
		extractVideos,
	}
)

func (a *Article) extract() error {
	for _, e := range extractors {
		err := e(a)
		if err != nil {
			return err
		}
	}

	return nil
}
