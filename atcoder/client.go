package atcoder

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	*http.Client
	Username string
}

type Token struct {
	Username string
	Cookies  []*http.Cookie
}

func NewClient() (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	return &Client{Client: &http.Client{Jar: jar}}, nil
}

func (cli *Client) LoadToken(token *Token) {
	cli.Username = token.Username
	u, _ := url.Parse(BaseURL)
	cli.Client.Jar.SetCookies(u, token.Cookies)
}

func (cli *Client) SaveToken() *Token {
	u, _ := url.Parse(BaseURL)
	cookies := []*http.Cookie{}
	for _, cookie := range cli.Client.Jar.Cookies(u) {
		switch cookie.Name {
		case "REVEL_SESSION":
			cookies = append(cookies, cookie)
		}
	}
	return &Token{
		Username: cli.Username,
		Cookies:  cookies,
	}
}
