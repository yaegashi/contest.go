package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	defaultConfigDir = "~/.contest-cli"
)

type App struct {
	ConfigDir string
}

func (app *App) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "contest-cli",
		Short:             "Contest CLI for Gophers",
		PersistentPreRunE: app.PersistentPreRunE,
		SilenceUsage:      true,
	}
	cmd.PersistentFlags().StringVarP(&app.ConfigDir, "config-dir", "", "", "config dir (env:CONTEST_CONFIG_DIR, default:\"~/.contest-cli\")")
	return cmd
}

func (app *App) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	if app.ConfigDir == "" {
		app.ConfigDir = os.Getenv("CONTEST_CONFIG_DIR")
	}
	if app.ConfigDir == "" {
		app.ConfigDir = defaultConfigDir
	}
	dir, err := homedir.Expand(app.ConfigDir)
	if err != nil {
		return err
	}
	app.ConfigDir = dir
	return nil
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

func (app *App) ConfigFile(path string) string {
	return filepath.Join(app.ConfigDir, path)
}

func (app *App) ReadConfigFile(path string) ([]byte, error) {
	return ioutil.ReadFile(app.ConfigFile(path))
}

func (app *App) WriteConfigFile(path string, b []byte, p os.FileMode) error {
	if _, err := os.Stat(app.ConfigDir); err != nil {
		err = os.MkdirAll(app.ConfigDir, 0755)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(app.ConfigFile(path), b, p)
}
