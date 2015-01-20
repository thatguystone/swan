package swan

import (
	"errors"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
)

var (
	// ErrNoLanguage is returned when a document's language could not be
	// determined and Config.ErrorOnNoLanguage is set
	ErrNoLanguage = errors.New("could not determine language")

	htmlTagMatcher         = cascadia.MustCompile("html")
	metaOpenGraphMatcher   = cascadia.MustCompile("[property^=og\\:]")
	metaMatcher            = cascadia.MustCompile("meta")
	metaMatcherCanonical   = cascadia.MustCompile("[name=canonical]")
	metaMatcherDescription = cascadia.MustCompile("[name=description]")
	metaMatcherDomain      = cascadia.MustCompile("[name=domain]")
	metaMatcherFavicon     = cascadia.MustCompile("link[rel~=icon]")
	metaMatcherKeywords    = cascadia.MustCompile("[name=keywords]")
	metaMatcherLangs       = []cascadia.Selector{
		cascadia.MustCompile("[http-equiv=Content-Language]"),
		cascadia.MustCompile("[name=lang]"),
	}
)

func extractMetaLanguage(a *Article, metas *goquery.Selection) bool {
	lang, _ := a.Doc.FindMatcher(htmlTagMatcher).Attr("lang")

	if lang == "" {
		for _, s := range metaMatcherLangs {
			lang, _ = metas.FilterMatcher(s).Attr("content")
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

func extractMetas(a *Article) error {
	metas := a.Doc.FindMatcher(metaMatcher)

	if !extractMetaLanguage(a, metas) && a.cfg.Error.OnNoLanguage {
		return ErrNoLanguage
	}

	t, _ := metas.FilterMatcher(metaMatcherCanonical).Attr("content")
	a.Meta.Canonical = strings.TrimSpace(t)

	t, _ = metas.FilterMatcher(metaMatcherDescription).Attr("content")
	a.Meta.Description = strings.TrimSpace(t)

	t, _ = metas.FilterMatcher(metaMatcherDomain).Attr("content")
	a.Meta.Domain = strings.TrimSpace(t)

	t, _ = a.Doc.FindMatcher(metaMatcherFavicon).Attr("href")
	a.Meta.Favicon = strings.TrimSpace(t)

	t, _ = metas.FilterMatcher(metaMatcherKeywords).Attr("content")
	a.Meta.Keywords = strings.TrimSpace(t)

	a.Meta.OpenGraph = make(map[string]string)
	metas.FilterMatcher(metaOpenGraphMatcher).Each(
		func(i int, s *goquery.Selection) {
			if content, exists := s.Attr("content"); exists {
				prop, _ := s.Attr("property")
				a.Meta.OpenGraph[prop[3:]] = content
			}
		})

	return nil
}
