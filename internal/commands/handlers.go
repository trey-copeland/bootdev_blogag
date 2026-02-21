package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trey.copeland/bootdev_blogag/internal/database"
)

func RegisterDefault(c *Commands) {
	c.Register("login", handlerLogin)
	c.Register("register", handlerRegister)
	c.Register("reset", handlerReset)
	c.Register("users", handlerUsers)
}

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
		return fmt.Errorf("User doesn't exist")
	}
	s.Config.SetUser(cmd.Args[0])

	fmt.Println("User has been set")
	return nil
}

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

	s.Config.SetUser(name)

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
	fmt.Printf("User created: %v", dbUser)

	return nil
}

func handlerReset(s *State, cmd Command) error {
	if err := s.Queries.ClearUsers(context.Background()); err != nil {
		return fmt.Errorf("Error clearing users from database: %w", err)
	}
	fmt.Println("Users cleared from database")

	return nil
}

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
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}
