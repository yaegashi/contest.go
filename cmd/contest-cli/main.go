package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const mainGo = `package main

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

const mainTestGo = `package main

import (
	"io"
	"testing"

	"github.com/yaegashi/contest.go/tester"
)

func TestContest(t *testing.T) {
	tester.Run(t, "*.txt", func(in io.Reader, out io.Writer) error {
		con := &contest{in: in, out: out}
		return con.main()
	})
}
`

const testCaseTxt = "go\n--\nhello, go\n"

func main() {
	for _, dir := range os.Args[1:] {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created %s", dir)

		files := []struct {
			path, content string
		}{
			{"main.go", mainGo},
			{"main_test.go", mainTestGo},
			{"testcase00.txt", testCaseTxt},
		}
		for _, file := range files {
			path := filepath.Join(dir, file.path)
			err := ioutil.WriteFile(path, []byte(file.content), 0644)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Created %s", path)
		}
	}
}
