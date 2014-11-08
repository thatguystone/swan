package swan

import (
	"fmt"
	"strings"
	"errors"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	doc     *goquery.Document
	cfg Config
	Cleaned string
	Meta    struct {
		Description string
		Favicon     string
		Keywords    string
		Language    string
		Title       string
	}
}

var (
	ErrorNoLanguage = errors.New("Could not determine language")

	metaDescription = cascadia.MustCompile("meta[name=description]")
	metaKeywords    = cascadia.MustCompile("meta[name=keywords]")
	metaFavicon     = cascadia.MustCompile("link[rel~=icon]")
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

	fmt.Printf("desc=%s\n", a.Meta.Description)
	fmt.Printf("favicon=%s\n", a.Meta.Favicon)
	fmt.Printf("keywords=%s\n", a.Meta.Keywords)
	fmt.Printf("language=%s\n", a.Meta.Language)
	fmt.Printf("title=%s\n", a.Meta.Title)

	return a, nil
}

func (a *Article) extractLanguage() bool {
	// a.setMeta("language", &a.Meta.Language)
	return false
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

	a.Meta.Title = title
}
