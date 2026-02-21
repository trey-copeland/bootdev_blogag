package commands

import (
	"context"
	"fmt"
)

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Queries.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error querying feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		if feed.Url.Valid {
			fmt.Printf("* %s (%s) - by %s\n", feed.Name, feed.Url.String, feed.UserName)
		} else {
			fmt.Printf("* %s - by %s\n", feed.Name, feed.UserName)
		}
	}

	return nil
}
