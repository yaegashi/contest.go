package main

import (
	"log"

	"github.com/spf13/cobra"
)

type App struct {
}

func (app *App) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "contest-cli",
		Short:        "Contest CLI for Gophers",
		SilenceUsage: true,
	}
	return cmd
}

func (app *App) Log(args ...interface{}) {
	log.Print(args...)
}

func (app *App) Logln(args ...interface{}) {
	log.Println(args...)
}

func (app *App) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
