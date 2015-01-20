package swan

import (
	"code.google.com/p/cascadia"
	// "github.com/PuerkitoBio/goquery"
)

var (
	linkMatcher = cascadia.MustCompile("a")
)

func extractLinks(a *Article) error {
	// Enable once there's a top-rated node to search through

	// a.TopNode.FindMatcher(linkMatcher).Each(func(i int, s *goquery.Selection) {
	// 	h, exists := s.Attr("href")
	// 	if exists && h != "" {
	// 		a.Meta.Links = append(a.Meta.Links, h)
	// 	}
	// })

	return nil
}
