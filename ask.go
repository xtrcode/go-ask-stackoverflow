package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xtrcode/go-ask-stackoverflow/cache"
	"github.com/xtrcode/go-ask-stackoverflow/engine"
	"github.com/xtrcode/go-ask-stackoverflow/website"
)

var (
	c cache.Cache
	e engine.Engine
	w website.Website
)

func main() {
	if len(os.Args) == 1 {
		help()
	}

	question := strings.Join(os.Args[1:], " ")

	// search in cache
	c = &cache.Bolt{}
	if err := c.Init(); err != nil {
		log.Fatal(err)
	}

	if err := c.Open(); err != nil {
		log.Fatal(err)
	}

	if answer, err := c.Get(question); err == nil && len(answer) > 0 {
		fmt.Println(answer)
		os.Exit(0)
	}

	// search online
	e = &engine.Google{}
	if err := e.Request(question); err != nil {
		log.Fatal(err)
	}

	content, err := e.Get()
	if err != nil {
		log.Fatal(err)
	}

	w = &website.StackOverflow{}
	if err := w.Get(content); err != nil {
		log.Fatal(err)
	}

	answer, err := w.Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)

	// save
	if err = c.Set(question, answer); err != nil {
		log.Fatal(err)
	}
}

func help() {
	fmt.Println("ask is a simple terminal command that tries to find " +
		"the most popular answer to your question")
	fmt.Println("from StackOverflow.com")
	fmt.Println("\nusage: ask your question")
	os.Exit(0)
}
