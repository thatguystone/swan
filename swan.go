// Package swan implementats the Goose HTML Content / Article Extractor
// algorithm
package swan

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// Version of the library
	Version = "1.0"
)

// FromURL does its best to extract an article from the given URL
func FromURL(url string) (a Article, err error) {
	body, resp, err := httpGet(url)
	if err != nil {
		return
	}

	defer body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("could not read response body: %s", err)
		return
	}

	return FromHTML(resp.Request.URL.String(), string(html))
}

// FromHTML does its best to extract an article from a single HTML page
func FromHTML(url string, html string) (Article, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("invalid HTML: %s", err)
		return Article{}, err
	}

	return FromDoc(url, doc)
}

// FromDoc does its best to extract an article from a single document
func FromDoc(url string, doc *goquery.Document) (Article, error) {
	a := Article{
		URL: url,
		Doc: doc,
	}

	err := a.extract()
	if err != nil {
		return Article{}, err
	}

	return a, nil
}
