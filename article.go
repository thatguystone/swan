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

	// Document backing this article
	Doc *goquery.Document

	cfg Config
}

var (
	// ErrNoLanguage is returned when a document's language could not be
	// determined and Config.ErrorOnNoLanguage is set
	ErrNoLanguage = errors.New("could not determine language")

	htmlTag         = cascadia.MustCompile("html")
	metaDescription = cascadia.MustCompile("meta[name=description]")
	metaKeywords    = cascadia.MustCompile("meta[name=keywords]")
	metaFavicon     = cascadia.MustCompile("link[rel~=icon]")
	metaLangs       = []cascadia.Selector{
		cascadia.MustCompile("meta[http-equiv=Content-Language]"),
		cascadia.MustCompile("meta[name=lang]"),
	}
)

func (a *Article) extract() error {
	if !a.determineLanguage() && a.cfg.Error.OnNoLanguage {
		return ErrNoLanguage
	}

	a.findMetas()
	a.findTitle()
	return a.extractContent()
}

func (a *Article) determineLanguage() bool {
	lang, _ := a.Doc.FindMatcher(htmlTag).Attr("lang")

	if lang == "" {
		for _, s := range metaLangs {
			lang, _ = a.Doc.FindMatcher(s).Attr("content")
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

func (a *Article) findMetas() {
	t, _ := a.Doc.FindMatcher(metaDescription).Attr("content")
	a.Meta.Description = strings.TrimSpace(t)

	t, _ = a.Doc.FindMatcher(metaFavicon).Attr("href")
	a.Meta.Favicon = strings.TrimSpace(t)

	t, _ = a.Doc.FindMatcher(metaKeywords).Attr("content")
	a.Meta.Keywords = strings.TrimSpace(t)
}

func (a *Article) findTitle() {
	title := a.Doc.Find("title").Text()
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
