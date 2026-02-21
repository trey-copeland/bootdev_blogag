package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trey.copeland/bootdev_blogag/internal/database"
)

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No argument provided to register")
	}
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Too many arguments provided to register")
	}

	name := cmd.Args[0]
	_, err := s.Queries.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("User already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Error checking existing user: %w", err)
	}

	now := time.Now()
	userParams := database.CreateUserParams{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		ID:        uuid.New(),
	}
	dbUser, err := s.Queries.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}
	if err := s.Config.SetUser(name); err != nil {
		return fmt.Errorf("Error setting current user: %w", err)
	}
	fmt.Printf("User created: %v", dbUser)

	return nil
}
