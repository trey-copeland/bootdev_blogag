package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No argument provided to login")
	}
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Too many arguments provided to login")
	}

	name := cmd.Args[0]
	_, err := s.Queries.GetUser(context.Background(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("User doesn't exist")
		}
		return fmt.Errorf("Error querying user: %w", err)
	}
	if err := s.Config.SetUser(name); err != nil {
		return fmt.Errorf("Error setting current user: %w", err)
	}

	fmt.Println("User has been set")
	return nil
}
