package swan

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type Result struct {
	URL      string
	Expected Expected
}

type ExpectedTopImage struct {
	Bytes      uint
	Confidence string `json:"confidence_score"`
	Height     uint
	Src        string
	Width      uint
}

type Expected struct {
	Authors         []string
	CleanedText     string `json:"cleaned_text"`
	Links           int
	MetaDescription string `json:"meta_description"`
	MetaFavicon     string `json:"meta_favicon"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaLang        string `json:"meta_lang"`
	OpenGraph       map[string]string
	PublishDate     string `json:"publish_date"`
	TopImage        ExpectedTopImage
	Tags            []string
	Title           string
}

var hijackOnce sync.Once

func hijiackHTTP() {
	hijackOnce.Do(func() {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sum := sha1.Sum([]byte(r.URL.String()))
			hash := hex.EncodeToString(sum[:])

			path := fmt.Sprintf("test_data/imgs/%s", hash)
			errPath := fmt.Sprintf("%s.err", path)

			_, dataErr := os.Stat(path)
			_, errErr := os.Stat(errPath)

			if os.IsNotExist(dataErr) && os.IsNotExist(errErr) {
				resp, err := http.Get(r.URL.String())
				if err != nil || resp.StatusCode != 200 {
					err = ioutil.WriteFile(errPath, nil, 0644)
					if err != nil {
						panic(err)
					}
				} else {
					defer resp.Body.Close()
					d, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						panic(err)
					}

					err = ioutil.WriteFile(path, d, 0644)
					if err != nil {
						panic(err)
					}
				}
			}

			if _, err := os.Stat(errPath); err == nil {
				http.Error(w, "not found", 404)
			} else {
				d, _ := ioutil.ReadFile(path)
				w.Write(d)
			}
		})

		s := httptest.NewServer(h)
		httpClient.Transport = http.DefaultTransport

		purl, _ := url.Parse(s.URL)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(purl),
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	})
}

func runPyTests(
	t *testing.T,
	dir string,
	fn func(t *testing.T, name string, a *Article, r *Result)) {

	hijiackHTTP()

	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				t.Fatalf("error while walking directory %s: %s", dir, err)
			}

			if info.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}

			r := Result{}
			name := strings.Replace(path, ".json", "", -1)

			jsonf, err := os.Open(path)
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			if err = json.NewDecoder(jsonf).Decode(&r); err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			h := strings.Replace(path, ".json", ".html", -1)
			htmlf, err := os.Open(h)
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			html, err := ioutil.ReadAll(htmlf)
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			a, err := FromHTML(r.URL, string(html))
			if err != nil {
				t.Fatalf("%s: %s", name, err)
			}

			fn(t, name, &a, &r)
			return nil
		})
}
