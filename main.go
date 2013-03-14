package main

import (
	"flag"
	"os"
	"net/http"
	"net/url"

	"github.com/zippoxer/RSS-Go"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		os.Stderr.WriteString("I need the URL.\n")
		os.Exit(1)
	}

	url, err := url.Parse(flag.Arg(0))
	maybeDie(err)

	resp, err := http.Get(url.String())
	maybeDie(err)
	defer resp.Body.Close()

	feed, err := rss.Get(resp.Body)
	maybeDie(err)

	os.Stdout.WriteString(feed.Title+"\n")
	os.Stdout.WriteString(feed.Link+"\n")
	for _, i := range feed.Items {
		os.Stdout.WriteString("\t"+i.Title+"\n")
		os.Stdout.WriteString("\t\t"+i.Link+"\n")
	}
}

func maybeDie(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
