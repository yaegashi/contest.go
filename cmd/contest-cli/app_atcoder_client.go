package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"
	"github.com/yaegashi/contest.go/atcoder"
)

type AppAtCoderClient struct {
	*AppAtCoder
	Output  string
	Data    string
	Forms   map[string]string
	Headers map[string]string
}

func (app *AppAtCoder) AtCoderClientCmder() cmder.Cmder {
	return &AppAtCoderClient{AppAtCoder: app}
}

func (app *AppAtCoderClient) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "client <method> <URL>",
		Short:        "AtCoder debug client",
		Args:         cobra.ExactArgs(2),
		RunE:         app.RunE,
		Hidden:       true,
		SilenceUsage: true,
	}
	cmd.Flags().StringVarP(&app.Output, "output", "o", "", "output file")
	cmd.Flags().StringVarP(&app.Data, "data", "d", "", "input file")
	cmd.Flags().StringToStringVarP(&app.Forms, "form", "F", map[string]string{}, "form")
	cmd.Flags().StringToStringVarP(&app.Headers, "header", "H", map[string]string{}, "header")
	return cmd
}

func (app *AppAtCoderClient) RunE(cmd *cobra.Command, args []string) error {
	method := strings.ToUpper(args[0])
	address := args[1]
	if !strings.HasPrefix(address, atcoder.BaseURL) {
		return fmt.Errorf("Unsupported URL")
	}

	app.LoadToken()
	defer app.SaveToken()

	header := http.Header{}
	for k, v := range app.Headers {
		header.Add(k, v)
	}

	var body io.Reader
	switch method {
	case http.MethodPost:
		if len(app.Forms) > 0 {
			f := url.Values{}
			for k, v := range app.Forms {
				f.Add(k, v)
			}
			body = strings.NewReader(f.Encode())
			header.Add("Content-Type", "application/x-www-form-urlencoded")
			break
		}
		fallthrough
	case http.MethodPut:
		if app.Data == "-" {
			body = os.Stdin
		} else if app.Data != "" {
			file, err := os.Open(app.Data)
			if err != nil {
				return err
			}
			defer file.Close()
			body = file
		}
	}

	req, err := http.NewRequest(method, address, body)
	if err != nil {
		return err
	}

	req.Header = header

	res, err := app.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out io.Writer
	if app.Output == "" || app.Output == "-" {
		out = os.Stdout
	} else {
		file, err := os.Create(app.Output)
		if err != nil {
			return err
		}
		defer file.Close()
		out = file
	}
	io.Copy(out, res.Body)

	return nil
}
