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
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type FeedReader struct {
	feed
	itemFilter feedFilter
	feedWriter feedWriter
	cnf        Config
	items      rssItems
	sync.Mutex
}

type feed struct {
	f atomic.Value
}

type feedFilter func(items []*rss.Item) []*rss.Item
type feedWriter func(items rssItems) (string, error)

// GetFeed returns the value set by the most recent Store.
func (fr *FeedReader) GetFeed() string {
	return fr.feed.f.Load().(string)
}

func newFeedReader(cnf Config, ff feedFilter, fw feedWriter) *FeedReader {
	r := &FeedReader{cnf: cnf, itemFilter: ff, feedWriter: fw}
	r.feed.f.Store("")
	return r
}

func (fr *FeedReader) run() {
	for _, feed := range fr.cnf.Feeds {
		go fr.poll(feed)
	}
}

// TODO(bep)
func (fr *FeedReader) genFeed() (string, error) {
	return fr.feedWriter(fr.items)
}

func (fr *FeedReader) poll(uri string) {
	// TODO(bep) options, maybe
	feed := rss.New(240, true, chanHandler, fr.itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s\n", uri, err)
			return
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}

}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	//fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (fr *FeedReader) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	newitems = fr.itemFilter(newitems)

	fmt.Printf("%d filtered item(s) in %s\n", len(newitems), feed.Url)

	if len(newitems) == 0 {
		return
	}

	fr.Lock()
	fr.items = append(fr.items, toRSSItems(newitems)...)
	fr.items = limit(fr.items, fr.cnf.Feed.MaxItems)
	rss, err := fr.genFeed()
	fr.Unlock()

	if err != nil {
		fmt.Println("error: Failed to create feed:", err)
		return
	}

	fr.feed.f.Store(rss)
}

func printItems(items []*rss.Item) {
	for _, item := range items {
		fmt.Println("Item:", item.Title)
	}
}