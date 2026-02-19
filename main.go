package main

import (
	"fmt"
	"os"

	"github.com/trey.copeland/bootdev_blogag/internal/config"
)

type command struct {
	name string
	args []string
}

type state struct {
	config *config.Config
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, exist := c.cmdMap[cmd.name]
	if !exist {
		return fmt.Errorf("Command not registered: %s", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.cmdMap[name] = f
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No argument provided to login")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Too many arguments provided to login")
	}

	s.config.SetUser(cmd.args[0])

	fmt.Println("User has been set")
	return nil
}

func main() {
	args := os.Args
	if (len(args) == 2) && (args[1] == "login") {
		fmt.Fprintln(os.Stderr, "Error: username required")
		os.Exit(1)
	}
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: too few arguments provided")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := run(cmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd command) error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("Read config: %w", err)
	}

	appState := state{
		config: &cfg,
	}

	cmdMap := make(map[string]func(*state, command) error)
	appCmds := commands{
		cmdMap: cmdMap,
	}
	appCmds.register("login", handlerLogin)

	if err := appCmds.run(&appState, cmd); err != nil {
		return err
	}

	return nil
}
