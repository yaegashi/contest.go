package atcoder

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func (cli *Client) Login(username string, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("Specify username and password")
	}

	res, err := cli.Get(LoginURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to get login: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	csrfToken := ""
	doc.Find("input").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if name, ok := s.Attr("name"); ok && name == "csrf_token" {
			if value, ok := s.Attr("value"); ok {
				csrfToken = value
				return false
			}
		}
		return true
	})
	if csrfToken == "" {
		return fmt.Errorf("Failed to get csrf_token")
	}

	values := url.Values{}
	values.Add("username", username)
	values.Add("password", password)
	values.Add("csrf_token", csrfToken)
	res, err = cli.PostForm(LoginURL, values)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to post login: %s", res.Status)
	}
	if res.Request.URL.String() == LoginURL {
		// No redirection probably means login failure
		return fmt.Errorf("Failed to login")
	}

	cli.Username = username

	return nil
}
