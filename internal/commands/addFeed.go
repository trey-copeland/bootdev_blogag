package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trey.copeland/bootdev_blogag/internal/database"
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

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	_, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %w", err)
	}

	user, err := s.Queries.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %w", err)
	}

	now := time.Now()
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       sql.NullString{String: feedURL, Valid: true},
		UserID:    user.ID,
	}

	createdFeed, err := s.Queries.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error creating feed: %w", err)
	}

	if createdFeed.Url.Valid {
		fmt.Printf("Feed created: %s (%s)\n", createdFeed.Name, createdFeed.Url.String)
	} else {
		fmt.Printf("Feed created: %s\n", createdFeed.Name)
	}

	return nil
}
