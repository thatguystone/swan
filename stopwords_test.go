package swan

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestDetectLang(t *testing.T) {
	t.Parallel()

	for l := range stopwords {
		path := fmt.Sprintf("test_data/stopwords/%s.txt", l)
		txt, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %s", path, err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(txt)))
		if err != nil {
			t.Fatalf("failed to create doc: %s", err)
		}

		a := &Article{
			Doc: doc,
		}

		err = cleanup{}.run(a)
		if err != nil {
			t.Fatalf("failed to clean doc: %s", err)
		}

		lang := detectLang(a)
		path = strings.Replace(path, ".txt", "", -1)
		if !strings.HasSuffix(path, lang) {
			t.Fatalf("incorrect language detected for %s: %s", path, lang)
		}
	}
}

func TestSplitText(t *testing.T) {
	t.Parallel()

	type test struct {
		in  string
		out []string
	}

	table := []test{
		test{
			in:  "there once  was a boy    .",
			out: []string{"there", "once", "was", "a", "boy"},
		},
		test{
			in:  "the boy's hat was green",
			out: []string{"the", "boy's", "hat", "was", "green"},
		},
		test{
			in:  "spaces.            ",
			out: []string{"spaces"},
		},
		test{
			in:  "                     more spaces.            ",
			out: []string{"more", "spaces"},
		},
		test{
			in:  "punct: everywhere!",
			out: []string{"punct", "everywhere"},
		},
	}

	for _, tc := range table {
		ws := splitText(tc.in)

		if !reflect.DeepEqual(ws, tc.out) {
			t.Fatalf("%#v != %#v", ws, tc.out)
		}
	}
}
