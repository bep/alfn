// Copyright © 2015 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package lib

import (
	"bytes"
	rss "github.com/jteeuwen/go-pkg-rss"
	"html/template"
	"net/url"
	"strings"
	"time"
)

var templateFuncs = template.FuncMap{
	"shtml": func(text string) template.HTML { return template.HTML(text) },
}

var rssTemplate *template.Template

func init() {
	rssTemplate = template.New("")
	rssTemplate.Funcs(templateFuncs)
	// The template below is kindly borrowed and adapted from https://github.com/spf13/hugo/blob/master/tpl/template_embedded.go
	_, err := rssTemplate.Parse(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ with $.Title }}{{.}}{{ end }}</title>
    <link>{{ .Link }}</link>
    {{ with $.Description }}<description>{{.}}</description>{{ end }}
    <generator>alfn</generator>{{ with $.LanguageCode }}
    <language>{{.}}</language>{{end}}
    {{ with $.Author.Name }}<managingEditor>{{.}}{{ with $.Author.Email }} ({{.}}){{end}}</managingEditor>{{end}}
    {{ with $.Author.Name }}<webMaster>{{.}}{{ with $.Author.Email }} ({{.}}){{end}}</webMaster>{{end}}{{ with $.Copyright }}
    <copyright>{{.}}</copyright>{{end}}{{ if not .PubDate.IsZero }}
    <lastBuildDate>{{ .PubDate.Format "Mon, 02 Jan 2006 15:04:05 -0700" | shtml }}</lastBuildDate>{{ end }}
    <atom:link href="{{.Link}}" rel="self" type="application/rss+xml" />
    {{ range .Items }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ (index .Links 0).Href }}</link>
      <pubDate>{{ .ParsedPubDate.Format "Mon, 02 Jan 2006 15:04:05 -0700" | shtml }}</pubDate>
      <author>{{ .AuthorFormatted }}</author>
      <guid>{{ .Key }}</guid>
      <description>{{ "<![CDATA[" | shtml  }}{{ with .Content }}{{ .Text | shtml }}{{ else }}{{ .Description | shtml }}{{ end }}{{ "]]>" | shtml }}</description>
    </item>
    {{ end }}
  </channel>
</rss>`)

	if err != nil {
		panic(err)
	}
}

// Feed describes the published feed.
type Feed struct {
	Title        string
	Description  string
	Link         string
	Author       Author
	LanguageCode string
	Copyright    string
	MaxItems     int
}

// Author describes the author of the published feed.
type Author struct {
	Name  string
	Email string
}

type rssItems []*rssItem

type aggregatedFeed struct {
	PubDate time.Time
	Items   rssItems
	Feed
}

type rssItem struct {
	// Embed it so we can add a method to it.
	*rss.Item
}

func (af rssItem) AuthorFormatted() string {
	url, err := url.Parse(af.Links[0].Href)

	var source string
	var formatted string

	if err == nil {
		source = url.Host
	}

	if source != "" {
		formatted = strings.TrimPrefix(source, "www.")
		formatted = strings.TrimPrefix(formatted, "WWW.")
	}

	if af.Author.Name != "" {
		if source != "" {
			formatted += ": "
		}
		formatted += af.Author.Name
	}

	return formatted
}

func feedFactory(f Feed) func(items rssItems) (string, error) {
	return func(items rssItems) (string, error) {
		var b bytes.Buffer

		af := aggregatedFeed{
			Feed:    f,
			PubDate: time.Now(),
			Items:   items}

		err := rssTemplate.Execute(&b, af)

		if err != nil {
			return "", err
		}

		return b.String(), nil
	}
}

func toRSSItems(items []*rss.Item) rssItems {
	newItems := make([]*rssItem, len(items))
	for i, item := range items {
		newItems[i] = &rssItem{item}
	}
	return newItems
}
