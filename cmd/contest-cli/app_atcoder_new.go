package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
	"github.com/yaegashi/contest.go/atcoder"
	"github.com/yaegashi/contest.go/tester"
)

type AppAtCoderNew struct {
	*AppAtCoder
	Dir string
}

func (app *AppAtCoder) AtCoderNewCmder() cmder.Cmder {
	return &AppAtCoderNew{AppAtCoder: app}
}

func (app *AppAtCoderNew) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "new <contest ID or URL> [task ID]",
		Short:        "Prepare solution dirs for AtCoder contests",
		Args:         cobra.RangeArgs(1, 2),
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.Dir, "dir", "d", "", "Output dir")
	return cmd
}

func (app *AppAtCoderNew) RunE(cmd *cobra.Command, args []string) error {
	app.LoadToken()
	defer app.SaveToken()

	contestID := ""
	taskID := ""
	if strings.HasPrefix(args[0], atcoder.ContestsURL) {
		dirs := strings.Split(strings.Trim(args[0][len(atcoder.ContestsURL):], "/"), "/")
		contestID = dirs[0]
		if len(args) != 1 {
			return fmt.Errorf("Specifying task ID is not allowed with URL")
		}
	} else {
		contestID = args[0]
		if len(args) > 1 {
			taskID = args[1]
		}
	}
	if contestID == "" {
		return fmt.Errorf("Bad AtCoder contest: %s", args[0])
	}

	app.Logf("I: Fetching AtCoder contest %s", contestID)
	outContest, err := app.Client.FetchContest(contestID)
	if err != nil {
		return err
	}

	var outTask *atcoder.Task
	if strings.HasPrefix(args[0], atcoder.ContestsURL) {
		ok := outContest.URL == args[0]
		for _, task := range outContest.Tasks {
			if task.URL == args[0] {
				ok = true
				outTask = task
				break
			}
		}
		if !ok {
			return fmt.Errorf("Bad AtCoder URL: %s", args[0])
		}
	} else if taskID != "" {
		ok := false
		for _, task := range outContest.Tasks {
			if task.ID == taskID {
				ok = true
				outTask = task
				break
			}
		}
		if !ok {
			return fmt.Errorf("Bad AtCoder Task: %s", strings.Join(args, " "))
		}
	}

	if outTask != nil {
		dir := filepath.Join(outContest.ID, outTask.ID)
		if app.Dir != "" {
			dir = app.Dir
		}
		err = app.CreateTaskDir(dir, outTask)
		if err != nil {
			return err
		}
	} else {
		dir := contestID
		if app.Dir != "" {
			dir = app.Dir
		}
		for _, task := range outContest.Tasks {
			err = app.CreateTaskDir(filepath.Join(dir, task.ID), task)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *AppAtCoderNew) CreateTaskDir(outDir string, task *atcoder.Task) error {
	opts := &tester.DirectoryOptions{
		OmitSample: true,
		Preamble:   fmt.Sprintf("// Solution for %s", task.URL),
	}
	err := tester.CreateDirectory(outDir, opts)
	if err != nil {
		return err
	}
	for i, testcase := range task.Testcases {
		fin := filepath.Join(outDir, fmt.Sprintf("sample%d.in.txt", i+1))
		app.Logf("I: Created %s", fin)
		ioutil.WriteFile(fin, []byte(testcase.Input), 0644)
		fout := filepath.Join(outDir, fmt.Sprintf("sample%d.out.txt", i+1))
		app.Logf("I: Created %s", fout)
		ioutil.WriteFile(fout, []byte(testcase.Output), 0644)
	}
	return nil
}
