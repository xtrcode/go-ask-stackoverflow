package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
)

const SEARCH_ENGINE = "https://www.google.de/search"
const QUERY_PARAM = "q"

type parseFunc func(*goquery.Document, *[]string)

func main() {
	if len(os.Args) == 1 {
		help()
	}

	question := strings.Join(os.Args[1:], " ")

	if answer := getCache(question); len(answer) != 0 {
		fmt.Println(answer)
		os.Exit(1)
	}

	content := googleRequest(question)

	google := parseGoogle(content)
	if len(google) == 0 {
		noAnswer(question)
	}

	content = get(google[0])
	if len(content) == 0 {
		noAnswer(question)
	}

	answer := parseStackOverflow(content)
	if len(answer) == 0 {
		noAnswer(question)
	}

	setCache(question, answer[0])
	fmt.Println(answer[0])
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// help prints the usage and other helpful information about the program
func help() {
	fmt.Println("ask is a simple terminal command that tries to find " +
		"the most popular answer to your question")
	fmt.Println("from StackOverflow.com")
	fmt.Println("\nusage: ask your question")
	os.Exit(0)
}

// noAnswer returns the sad message and caches that there were no answer
func noAnswer(question string) {
	fmt.Println("No answers found :(")
	setCache(question, "No answers found :x")
	os.Exit(0)
}

// openCache opens the bold database
// and returns its handle
func openCache() (*bolt.DB) {
	// create cache dir if not exists
	usr, err := user.Current()
	check(err)

	path := usr.HomeDir + "/.ask"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		check(os.MkdirAll(path, 0600))
	}

	// open connection to cache db
	db, err := bolt.Open(path+"/cache.db", 0600, nil)
	check(err)

	return db
}

// getCache returns the answer to the given question
// if it were cached previously
func getCache(question string) (string) {
	db := openCache()
	defer db.Close()

	var answer []byte
	db.View(func(tx *bolt.Tx) (error) {
		b := tx.Bucket([]byte("42"))
		if b == nil {
			return nil
		}

		answer = b.Get([]byte(question))

		return nil
	})

	return string(answer)
}

// setCache caches the given question and corresponding answer
func setCache(question, answer string) {
	db := openCache()

	err := db.Update(func(tx *bolt.Tx) (error) {
		b, err := tx.CreateBucketIfNotExists([]byte("42"))
		check(err)

		return b.Put([]byte(question), []byte(answer))
	})

	check(err)
}

// googleRequest send the request to the specified search engine
// and returns the response content
func googleRequest(question string) ([]byte) {
	// build url
	u, err := url.Parse(SEARCH_ENGINE)
	check(err)

	// encode question
	param := url.Values{}
	param.Add(QUERY_PARAM, question)
	u.RawQuery = param.Encode()

	return get(u.String())
}

// get send a http-get request to the given url
// and returns it body/content
func get(url string) ([]byte) {
	resp, err := http.Get(url)
	check(err)

	defer resp.Body.Close()

	// retrieve content/body
	content, err := ioutil.ReadAll(resp.Body)
	check(err)

	return content
}

func parse(content []byte, fn parseFunc) ([]string) {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	check(err)

	var matches []string

	fn(document, &matches)

	return matches
}

// parseGoogle parses the returned content from the request and
// returns all matches found
func parseGoogle(content []byte) ([]string) {
	return parse(content, func(document *goquery.Document, matches *[]string) {
		// find all links containing "https://stackoverflow.com" in attribute href
		// and append them to matches array
		document.Find("a").Each(func(i int, s *goquery.Selection) {
			if href, ok := s.Attr("href"); ok && strings.Contains(href, "/url?q=https://stackoverflow.com") {
				*matches = append(*matches, strings.Split(strings.Trim(href, "/url?q="), "&")[:1][0])
			}
		})
	})
}

// parseStackOverflow parses the returned content from request and
// prints the best voted answer
func parseStackOverflow(content []byte) ([]string) {
	return parse(content, func(document *goquery.Document, matches *[]string) {
		// find the accepted answer
		document.Find(".accepted-answer").Each(func(i int, s *goquery.Selection) {
			*matches = append(*matches, s.Find(".post-text").Text())

			return
		})
	})
}
