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
	"regexp"
	"sync"
)

type Matcher struct {
	Pattern      string
	MatchBreaker bool
	Negate       bool
}

type Config struct {
	Feeds    []string
	Matchers []Matcher
	Feed
}

type App struct {
	cnf             Config
	includeMatchers []matcher
	shutdown        chan int
	wg              *sync.WaitGroup
	*FeedReader
}

type regexpMatcher struct {
	matchBreaker bool
	negate       bool
	*regexp.Regexp
}

func (rm regexpMatcher) MatchBreaker() bool {
	return rm.matchBreaker
}

func (rm regexpMatcher) Negate() bool {
	return rm.negate
}

func NewApp(config Config) *App {
	wg := &sync.WaitGroup{}
	shutdown := make(chan int)
	app := &App{cnf: config, shutdown: shutdown, wg: wg}

	for _, matchRe := range app.cnf.Matchers {
		// make all the Regexps case insensitive
		app.includeMatchers =
			append(app.includeMatchers, regexpMatcher{matchBreaker: matchRe.MatchBreaker, negate: matchRe.Negate,
				Regexp: regexp.MustCompile("(?i)" + matchRe.Pattern)})
	}

	app.FeedReader = newFeedReader(app.cnf, app.feedItemsFilter, feedFactory(app.cnf.Feed))

	return app
}

func (a *App) Run() *App {
	a.FeedReader.run(a.shutdown, a.wg)
	return a
}

func (a *App) Shutdown() {
	close(a.shutdown)
	a.wg.Wait()
}

func (a App) feedItemsFilter(items []*rss.Item) []*rss.Item {
	return applyFilters(items, a.includeMatchers)
}
