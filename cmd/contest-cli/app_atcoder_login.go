package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
	"golang.org/x/crypto/ssh/terminal"
)

type AppAtCoderLogin struct {
	*AppAtCoder
	Username string
	Password string
}

func (app *AppAtCoder) AtCoderLoginCmder() cmder.Cmder {
	return &AppAtCoderLogin{AppAtCoder: app}
}

func (app *AppAtCoderLogin) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "login",
		Short:        "Log in and preserve auth tokens for AtCoder contests",
		RunE:         app.RunE,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.Username, "username", "u", "", "username (env:CONTEST_ATCODER_USERNAME)")
	cmd.Flags().StringVarP(&app.Password, "password", "p", "", "password (env:CONTEST_ATCODER_PASSWORD)")
	return cmd
}

func (app *AppAtCoderLogin) RunE(cmd *cobra.Command, args []string) error {
	if app.Username == "" {
		app.Username = os.Getenv("CONTEST_ATCODER_USERNAME")
	}
	if app.Password == "" {
		app.Password = os.Getenv("CONTEST_ATCODER_PASSWORD")
	}

	if app.Username == "" {
		fmt.Fprintf(os.Stderr, "AtCoder Login: ")
		fmt.Scan(&app.Username)
	}
	if app.Password == "" {
		fmt.Fprintf(os.Stderr, "Password: ")
		passwordBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return err
		}
		app.Password = string(passwordBytes)
	}

	if app.Username == "" || app.Password == "" {
		return fmt.Errorf("Specify username and password")
	}

	err := app.Client.Login(app.Username, app.Password)
	if err != nil {
		return err
	}

	err = app.SaveToken()
	if err != nil {
		return err
	}

	return nil
}
