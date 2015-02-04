package swan

import (
	"bytes"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify"
)

type formatCleanedText struct{}

var (
	brTags              = cascadia.MustCompile("br")
	replaceWithTextTags = cascadia.MustCompile("a, b, strong, i, sup")
)

func (f formatCleanedText) run(a *Article) error {
	if a.TopNode == nil {
		return nil
	}

	f.dropNegativeScored(a)

	var b bytes.Buffer
	html, _ := a.TopNode.Html()
	minify.NewMinifier().HTML(&b, strings.NewReader(html))
	doc, _ := goquery.NewDocumentFromReader(&b)

	a.TopNode = doc.Selection

	// Quick-and-dirty node-to-text replacement
	a.TopNode.FindMatcher(replaceWithTextTags).Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml(s.Text())
	})

	a.TopNode.FindMatcher(brTags).ReplaceWithHtml("\n")

	a.CleanedText = f.getText(a.TopNode)

	return nil
}

func (f formatCleanedText) dropNegativeScored(a *Article) {
	for n, score := range a.scores {
		if score <= 0 {
			if n.Parent != nil {
				n.Parent.RemoveChild(n)
			}
		}
	}
}

func (f formatCleanedText) getText(s *goquery.Selection) string {
	s.FindMatcher(pTags).Each(func(i int, s *goquery.Selection) {
		s.AfterHtml("\n\n")
	})

	return strings.TrimSpace(s.Text())
}
