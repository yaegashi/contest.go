package tester

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/diff"
)

var Delimiter = "--"

func Run(t *testing.T, pattern string, solve func(in io.Reader, out io.Writer) error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	for i, file := range files {
		t.Run(fmt.Sprintf("%d:%s", i, file), func(t *testing.T) {
			pass := 0
			input := ""
			want := []string{}
			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())
				if text == Delimiter {
					pass++
					continue
				}
				switch pass {
				case 0:
					input += text + "\n"
				case 1:
					want = append(want, text)
				}
			}
			err = scanner.Err()
			if err != nil {
				t.Fatal(err)
			}
			inBuf := bytes.NewBufferString(input)
			outBuf := &bytes.Buffer{}
			err = solve(inBuf, outBuf)
			if err != nil {
				t.Error(err)
			}
			got := strings.SplitAfter(outBuf.String(), "\n")
			if got[len(got)-1] == "" {
				got = got[:len(got)-1]
			}
			for j := range got {
				got[j] = strings.TrimRight(got[j], "\r\n")
			}
			ab := diff.Strings(want, got)
			e := diff.Myers(context.Background(), ab)
			if !e.IsIdentity() {
				diffBuf := &bytes.Buffer{}
				e.WriteUnified(diffBuf, ab, diff.Names("want", "got"), diff.TerminalColor())
				s := diffBuf.String()
				if strings.HasSuffix(s, "\n\x1b[0m") {
					s = s[:len(s)-5] + "\x1b[0m"
				}
				t.Error("Wrong answer:\n" + s)
			}
		})
	}
}
