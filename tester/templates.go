package tester

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const MainGoTemplate = `package main

{{if .Preamble}}{{.Preamble}}{{end}}

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type contest struct {
	in  io.Reader
	out io.Writer
}

func (con *contest) Scan(a ...interface{}) (int, error) {
	return fmt.Fscan(con.in, a...)
}
func (con *contest) Scanln(a ...interface{}) (int, error) {
	return fmt.Fscanln(con.in, a...)
}
func (con *contest) Scanf(f string, a ...interface{}) (int, error) {
	return fmt.Fscanf(con.in, f, a...)
}
func (con *contest) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(con.out, a...)
}
func (con *contest) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(con.out, a...)
}
func (con *contest) Printf(f string, a ...interface{}) (int, error) {
	return fmt.Fprintf(con.out, f, a...)
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	con := &contest{in: in, out: out}
	con.main()
}

func (con *contest) main() error {
	var s string
	con.Scan(&s)
	con.Println("hello,", s)
	return nil
}
`

const MainTestGoTemplate = `package main

import (
	"io"
	"testing"

	"github.com/yaegashi/contest.go/tester"
)

func TestContest(t *testing.T) {
	tester.Run(t, "*.in.txt", func(in io.Reader, out io.Writer) error {
		con := &contest{in: in, out: out}
		return con.main()
	})
}
`

const SampleInTemplate = "go\n"

const SampleOutTemplate = "hello, go\n"

type DirectoryOptions struct {
	Preamble   string
	OmitSample bool
}

func CreateDirectory(dir string, opts *DirectoryOptions) error {
	if opts == nil {
		opts = &DirectoryOptions{}
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	files := []struct {
		path, content string
		sample        bool
	}{
		{"main.go", MainGoTemplate, false},
		{"main_test.go", MainTestGoTemplate, false},
		{"sample1.in.txt", SampleInTemplate, true},
		{"sample1.out.txt", SampleOutTemplate, true},
	}
	for _, file := range files {
		if opts.OmitSample && file.sample {
			continue
		}
		tmpl, err := template.New(file.path).Parse(file.content)
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, opts)
		if err != nil {
			return err
		}
		b := buf.Bytes()
		if !file.sample {
			b, err = format.Source(b)
			if err != nil {
				return err
			}
		}
		path := filepath.Join(dir, file.path)
		err = ioutil.WriteFile(path, b, 0644)
		if err != nil {
			return err
		}
		log.Printf("I: Created %s", path)
	}
	return nil
}
