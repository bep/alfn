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
	rss "github.com/jteeuwen/go-pkg-rss"
	"sort"
)

type matcher interface {
	MatchString(s string) bool
	MatchBreaker() bool
	Negate() bool
}

func appyFilters(items []*rss.Item, matchers []matcher) []*rss.Item {
	var filtered []*rss.Item

OUTER:
	for _, item := range items {
		for _, m := range matchers {
			if !m.MatchBreaker() {
				continue
			}
			content := textContent(item)
			match := m.MatchString(content)
			if match == m.Negate() {
				continue OUTER
			}

		}

		for _, m := range matchers {
			if m.MatchBreaker() {
				continue
			}
			content := textContent(item)
			match := m.MatchString(content)
			if match != m.Negate() {
				filtered = append(filtered, item)
				continue OUTER
			}
		}
	}
	return filtered
}

func textContent(item *rss.Item) string {
	var content string
	if item.Title != "" {
		content += item.Title + " "
	}
	if item.Description != "" {
		content += item.Description + " "
	}

	if item.Content != nil {
		content += item.Content.Text
	}
	return content
}

func limit(items rssItems, n int) []*rssItem {
	sort.Sort(items)

	if n < len(items) {
		for i := n; i < len(items); i++ {
			items[i] = nil
		}
		return items[:n]
	}

	return items

}

func (slice rssItems) Len() int {
	return len(slice)
}

//TODO(bep): Add some weight?
func (slice rssItems) Less(i, j int) bool {
	p1, err1 := slice[i].ParsedPubDate()
	p2, err2 := slice[j].ParsedPubDate()
	if err1 != nil && err2 == nil {
		return true
	}
	if err2 != nil && err1 == nil {
		return false
	}
	return p1.After(p2)
}

func (slice rssItems) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
