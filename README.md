rss
===

Usage:

	rss http://site.com/path/to/feed
	rss -f file

The first form retrieves the summary for a single feed,
while the second form retrieves the summaries for all feed URLs listed in a line-delimted file.
The t flag can be used to filter out old items in each feed:

	rss -f file -t '2013-03-17 14:03:49 -0400 EDT'

I've also provided a script, `feed`, which makes some assumptions about where feedlist and time
info are stored so that I don't have to use those flags at all. I'm not a fan of dotfiles, so I chose
to not force my config style on others that wanted to use `rss`. As `feed` shows, it's easy to use
the -f and -t flags to store the info wherever one pleases.
