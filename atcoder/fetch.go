package atcoder

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (cli *Client) FetchContest(id string) (*Contest, error) {
	contest := &Contest{
		ID:  id,
		URL: ContestsURL + "/" + id,
	}
	tasksURL, err := url.Parse(contest.URL + "/tasks")
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(tasksURL.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch %s: %s", tasksURL, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	doc.Find("div > h2:first-child").Each(func(i int, s1 *goquery.Selection) {
		if s1.Text() != "Tasks" {
			return
		}
		s1.Parent().Find("tbody tr").Each(func(i int, s2 *goquery.Selection) {
			a := s2.Find("td").First().Find("a")
			p := a.Text()
			href, ok := a.Attr("href")
			if p == "" || !ok {
				return
			}
			taskURL, err := tasksURL.Parse(href)
			if err != nil {
				return
			}
			taskID := strings.ToLower(p)
			task, err := cli.FetchTask(taskID, taskURL.String())
			if err != nil {
				return
			}
			contest.Tasks = append(contest.Tasks, task)
		})
	})
	return contest, nil
}

func (cli *Client) FetchTask(id string, u string) (*Task, error) {
	task := &Task{
		ID:  id,
		URL: u,
	}
	res, err := cli.Get(task.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch %s: %s", task.URL, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	inputs := []string{}
	outputs := []string{}
	doc.Find("section").Each(func(i int, s *goquery.Selection) {
		h3 := s.Find("h3").First().Text()
		if strings.HasPrefix(h3, "Sample Input ") {
			inputs = append(inputs, s.Find("pre").First().Text())
		}
		if strings.HasPrefix(h3, "Sample Output ") {
			outputs = append(outputs, s.Find("pre").First().Text())
		}
	})
	if len(inputs) == len(outputs) {
		for i := 0; i < len(inputs); i++ {
			task.Testcases = append(task.Testcases, &Testcase{Input: inputs[i], Output: outputs[i]})
		}
	}
	return task, nil
}
