package swan

import (
	"errors"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
)

// Article is a fully extracted and cleaned document.
type Article struct {
	doc         *goquery.Document
	cfg         Config
	Cleaned     string
	CleanedHTML string
	Meta        struct {
		Description string
		Favicon     string
		Keywords    string
		Lang        string
		Title       string
	}
}

var (
	// ErrorNoLanguage is returned when a document's language could not be
	// determined and Config.ErrorOnNoLanguage is set
	ErrorNoLanguage = errors.New("Could not determine language")

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

	if !a.extractLanguage() && cfg.ErrorOnNoLanguage {
		return Article{}, ErrorNoLanguage
	}

	a.extractMetas()
	a.extractTitle()

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
