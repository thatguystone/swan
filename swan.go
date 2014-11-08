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
	Version = "1.0"
)

var (
	client = &http.Client{
		Timeout: time.Second * 10,
	}
)

func FromUrl(url string, cfg Config) (a Article, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("Could not create new request: %s", err)
		return
	}

	req.Header.Set("User-Agent", "swan/"+Version)
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Could not load URL: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Could not read response body: %s", err)
		return
	}

	fmt.Println(http.DetectContentType(body))

	a, err = FromHtml(string(body), cfg)
	return
}

func FromHtml(html string, cfg Config) (Article, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("Invalid HTML: %s", err)
		return Article{}, err
	}

	return FromDoc(doc, cfg)
}

func FromDoc(doc *goquery.Document, cfg Config) (Article, error) {
	return prepareArticle(doc, cfg)
}
