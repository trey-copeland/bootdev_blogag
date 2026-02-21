package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/trey.copeland/bootdev_blogag/internal/rss"
)

func handlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return errors.New("Incorrect number of arguments. 2 required")
	}

	currentUser := s.Config.CurrentUserName()
	if currentUser == "" {
		return errors.New("No user logged in")
	}

	// feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %w", err)
	}

	_, err = s.Queries.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %w", err)
	}

	fmt.Println(feed)

	return nil
}
