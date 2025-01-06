package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"io"
	"net/http"
	"time"
)

// iterate over the items and print them to the console
func printFeed(rssFeed *RSSFeed) {
	fmt.Printf("Feed: %v\n", rssFeed.Channel.Title)
	fmt.Printf("Link: %v\n", rssFeed.Channel.Link)
	fmt.Printf("Description: %v\n", rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("%v.\n", i)
		fmt.Printf("Title: %v\n", item.Title)
		fmt.Printf("Link: %v\n", item.Link)
		fmt.Printf("Description: %v\n", item.Description)
		fmt.Printf("Publication Date: %v\n", item.PubDate)
		fmt.Println()
	}
	fmt.Println("--------------------------------------------------")
}

func scrapeFeeds(s *state) error {
	var (
		feed    database.Feed
		rssFeed *RSSFeed
		err     error
	)

	// get the next feed to fetch from our database
	feed, err = s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	// fetch the feed from the internet
	rssFeed, err = fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)

	printFeed(rssFeed)
	return err
}

func parseXML(resp *http.Response, feed *RSSFeed) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[!] Status error: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, feed)
	if err != nil {
		return err
	}

	// unescape HTML characters in the struct fields in both Channel and Items.
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Link = html.UnescapeString(feed.Channel.Link)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Link = html.UnescapeString(feed.Channel.Item[i].Link)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].PubDate = html.UnescapeString(feed.Channel.Item[i].PubDate)
	}
	return err
}

func fetchFeed(_ctx context.Context, feedURL string) (*RSSFeed, error) {
	var (
		httpClient *http.Client
		ctx        context.Context
		cancel     context.CancelFunc
		req        *http.Request
		resp       *http.Response
		err        error
		feed       *RSSFeed = &RSSFeed{}
	)
	ctx, cancel = context.WithTimeout(_ctx, 3*time.Second) // request timeout after 3s
	defer cancel()

	httpClient = &http.Client{}
	req, err = http.NewRequestWithContext(ctx, "GET", feedURL, nil) // create a new request with the timeout context
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	c := make(chan error, 1)
	go func() {
		var _err error
		resp, _err = httpClient.Do(req) // send the request
		if _err != nil {
			c <- _err
			return
		}

		defer resp.Body.Close()
		err = parseXML(resp, feed)
		c <- _err
	}()

	select {
	case <-ctx.Done(): // case timeout
		go func() { <-c }()
		return feed, ctx.Err()
	case err = <-c: // case finished parsing before timeout
		return feed, err
	}
}

func handlerAgg(s *state, cmd command) error {
	var (
		err                 error
		timeBetweenRequests time.Duration
	)

	if len(cmd.args) < 1 {
		return fmt.Errorf("[!] Usage: agg <time_between_reqs>\n")
	}

	// run the scrapeFeeds function every timeBetweenRequests seconds
	timeBetweenRequests, err = time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return err
}
