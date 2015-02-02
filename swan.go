// Package swan implementats the Goose HTML Content / Article Extractor algorithm, with some extra, pretty goodies
package swan

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	// Version of the library
	Version = "1.0"
)

var (
	client = &http.Client{
		Timeout: time.Second * 10,
	}
)

// FromURL does its best to extract an article from the given URL
func FromURL(url string) (a Article, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("could not create new request: %s", err)
		return
	}

	req.Header.Set("User-Agent", "swan/"+Version)
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("could not load URL: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("could not read response body: %s", err)
		return
	}

	fmt.Println(http.DetectContentType(body))

	return FromHTML(string(body))
}

// FromHTML does its best to extract an article from a single HTML page
func FromHTML(html string) (Article, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("invalid HTML: %s", err)
		return Article{}, err
	}

	return FromDoc(doc)
}

// FromDoc does its best to extract an article from a single document
func FromDoc(doc *goquery.Document) (Article, error) {
	a := Article{
		Doc: doc,
	}

	err := a.extract()
	if err != nil {
		return Article{}, err
	}

	return a, nil
}
