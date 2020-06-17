package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yaegashi/contest.go/tester"
)

const AtCoderBaseURL = "https://atcoder.jp/contests/"

func createAtCoderTask(outDir string, taskURL string) error {
	res, err := http.Get(taskURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to fetch %s: %s", taskURL, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	opts := &tester.DirectoryOptions{
		OmitSample: true,
		Preamble:   fmt.Sprintf("// Solution for %s", taskURL),
	}
	err = tester.CreateDirectory(outDir, opts)
	if err != nil {
		return err
	}
	doc.Find("section").Each(func(i int, s *goquery.Selection) {
		h3 := s.Find("h3").First().Text()
		if strings.HasPrefix(h3, "Sample Input ") {
			pre := s.Find("pre").First().Text()
			fn := filepath.Join(outDir, fmt.Sprintf("sample%s.in.txt", h3[13:]))
			ioutil.WriteFile(fn, []byte(pre), 0644)
			log.Println("I: Created", fn)
		}
		if strings.HasPrefix(h3, "Sample Output ") {
			pre := s.Find("pre").First().Text()
			fn := filepath.Join(outDir, fmt.Sprintf("sample%s.out.txt", h3[14:]))
			ioutil.WriteFile(fn, []byte(pre), 0644)
			log.Println("I: Created", fn)
		}
	})
	return nil
}

func createAtCoderContest(outDir string, tasksURL string) error {
	u1, err := url.Parse(tasksURL)
	if err != nil {
		return err
	}
	res, err := http.Get(tasksURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to fetch %s: %s", tasksURL, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	doc.Find("div > h2:first-child").Each(func(i int, s1 *goquery.Selection) {
		if s1.Text() != "Tasks" {
			return
		}
		s1.Parent().Find("tbody tr").Each(func(i int, s2 *goquery.Selection) {
			a := s2.Find("td").First().Find("a")
			p := a.Text()
			href, ok := a.Attr("href")
			if p == "" || !ok {
				return
			}
			u2, err := url.Parse(href)
			if err != nil {
				return
			}
			p = strings.ToLower(p) + "1"
			u3 := u1.ResolveReference(u2)
			createAtCoderTask(filepath.Join(outDir, p), u3.String())
		})
	})
	return nil
}
