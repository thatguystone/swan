package swan

import (
	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/atom"
)

type extractContent struct {
}

var (
	pTags = cascadia.MustCompile("p")
)

func (e extractContent) run(a *Article) error {
	if a.TopNode == nil {
		return nil
	}

	known := useKnownArticles{}
	if !known.isTopKnown(a) {
		e.addSiblings(a)
	}

	a.TopNode.Children().FilterFunction(func(i int, s *goquery.Selection) bool {
		if !nodeIs(s.Nodes[0], atom.P) {
			cc := a.getCCache(s.Nodes[0])
			return cc.highLinkDensity ||
				e.noParasWithoutTable(s) ||
				e.isNodeScoreThreshMet(a, s)
		}

		return false
	}).Remove()

	return nil
}

func (e extractContent) addSiblings(a *Article) {
	newTop := goquery.Selection{}
	baseScore := e.getSiblingBaseScore(a)

	a.TopNode.PrevAll().Each(func(i int, s *goquery.Selection) {
		newTop.AppendSelection(e.getSiblingContent(a, s, baseScore))
	})

	if len(newTop.Nodes) > 0 {
		newTop.AppendSelection(a.TopNode)
		a.TopNode = &newTop
	}
}

func (e extractContent) getSiblingContent(
	a *Article,
	s *goquery.Selection,
	baseScore uint) *goquery.Selection {

	if nodeIs(s.Nodes[0], atom.P) && len(s.Text()) > 0 {
		return s
	}

	ret := goquery.Selection{}
	ps := s.FindMatcher(pTags)

	for _, n := range ps.Nodes {
		cc := a.getCCache(n)
		if len(cc.text) > 0 {
			if cc.stopwords > baseScore && !cc.highLinkDensity {
				ret.AppendNodes(createNode(atom.P, "p", cc.text))
			}
		}
	}

	return &ret
}

func (e extractContent) getSiblingBaseScore(a *Article) uint {
	base := uint(100000)
	pCount := uint(0)
	pScore := uint(0)

	for _, n := range a.TopNode.FindMatcher(pTags).Nodes {
		cc := a.getCCache(n)

		if cc.stopwords > 2 && !cc.highLinkDensity {
			pCount++
			pScore += cc.stopwords
		}
	}

	if pCount > 0 {
		base = pScore / pCount
	}

	base = uint(float32(base) * float32(0.3))

	return base
}

func (e extractContent) noParasWithoutTable(s *goquery.Selection) bool {
	s.FindMatcher(pTags).Each(func(i int, s *goquery.Selection) {
		if len(s.Text()) < 25 {
			s.Remove()
		}
	})

	return s.FindMatcher(pTags).Length() == 0 && !nodeIs(s.Nodes[0], atom.Td)
}

func (e extractContent) isNodeScoreThreshMet(a *Article, s *goquery.Selection) bool {
	topNodeScore := a.scores[a.TopNode.Nodes[0]]
	currNodeScore := a.scores[s.Nodes[0]]
	threshScore := int(float32(topNodeScore) * 0.08)

	return (currNodeScore < threshScore) && !nodeIs(s.Nodes[0], atom.Td)
}
