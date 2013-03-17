Usage:

	rss http://site.com/path/to/feed
	rss -f file

The first form retrieves the summary for a single feed,
while the second form retrieves the summaries for all feed URLs listed in a line-delimted file.
The t flag can be used to filter out old items in each feed:

	rss -f file -t '2013-03-17 14:03:49 -0400 EDT'

I've also provided a script, `feed`, which makes some assumptions about where feedlist and time
info are stored so that I don't have to use those flags at all. I'm not a fan of dotfiles, so I chose
to not force my config style on others who wanted to use `rss`. As `feed` shows, it's easy to use
the -f and -t flags to store the info wherever one pleases.

Example output:

	The Daily WTF
	http://thedailywtf.com/
		Error'd: Warning - Upgrade Have a Risk, Working Wriness
			http://thedailywtf.com/Articles/Warning--Upgrade-Have-a-Risk,-Working-Wriness.aspx
		Less is More
			http://thedailywtf.com/Articles/Less-is-More.aspx

	Homages, Ripoffs, and Coincidences
	http://www.blogger.com/feeds/7279306243558348368/posts/default?start-index=26&max-results=25
		
			http://shotcontext.blogspot.com/2013/03/i-just-discovered-mike-geberts.html
