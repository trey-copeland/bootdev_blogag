package commands

import (
	"fmt"

	"github.com/trey.copeland/bootdev_blogag/internal/config"
	"github.com/trey.copeland/bootdev_blogag/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type State struct {
	Config  *config.Config
	Queries *database.Queries
}

type HandlerFunc func(*State, Command) error

type Commands struct {
	cmdMap map[string]HandlerFunc
}

func New() *Commands {
	return &Commands{
		cmdMap: make(map[string]HandlerFunc),
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, exists := c.cmdMap[cmd.Name]
	if !exists {
		return fmt.Errorf("Command not registered: %s", cmd.Name)
	}
	return f(s, cmd)
}

func (c *Commands) Register(name string, f HandlerFunc) {
	c.cmdMap[name] = f
}
