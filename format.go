package swan

import (
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
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

	// From here on, don't work on the main document as this mangles the
	// document into text, which is not what we want for an HTML document.
	s := a.TopNode.Clone()

	// Quick-and-dirty node-to-text replacement
	s.FindMatcher(replaceWithTextTags).Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml(s.Text())
	})

	s.FindMatcher(brTags).ReplaceWithHtml("\n")

	a.CleanedText = f.getText(s)

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
