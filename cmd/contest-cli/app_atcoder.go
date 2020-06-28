package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
	"github.com/yaegashi/contest.go/tester"
)

const AtCoderBaseURL = "https://atcoder.jp/contests/"

type AppAtCoder struct {
	*App
	Dir string
}

func (app *App) AtCoderCmder() cmder.Cmder {
	return &AppAtCoder{App: app}
}

func (app *AppAtCoder) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "atcoder",
		Aliases:      []string{"ac"},
		Short:        "Sub commands for AtCoder https://atcoder.jp",
		SilenceUsage: true,
	}
	return cmd
}

type AppAtCoderLogin struct {
	*AppAtCoder
}

func (app *AppAtCoder) AtCoderLoginCmder() cmder.Cmder {
	return &AppAtCoderLogin{AppAtCoder: app}
}

func (app *AppAtCoderLogin) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "login",
		Short:        "Log in and preserve auth cookies for AtCoder contests",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	return cmd
}

func (app *AppAtCoderLogin) RunE(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not yet implemented")
}

type AppAtCoderNew struct {
	*AppAtCoder
}

func (app *AppAtCoder) AtCoderNewCmder() cmder.Cmder {
	return &AppAtCoderNew{AppAtCoder: app}
}

func (app *AppAtCoderNew) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "new <contest name or URL> ...",
		Short:        "Prepare solution dirs for AtCoder contests",
		Args:         cobra.MinimumNArgs(1),
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.Dir, "dir", "d", "", "Output dir")
	return cmd
}

func (app *AppAtCoderNew) RunE(cmd *cobra.Command, args []string) error {
	errs := 0
	for _, arg := range args {
		var err error
		dir := ""
		contest := ""
		if strings.HasPrefix(arg, AtCoderBaseURL) {
			contest = strings.Split(strings.Trim(arg[len(AtCoderBaseURL):], "/"), "/")[0]
			dir = contest
		} else {
			dirs := strings.Split(arg, "/")
			for i := len(dirs) - 1; contest == "" && i >= 0; i-- {
				contest = strings.ToLower(dirs[i])
			}
			dir = arg
		}
		if contest == "" {
			app.Logf("E: %s: Bad AtCoder contest", arg)
			errs++
			continue
		}
		if app.Dir != "" {
			dir = app.Dir
		}
		contestURL := AtCoderBaseURL + contest
		tasksURL := contestURL + "/tasks"
		app.Logf("I: Fetching contest %s into %s", contest, dir)
		err = app.createAtCoderContest(dir, tasksURL)
		if err != nil {
			app.Logf("E: %s: %s", arg, err)
			errs++
		}
	}
	if errs > 0 {
		return fmt.Errorf("failed to create some dirs")
	}
	return nil
}

func (app *AppAtCoderNew) createAtCoderTask(outDir string, taskURL string) error {
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
			app.Logln("I: Created", fn)
		}
		if strings.HasPrefix(h3, "Sample Output ") {
			pre := s.Find("pre").First().Text()
			fn := filepath.Join(outDir, fmt.Sprintf("sample%s.out.txt", h3[14:]))
			ioutil.WriteFile(fn, []byte(pre), 0644)
			app.Logln("I: Created", fn)
		}
	})
	return nil
}

func (app *AppAtCoderNew) createAtCoderContest(outDir string, tasksURL string) error {
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
			p = strings.ToLower(p)
			u3 := u1.ResolveReference(u2)
			app.createAtCoderTask(filepath.Join(outDir, p), u3.String())
		})
	})
	return nil
}
