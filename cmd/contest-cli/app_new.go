package main

import (
	"fmt"

	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
	"github.com/yaegashi/contest.go/tester"
)

type AppNew struct {
	*App
}

func (app *App) AppNewCmder() cmder.Cmder {
	return &AppNew{App: app}
}

func (app *AppNew) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "new <path> ...",
		Short:        "Create new solution dirs",
		Args:         cobra.MinimumNArgs(1),
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	return cmd
}

func (app *AppNew) RunE(cmd *cobra.Command, args []string) error {
	errs := 0
	for _, arg := range args {
		err := tester.CreateDirectory(arg, nil)
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
