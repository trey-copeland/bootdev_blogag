package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/trey.copeland/bootdev_blogag/internal/commands"
	"github.com/trey.copeland/bootdev_blogag/internal/config"
	"github.com/trey.copeland/bootdev_blogag/internal/database"
)

type configAdapter struct {
	cfg *config.Config
}

func (a configAdapter) SetUser(currentUserName string) error {
	return a.cfg.SetUser(currentUserName)
}

func (a configAdapter) CurrentUserName() string {
	return a.cfg.CurrentUserName
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: command required")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: args[1],
		Args: args[2:],
	}

	if err := run(cmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd commands.Command) error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("Read config: %w", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		return fmt.Errorf("Database access error: %w", err)
	}
	dbQueries := database.New(db)

	appState := commands.State{
		Config:  configAdapter{cfg: &cfg},
		Queries: dbQueries,
	}

	appCmds := commands.New()
	commands.RegisterDefault(appCmds)

	if err := appCmds.Run(&appState, cmd); err != nil {
		return err
	}

	return nil
}
