package swan

import (
	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
)

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
		Links       []string
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
	processors = []func(a *Article) error{
		extractMetas,

		extractAuthors,
		extractPublishDate,
		extractTags,
		extractTitle,

		useKnownArticles,
		cleanup,
		extractContent,

		extractLinks,
		extractImages,
		extractVideos,
	}

	// Don't match all-at-once: there's precedence here
	knownArticles = []goquery.Matcher{
		cascadia.MustCompile("[itemprop=articleBody]"),
		cascadia.MustCompile(".post-content"),
		cascadia.MustCompile("article"),
	}
)

func useKnownArticles(a *Article) error {
	for _, m := range knownArticles {
		s := a.Doc.FindMatcher(m)
		if s.Size() > 0 {
			// Remove from document so that memory can be freed
			f := s.First().Remove()
			a.Doc = goquery.NewDocumentFromNode(f.Nodes[0])
			break
		}
	}

	return nil
}

func (a *Article) extract() error {
	for _, f := range processors {
		err := f(a)
		if err != nil {
			return err
		}
	}

	return nil
}
