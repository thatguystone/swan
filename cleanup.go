package swan

import (
	"bytes"
	"fmt"
	"strings"

	"code.google.com/p/cascadia"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type cleanup struct{}
type commentMatcher struct{}
type childTextMatcher struct{}

var (
	tablinesReplacer = strings.NewReplacer("\n", "\n\n", "\t", "")

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

func (m childTextMatcher) matchAll(n *html.Node, ns []*html.Node) []*html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			ns = append(ns, c)
		} else {
			ns = m.matchAll(c, ns)
		}
	}

	return ns
}

func (m childTextMatcher) MatchAll(n *html.Node) []*html.Node {
	return m.matchAll(n, nil)
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

func nodeIs(n *html.Node, a atom.Atom) bool {
	return n != nil && n.Type == html.ElementNode && n.DataAtom == a
}

func getReplacements(s *goquery.Selection) []*html.Node {
	var ns []*html.Node
	var b bytes.Buffer

	flushB := func() {
		ns = append(ns, &html.Node{
			Type: html.TextNode,
			Data: b.String(),
		})
		b.Reset()
	}

	cs := s.FindMatcher(childTextMatcher{})

	for _, n := range cs.Nodes {
		switch {
		case n.Type == html.TextNode:
			n.Data = tablinesReplacer.Replace(n.Data)
			n.Data = strings.TrimSpace(n.Data)

			if len(n.Data) == 0 {
				continue
			}

			an := n
			flushB()

			// Rewind to first <a> before this text
			for ; nodeIs(an.PrevSibling, atom.A); an = an.PrevSibling {
			}

			// Run through all previous and trailing <a>s, injecting the node in
			// the mix
			for ; nodeIs(an, atom.A) || an == n; an = an.NextSibling {
				if an.Parent != nil {
					an.Parent.RemoveChild(an)
				}

				ns = append(ns, an)
			}

		case nodeIs(n, atom.P) && b.Len() > 0:
			flushB()
			fallthrough
		default:
			ns = append(ns, n)
		}
	}

	flushB()

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

func (c cleanup) run(a *Article) error {
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

	return nil
}