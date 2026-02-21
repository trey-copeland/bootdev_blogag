package commands

import (
	"context"
	"fmt"

	"github.com/trey.copeland/bootdev_blogag/internal/rss"
)

func handlerAgg(s *State, cmd Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("handlerAgg fetch feed %q: %w", feedURL, err)
	}
	fmt.Println(feed)

	return nil
}
