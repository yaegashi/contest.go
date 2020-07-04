package main

import (
	"encoding/json"

	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
	"github.com/yaegashi/contest.go/atcoder"
)

const (
	AtCoderConfigFile = "atcoder.json"
)

type AppAtCoder struct {
	*App
	Client *atcoder.Client
}

func (app *App) AtCoderCmder() cmder.Cmder {
	return &AppAtCoder{App: app}
}

func (app *AppAtCoder) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "atcoder",
		Aliases:           []string{"ac"},
		Short:             "Sub commands for AtCoder https://atcoder.jp",
		PersistentPreRunE: app.PersistentPreRunE,
		SilenceUsage:      true,
	}
	return cmd
}

func (app *AppAtCoder) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	err := app.App.PersistentPreRunE(cmd, args)
	if err != nil {
		return err
	}
	cli, err := atcoder.NewClient()
	if err != nil {
		return err
	}
	app.Client = cli
	return nil
}

func (app *AppAtCoder) LoadToken() error {
	b, err := app.ReadConfigFile(AtCoderConfigFile)
	if err != nil {
		return err
	}
	var token *atcoder.Token
	err = json.Unmarshal(b, &token)
	if err != nil {
		return err
	}
	app.Client.LoadToken(token)
	return nil
}

func (app *AppAtCoder) SaveToken() error {
	b, err := json.MarshalIndent(app.Client.SaveToken(), "", "  ")
	if err != nil {
		return err
	}
	return app.WriteConfigFile(AtCoderConfigFile, b, 0600)
}
