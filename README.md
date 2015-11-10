# alfn 

**alfn** (Aggregated Latest Filtered News) is a great little tool to create your own RSS feed about only the stuff you´re really interested in! Just configure it with the feeds to watch and what to watch for, and just sit back and wait for the updates. Oh, and it supports hotreloading configuration changes!

[![Build Status](https://travis-ci.org/bep/alfn.svg)](https://travis-ci.org/bep/alfn)

## Install

The easiest way to install the binary is via `go get`:

``` bash
go get github.com/bep/alfn
```

Then run it with the `-h` flag to see the options available:

```
alfn -h
```

## Configure

`alfn` looks, by default, for configuration named either `config.toml`,  `config.yaml` or  `config.json` in either `$HOME/.alfn` or `/etc/alfn/`. The config file to use can also be set in the `--config` flag.

Below is an imaginary `liverpool-feed.toml` -- a feed about the Premier League football club Liverpool:

```
feeds = [
 	"http://www.telegraph.co.uk/sport/football/teams/liverpool/rss",
	"http://feeds.bbci.co.uk/sport/0/football/premier-league/rss.xml",
	"http://www.theguardian.com/football/premierleague/rss"
]


[[matchers]]
# If it doesn't contain any of these, stop looking.
# Note: The search also includes the RSS URL
pattern="football|premier-?\\s?league"
matchBreaker=true
 
[[matchers]]
# Beatles was from Liverpool, but have little to do with the football team, so skip those.
pattern="Beatles"
matchBreaker=true
negate=true

# The matcher(s) below decide what to include. Any article matching any of the 
# regular expressions will be added to the feed.
# Add more matchers sections as needed.
[[matchers]]
pattern="Liverpool|J(u|ü)rgen (Norbert)? Klopp"

[[matchers]]
pattern="Joe Gomez|Philippe Coutinho|Simon Mignolet"

[feed]
title="My Feed About Liverpool FC"
description="An RSS feed about the great football team!"
link="" # Defaults to the http://<server>:<port>
languageCode="EN-us"
copyright="2015 @ yourname"
maxItems=20
[feed.author]
name="Your Name"
email="yourname@yourname.com"
```

Then start your server pointing to the config file above:

```bash
alfn --config=/path/to/config/liverpool-feed.toml

Using config file: /path/to/config/liverpool-feed.toml

Starting server on http://127.0.0.1:1926 ...

2 filtered item(s) in http://feeds.bbci.co.uk/sport/0/football/premier-league/rss.xml
1 filtered item(s) in http://www.theguardian.com/football/premierleague/rss
17 filtered item(s) in http://www.telegraph.co.uk/sport/football/teams/liverpool/rss

```

## Credits

The best open source projects lean lazily on the great work of others:

* [go-pkg-rss](https://github.com/jteeuwen/go-pkg-rss)
* [Graceful](https://github.com/tylerb/graceful)
* [Cobra](https://github.com/spf13/cobra)
* [Viper](https://github.com/spf13/viper)



