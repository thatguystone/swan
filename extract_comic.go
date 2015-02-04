package swan

import (
	"fmt"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type extractComic struct{}

var (
	comicProcessor = &processor{
		probe: comicProbe,
		runners: []runner{
			extractComic{},
		},
	}

	comicDomains = []string{
		"xkcd.com",
	}

	comicKeywords = []string{
		"webcomic",
		"comic strip",
	}
)

func comicProbe(a *Article) uint {
	for _, d := range comicDomains {
		if a.baseURL.Host == d {
			return 100
		}
	}

	score := uint(0)
	for _, kw := range comicKeywords {
		if strings.Contains(a.Meta.Keywords, kw) {
			score += 10
		}
	}

	return score
}

func (e extractComic) run(a *Article) error {
	if e.checkOpenGraph(a) {
		return nil
	}

	e.findBestImage(a)

	return nil
}

func (e extractComic) setImage(a *Article, img *goquery.Selection) bool {
	if img.Length() == 0 {
		return false
	}

	img = img.First()

	src, ok := img.Attr("src")
	if !ok {
		return false
	}

	i := hitImage(src)
	if i == nil {
		return false
	}

	title, _ := img.Attr("title")
	if title == "" {
		title, _ = img.Attr("alt")
	}

	a.Img = i
	a.CleanedText = title
	a.TopNode = goquery.NewDocumentFromNode(&html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Span,
		Data:     "span",
	}).AppendSelection(img.Clone())

	return true
}

// If the opengraph image exists on the page, that's probably the comic
func (e extractComic) checkOpenGraph(a *Article) bool {
	ogimg := a.Meta.OpenGraph["image"]
	if ogimg == "" {
		return false
	}

	m, err := cascadia.Compile(fmt.Sprintf("img[src=\"%s\"]", ogimg))
	if err != nil {
		return false
	}

	return e.setImage(a, a.Doc.FindMatcher(m))
}

func (e extractComic) findBestImage(a *Article) bool {
	a.TopNode = a.Doc.Selection
	eImgs := extractImages{}
	eImgs.run(a)

	if a.Img != nil {
		return e.setImage(a, a.Img.Sel)
	}

	return false
}
