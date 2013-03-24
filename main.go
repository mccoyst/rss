// © 2013 Steve McCoy. Licensed under the MIT License.

package main

import (
	"bytes"
	"errors"
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

	n := 0
	feeds := []*rss.Feed{}
	errs := []error{}
	fc := make(chan *rss.Feed)
	ec := make(chan error)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		n++
		go getFeed(string(line), fc, ec)
	}

	for i := 0; i < n; i++ {
		select {
		case f := <- fc:
			feeds = append(feeds, f)
		case e := <- ec:
			errs = append(errs, e)
		}
	}

	for _, f := range feeds {
		showFeed(f)
	}

	for _, e := range errs {
		os.Stdout.WriteString("Problem: "+e.Error()+"\n")
	}
}

func getFeed(s string, fc chan *rss.Feed, ec chan error) {
	url, err := url.Parse(s)
	if err != nil {
		ec <- errors.New(s + ": " + err.Error())
		return
	}

	resp, err := http.Get(url.String())
	if err != nil {
		ec <- errors.New(s + ": " + err.Error())
		return
	}
	defer resp.Body.Close()

	feed, err := rss.Get(resp.Body)
	if err != nil {
		ec <- errors.New(s + ": " + err.Error())
		return
	}

	newItems := make([]*rss.Item, 0, len(feed.Items))
	for _, i := range feed.Items {
		if i.When.After(begin) {
			newItems = append(newItems, i)
		}
	}
	feed.Items = newItems

	fc <- feed
}

func showFeed(feed *rss.Feed) {
	if len(feed.Items) == 0 {
		return
	}

	os.Stdout.WriteString(feed.Title + "\n")
	os.Stdout.WriteString(feed.Link + "\n")
	for _, i := range feed.Items {
		title := i.Title
		if title == "" {
			m := 40
			if n := len(i.Description); n < m {
				m = n
			}
			title = i.Description[:m]
		}
		os.Stdout.WriteString("\t" + strings.Replace(title, "\n", " ", -1) + "\n")
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
