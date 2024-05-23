package main

import (
	"fmt"
	"testing"
)

func TestScrapeURL(t *testing.T) {
	const url = "https://blog.boot.dev/index.xml"
	rss, err := fetchData(url)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	fmt.Println(rss.Channel.Title)
	fmt.Println(rss.Channel.Link.Href)
	if rss.Channel.Title != "Boot.dev Blog" {
		t.FailNow()
	}
	if rss.Channel.Link.Href != url {
		t.FailNow()
	}
	if len(rss.Channel.Item) == 0 {
		t.FailNow()
	}
}
