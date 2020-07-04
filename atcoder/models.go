package atcoder

const (
	BaseURL     = "https://atcoder.jp"
	LoginURL    = "https://atcoder.jp/login"
	ContestsURL = "https://atcoder.jp/contests"
)

type Contest struct {
	ID    string
	URL   string
	Title string
	Tasks []*Task
}

type Task struct {
	ID        string
	URL       string
	Title     string
	Testcases []*Testcase
}

type Testcase struct {
	Input  string
	Output string
}
