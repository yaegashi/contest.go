package tester_test

import (
	"io"
	"testing"

	"github.com/yaegashi/contest.go/tester"
)

func TestTester(t *testing.T) {
	tester.Run(t, "testdata/*.inout.txt", func(in io.Reader, out io.Writer) error {
		_, err := io.Copy(out, in)
		return err
	})
	tester.Run(t, "testdata/*.in.txt", func(in io.Reader, out io.Writer) error {
		_, err := io.Copy(out, in)
		return err
	})
}
