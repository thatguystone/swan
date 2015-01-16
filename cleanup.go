package swan

import (
	"fmt"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type commentMatcher struct{}
type childTextMatcher struct{}

var (
	remove   = []cascadia.Selector{}
	badNames = []string{
		"ajoutVideo",
		"articleheadings",
		"author-dropdown",
		"breadcrumbs",
		"byline",
		"cnn_html_slideshow",
		"cnn_strycaptiontxt",
		"cnn_strylftcntnt",
		"cnn_stryspcvbx",
		"cnnStryHghLght",
		"combx",
		"comment",
		"communitypromo",
		"contact",
		"contentTools2",
		"date",
		"foot",
		"footer",
		"Footer",
		"footnote",
		"inline-share-tools",
		"js_replies",
		"konafilter",
		"KonaFilter",
		"legende",
		"mediaarticlerelated",
		"menucontainer",
		"navbar",
		"pagetools",
		"PopularQuestions",
		"popup",
		"post-attributes",
		"retweet",
		"runaroundLeft",
		"shoutbox",
		"socialnetworking",
		"socialNetworking",
		"socialtools",
		"sponsor",
		"storytopbar-bucket",
		"subscribe",
		"tags",
		"the_answers",
		"timestamp",
		"tools",
		"utility-bar",
		"vcard",
		"welcome_form",
		"wp-caption-text",

		"facebook",
		"facebook-broadcasting",
		"google",
		"twitter",
	}
	badNamesExact = []string{
		"caption",
		"fn",
		"inset",
		"links",
		"print",
		"side",
	}
	badNamesStartsWith = []string{
		"more",
	}
	badNamesEndsWith = []string{
		"meta",
	}

	divSpanTags     = cascadia.MustCompile("div, span")
	emTags          = cascadia.MustCompile("em")
	imgTags         = cascadia.MustCompile("img")
	safeTags        = cascadia.MustCompile("body, article")
	scriptStyleTags = cascadia.MustCompile("script, style")
	unwraps         = cascadia.MustCompile("span[class~=dropcap]," +
		"span[class~=drop_cap]," +
		"p span")
	keepTags = cascadia.MustCompile("a, blockquote, dl, div," +
		"img, ol, p, pre, table, ul")
)

func init() {
	attrs := []string{
		"id",
		"class",
		"name",
	}

	for _, attr := range attrs {
		for _, s := range badNames {
			sel := fmt.Sprintf("[%s*=%s]", attr, s)
			remove = append(remove, cascadia.MustCompile(sel))
		}

		for _, s := range badNamesExact {
			sel := fmt.Sprintf("[%s=%s]", attr, s)
			remove = append(remove, cascadia.MustCompile(sel))
		}

		for _, s := range badNamesStartsWith {
			sel := fmt.Sprintf("[%s^=%s]", attr, s)
			remove = append(remove, cascadia.MustCompile(sel))
		}

		for _, s := range badNamesEndsWith {
			sel := fmt.Sprintf("[%s$=%s]", attr, s)
			remove = append(remove, cascadia.MustCompile(sel))
		}
	}
}

func (m commentMatcher) Match(n *html.Node) bool {
	return n.Type == html.CommentNode
}

func (m commentMatcher) matchAll(n *html.Node, ns []*html.Node) []*html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.CommentNode {
			ns = append(ns, c)
		} else {
			ns = m.matchAll(c, ns)
		}
	}

	return ns
}

func (m commentMatcher) MatchAll(n *html.Node) []*html.Node {
	return m.matchAll(n, nil)
}

func (m commentMatcher) Filter(n []*html.Node) []*html.Node {
	return nil
}

func (m childTextMatcher) Match(n *html.Node) bool {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode || m.Match(c) {
			return true
		}
	}

	return false
}

func (m childTextMatcher) MatchAll(n *html.Node) []*html.Node {
	return nil
}

func (m childTextMatcher) Filter(n []*html.Node) []*html.Node {
	var ns []*html.Node

	for _, c := range n {
		if m.Match(c) {
			ns = append(ns, c)
		}
	}

	return ns
}

func getReplacements(s *goquery.Selection) []*html.Node {
	var ns []*html.Node

	cs := s.FilterMatcher(childTextMatcher{})

	for _, n := range cs.Nodes {
		ns = append(ns, n)
	}

	return ns
}

func divToPara(i int, s *goquery.Selection) {
	if s.FindMatcher(keepTags).Length() == 0 {
		s.Nodes[0].Data = "p"
		s.Nodes[0].DataAtom = atom.P
	} else {
		ns := getReplacements(s.Empty())
		s.AppendNodes(ns...)
	}
}

// Cleanup runs basic article cleanup, discarding elements that are typically
// junk while regrouping text into larger chunks.
func (a *Article) Cleanup() {
	a.Doc.FindMatcher(safeTags).
		RemoveAttr("class").
		RemoveAttr("id").
		RemoveAttr("name")

	a.Doc.FindMatcher(commentMatcher{}).Remove()

	for _, cs := range remove {
		a.Doc.FindMatcher(cs).Remove()
	}

	a.Doc.FindMatcher(scriptStyleTags).Remove()
	a.Doc.FindMatcher(unwraps).Unwrap()

	ems := a.Doc.FindMatcher(emTags)
	ems.NotSelection(ems.HasMatcher(imgTags)).Unwrap()

	a.Doc.FindMatcher(divSpanTags).Each(divToPara)
}
