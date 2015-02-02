package swan

import (
	"bytes"
	"math"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type extractTopNode struct {
	a      *Article
	cache  map[*html.Node]*contentCache
	scores map[*html.Node]uint
}

type contentCache struct {
	text      string
	wordCount uint
	stopwords uint
	s         *goquery.Selection
}

const (
	maxStepsFromNode = 3
	minStopwordCount = 5
)

var (
	linkTags     = cascadia.MustCompile("a")
	nodesToCheck = cascadia.MustCompile("p, pre, td")
)

func (e extractTopNode) run(a *Article) error {
	var ccs []*contentCache

	e.a = a
	e.cache = make(map[*html.Node]*contentCache)
	e.scores = make(map[*html.Node]uint)

	for _, n := range a.Doc.FindMatcher(nodesToCheck).Nodes {
		cc := e.hitCache(n)
		if cc.stopwords > 2 && !e.highLinkDensity(cc) {
			ccs = append(ccs, cc)
		}
	}

	startingBoost := 1.0
	bottomNegativeScore := int(float32(len(ccs)) * 0.25)

	for i, cc := range ccs {
		boostScore := 0.0
		if i > 0 && e.isBoostable(cc) {
			boostScore = (1.0 / startingBoost) * 50
			startingBoost++
		}

		if len(ccs) > 15 {
			if (len(ccs) - i) <= bottomNegativeScore {
				booster := float64(bottomNegativeScore - (len(ccs) - i))
				boostScore = -math.Pow(booster, 2.0)
				if math.Abs(boostScore) > 40 {
					boostScore = 5.0
				}
			}
		}

		upscore := cc.stopwords + uint(boostScore)

		p := cc.s.Nodes[0].Parent
		if p == nil {
			continue
		}

		score, _ := e.scores[p]
		e.scores[p] = score + upscore

		p = p.Parent
		if p != nil {
			pscore, _ := e.scores[p]
			e.scores[p] = pscore + (upscore / 2)
		}
	}

	var topNode *html.Node
	topScore := uint(0)
	for n, score := range e.scores {
		if score > topScore {
			topNode = n
			topScore = score
		}

		if topNode == nil {
			topNode = n
		}
	}

	if topNode != nil {
		a.TopNode = goquery.NewDocumentFromNode(topNode).Selection
	}

	return nil
}

func (e *extractTopNode) hitCache(n *html.Node) *contentCache {
	cc, ok := e.cache[n]
	if !ok {
		s := goquery.NewDocumentFromNode(n).Selection
		cc = &contentCache{
			text: s.Text(),
			s:    s,
		}

		ws := splitText(cc.text)
		cc.wordCount = uint(len(ws))
		cc.stopwords = stopwordCountWs(e.a.Meta.Lang, ws)
		e.cache[n] = cc
	}

	return cc
}

func (e *extractTopNode) isBoostable(cc *contentCache) bool {
	stepsAway := 0
	for sib := cc.s.Nodes[0].PrevSibling; sib != nil; sib = sib.PrevSibling {
		if sib.Type == html.ElementNode && sib.DataAtom == atom.P {
			if stepsAway > maxStepsFromNode {
				return false
			}

			scc := e.hitCache(sib)
			if scc.stopwords > minStopwordCount {
				return true
			}
		}

		stepsAway++
	}

	return false
}

func (e *extractTopNode) highLinkDensity(cc *contentCache) bool {
	var b bytes.Buffer

	links := cc.s.FindMatcher(linkTags)

	if links.Size() == 0 {
		return false
	}

	links.Each(func(i int, l *goquery.Selection) {
		b.WriteString(l.Text())
	})

	linkWords := float32(strings.Count(b.String(), " "))

	return ((linkWords / float32(cc.wordCount)) * float32(len(cc.s.Nodes))) >= 1
}
