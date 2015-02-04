package swan

import (
	"bytes"
	"net/url"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Article is a fully extracted and cleaned document.
type Article struct {
	// Final URL after all redirects
	URL string

	// Newline-separated and cleaned content
	CleanedText string

	// Node from which CleanedText was created
	TopNode *goquery.Selection

	// A header image to use for the article
	Img *Image

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

	// Full document backing this article
	Doc *goquery.Document

	// For use resolving URLs in the document
	baseURL *url.URL

	// Caches information about nodes so that it doesn't have to be updated
	cCache map[*html.Node]*contentCache

	// Scores that have been calculated
	scores map[*html.Node]int
}

// Image contains information about the header image associated with an article
type Image struct {
	Src        string
	Width      int
	Height     int
	Bytes      int64
	Confidence uint
}

type contentCache struct {
	text            string
	wordCount       uint
	stopwords       uint
	highLinkDensity bool
	s               *goquery.Selection
}

type runner interface {
	run(a *Article) error
}

type useKnownArticles struct{}

var (
	runners = []runner{
		extractMetas{},

		extractAuthors{},
		extractPublishDate{},
		extractTags{},
		extractTitle{},

		cleanup{},
		useKnownArticles{},
		metaDetectLanguage{},

		extractTopNode{},
		extractLinks{},
		extractImages{},
		extractVideos{},

		// Does more document mangling and TopNode resetting
		extractContent{},
	}

	// Don't match all-at-once: there's precedence here
	knownArticles = []goquery.Matcher{
		cascadia.MustCompile("[itemprop=articleBody]"),
		cascadia.MustCompile(".post-content"),
		cascadia.MustCompile("article"),
	}
)

func (u useKnownArticles) run(a *Article) error {
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

// Checks to see if TopNode is a known article tag that was picked before
// scoring
func (u useKnownArticles) isTopKnown(a *Article) bool {
	for _, m := range knownArticles {
		if a.TopNode.IsMatcher(m) {
			return true
		}
	}

	return false
}

func (a *Article) extract() error {
	a.cCache = make(map[*html.Node]*contentCache)
	a.scores = make(map[*html.Node]int)

	for _, r := range runners {
		err := r.run(a)
		if err != nil {
			return err
		}
	}

	a.cCache = nil
	a.scores = nil

	return nil
}

func (a *Article) getCCache(n *html.Node) *contentCache {
	cc, ok := a.cCache[n]
	if !ok {
		s := goquery.NewDocumentFromNode(n).Selection
		cc = &contentCache{
			text: strings.TrimSpace(s.Text()),
			s:    s,
		}

		ws := splitText(cc.text)
		cc.wordCount = uint(len(ws))
		cc.stopwords = stopwordCountWs(a.Meta.Lang, ws)
		cc.highLinkDensity = highLinkDensity(cc)
		a.cCache[n] = cc
	}

	return cc
}

func highLinkDensity(cc *contentCache) bool {
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
