package swan

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	contentDirs = []string{
		"test_data/processors/comics/",
		"test_data/processors/default/",
	}
)

func TestProcessors(t *testing.T) {
	t.Parallel()
	hijiackHTTP()

	for _, dir := range contentDirs {
		filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					t.Fatalf("error while walking directory %s: %s",
						dir,
						err)
				}

				if info.IsDir() || !strings.HasSuffix(path, ".in") {
					return nil
				}

				baseName := strings.Replace(path, ".in", "", -1)
				htmlOut := fmt.Sprintf("%s.html.out", baseName)
				textOut := fmt.Sprintf("%s.text.out", baseName)

				in, err := ioutil.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read input for %s: %s", baseName, err)
				}

				html, htmlErr := ioutil.ReadFile(htmlOut)
				text, textErr := ioutil.ReadFile(textOut)
				if htmlErr != nil && textErr != nil {
					t.Fatalf("%s: \n"+
						"	htmlErr: %s\n"+
						"	textErr: %s",
						baseName,
						htmlErr,
						textErr)
				}

				parts := strings.SplitN(string(in), "\n", 2)

				a, err := FromHTML(parts[0], parts[1])
				if err != nil {
					t.Fatalf("%s: %s", baseName, err)
				}

				aHTML := ""
				if a.TopNode != nil {
					aHTML, _ = a.TopNode.Html()
					aHTML = strings.TrimSpace(aHTML)
				}
				eHTML := strings.TrimSpace(string(html))

				if htmlErr == nil && eHTML != aHTML {
					t.Fatalf(
						"%s: HTML does not match:\n"+
							"	Got:      %s\n"+
							"	Expected: %s",
						baseName,
						aHTML,
						eHTML)
				}

				eText := strings.TrimSpace(string(text))
				if textErr == nil && eText != a.CleanedText {
					t.Fatalf(
						"%s: CleanedText does not match:\n"+
							"	Got:      %s\n"+
							"	Expected: %s",
						baseName,
						a.CleanedText,
						eText)
				}

				return nil
			})
	}
}
