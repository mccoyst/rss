package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zippoxer/RSS-Go"
)

var feeds = flag.String("f", "", "file containing a list of feeds")

func main() {
	flag.Parse()

	if flag.NArg() == 0 && *feeds == "" {
		os.Stderr.WriteString("I need the feed URL.\n")
		os.Exit(1)
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
		os.Stdout.WriteString("\n")
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

	os.Stdout.WriteString(feed.Title + "\n")
	os.Stdout.WriteString(feed.Link + "\n")
	for _, i := range feed.Items {
		os.Stdout.WriteString("\t" + strings.Replace(i.Title, "\n", " ", -1) + "\n")
		os.Stdout.WriteString("\t\t" + i.Link + "\n")
	}
}

func maybeDie(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
