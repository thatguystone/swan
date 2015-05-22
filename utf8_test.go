package swan

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToUtf8(t *testing.T) {
	dir := "test_data/utf8/"
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				t.Fatalf("error while walking directory %s: %s", dir, err)
			}

			if info.IsDir() || !strings.HasSuffix(path, ".in") {
				return nil
			}

			in, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read %s: %s", path, err)
			}

			name := strings.Replace(path, ".in", ".out", -1)
			out, err := ioutil.ReadFile(name)
			if err != nil {
				t.Fatalf("failed to read %s: %s", path, err)
			}

			uout, err := ToUtf8(in)
			if err != nil {
				t.Fatalf("failed to convert %s: %s", path, err)
			}

			if !bytes.Equal(out, uout) {
				t.Fatalf("conversion doesn't match for %s:\n\t%s\n\t%s",
					path,
					out,
					uout)
			}

			return nil
		})
}
