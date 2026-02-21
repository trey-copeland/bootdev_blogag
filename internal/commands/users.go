package commands

import (
	"context"
	"fmt"
)

func handlerUsers(s *State, cmd Command) error {
	users, err := s.Queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error querying users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	for _, user := range users {
		if user == s.Config.CurrentUserName() {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}
