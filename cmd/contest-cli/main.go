package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yaegashi/contest.go/tester"
)

func main() {
	for _, arg := range os.Args[1:] {
		var err error
		if strings.HasPrefix(arg, AtCoderBaseURL) {
			dir := strings.Trim(arg[len(AtCoderBaseURL):], "/")
			problem := strings.Split(dir, "/")[0]
			if problem == "" {
				log.Printf("E: Wrong AtCoder URL: %s", arg)
				continue
			}
			tasksURL := AtCoderBaseURL + filepath.Join(problem, "tasks")
			err = createAtCoderContest(problem, tasksURL)
		} else {
			err = tester.CreateDirectory(arg, nil)
		}
		if err != nil {
			log.Printf("E: Failed to process %s: %s", arg, err)
		}
	}
}
