package engine

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Google struct {
	content []byte
}

const SEARCH_ENGINE = "https://www.google.de/search"
const QUERY_PARAM = "q"

func (g *Google) Request(str string) error {
	var err error

	u, err := url.Parse(SEARCH_ENGINE)
	if err != nil {
		return err
	}

	// encode question
	param := url.Values{}
	param.Add(QUERY_PARAM, str)
	u.RawQuery = param.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// retrieve content/body
	g.content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (g *Google) Get() (string, error) {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(g.content))
	if err != nil {
		return "", err
	}

	var matches []string

	document.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok && strings.Contains(href, "/url?q=https://stackoverflow.com") {
			matches = append(matches, strings.Split(strings.Trim(href, "/url?q="), "&")[:1][0])
		}
	})

	return matches[0], nil
}
