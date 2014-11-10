package swan

import (
	"errors"
	"strings"

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
		Description string
		Favicon     string
		Keywords    string
		Lang        string
		Title       string
	}

	doc *goquery.Document
	cfg Config
}

var (
	// ErrNoLanguage is returned when a document's language could not be
	// determined and Config.ErrorOnNoLanguage is set
	ErrNoLanguage = errors.New("could not determine language")

	html            = cascadia.MustCompile("html")
	metaDescription = cascadia.MustCompile("meta[name=description]")
	metaKeywords    = cascadia.MustCompile("meta[name=keywords]")
	metaFavicon     = cascadia.MustCompile("link[rel~=icon]")
	metaLangs       = []cascadia.Selector{
		cascadia.MustCompile("meta[http-equiv=Content-Language]"),
		cascadia.MustCompile("meta[name=lang]"),
	}
)

func prepareArticle(doc *goquery.Document, cfg Config) (Article, error) {
	a := Article{
		cfg: cfg,
		doc: doc,
	}

	if !a.extractLanguage() && cfg.Error.OnNoLanguage {
		return Article{}, ErrNoLanguage
	}

	a.extractMetas()
	a.extractTitle()
	a.extractContent()

	return a, nil
}

func (a *Article) extractLanguage() bool {
	lang, _ := a.doc.FindMatcher(html).Attr("lang")

	if lang == "" {
		for _, s := range metaLangs {
			lang, _ = a.doc.FindMatcher(s).Attr("content")
			if lang != "" {
				break
			}
		}
	}

	if lang != "" {
		a.Meta.Lang = lang[:2]
	}

	return a.Meta.Lang != ""
}

func (a *Article) extractMetas() {
	t, _ := a.doc.FindMatcher(metaDescription).Attr("content")
	a.Meta.Description = strings.TrimSpace(t)

	t, _ = a.doc.FindMatcher(metaFavicon).Attr("href")
	a.Meta.Favicon = strings.TrimSpace(t)

	t, _ = a.doc.FindMatcher(metaKeywords).Attr("content")
	a.Meta.Keywords = strings.TrimSpace(t)
}

func (a *Article) extractTitle() {
	title := a.doc.Find("title").Text()
	if title == "" {
		return
	}

	delim := ""
	switch {
	case strings.Contains(title, "|"):
		delim = "|"
	case strings.Contains(title, "-"):
		delim = " - "
	case strings.Contains(title, "»"):
		delim = "»"
	case strings.Contains(title, ":"):
		delim = ":"
	}

	if delim != "" {
		longest := 0
		parts := strings.Split(title, delim)

		for i, t := range parts {
			if len(t) > len(parts[longest]) {
				longest = i
			}
		}

		title = parts[longest]
	}

	a.Meta.Title = strings.TrimSpace(title)
}
