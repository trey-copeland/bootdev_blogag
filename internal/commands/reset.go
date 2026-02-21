package commands

import (
	"context"
	"fmt"
)

func handlerReset(s *State, cmd Command) error {
	if err := s.Queries.ClearUsers(context.Background()); err != nil {
		return fmt.Errorf("Error clearing users from database: %w", err)
	}
	fmt.Println("Users cleared from database")

	return nil
}
