package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/zippoxer/RSS-Go"
)

var feeds = flag.String("f", "", "file containing a list of feeds")
var lastRun = flag.String("t", "", "omit items preceding this date/time, in format \"2006-01-02 15:04:05 -0700 MST\"")

var begin = time.Time{}

func main() {
	flag.Parse()

	if flag.NArg() == 0 && *feeds == "" {
		os.Stderr.WriteString("I need the feed URL.\n")
		os.Exit(1)
	}

	if *lastRun != "" {
		var err error
		begin, err = time.Parse("2006-01-02 15:04:05 -0700 MST", *lastRun)
		maybeDie(err)
	}

	lines := [][]byte{[]byte(flag.Arg(0))}

	if *feeds != "" {
		f, err := os.Open(*feeds)
		maybeDie(err)

		data, err := ioutil.ReadAll(f)
		f.Close()
		maybeDie(err)

		lines = bytes.Split(data, []byte{'\n'})
	}

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		showFeed(string(line))
	}
}

func showFeed(s string) {
	url, err := url.Parse(s)
	maybeDie(err)

	resp, err := http.Get(url.String())
	maybeDie(err)
	defer resp.Body.Close()

	feed, err := rss.Get(resp.Body)
	maybeDie(err)

	newItems := make([]*rss.Item, 0, len(feed.Items))
	for _, i := range feed.Items {
		if i.When.After(begin) {
			newItems = append(newItems, i)
		}
	}
	feed.Items = newItems

	if len(feed.Items) == 0 {
		return
	}

	os.Stdout.WriteString(feed.Title + "\n")
	os.Stdout.WriteString(feed.Link + "\n")
	for _, i := range feed.Items {
		os.Stdout.WriteString("\t" + strings.Replace(i.Title, "\n", " ", -1) + "\n")
		os.Stdout.WriteString("\t\t" + i.Link + "\n")
	}
	os.Stdout.WriteString("\n")
}

func maybeDie(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
