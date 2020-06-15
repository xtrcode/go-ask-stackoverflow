package website

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type StackOverflow struct {
	content []byte
}

func (s *StackOverflow) Get(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	s.content, err = ioutil.ReadAll(resp.Body)
	return err
}

func (s *StackOverflow) Parse() (string, error) {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(s.content))
	if err != nil {
		return "", err
	}

	var matches []string

	// find the accepted answer
	document.Find(".accepted-answer").Each(func(i int, s *goquery.Selection) {
		matches = append(matches, s.Find(".post-text").Text())
		return
	})

	if matches != nil {
		return matches[0], nil
	}

	// find another answer
	document.Find("#answers > .answer").Each(func(i int, s *goquery.Selection) {
		matches = append(matches, s.Find(".post-text").Text())

		return
	})

	if matches != nil {
		return matches[0], nil
	}

	return "", nil
}
