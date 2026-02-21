package commands

import (
	"context"
	"fmt"

	"github.com/trey.copeland/bootdev_blogag/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type ConfigStore interface {
	SetUser(currentUserName string) error
	CurrentUserName() string
}

type QueryStore interface {
	GetUser(ctx context.Context, name string) (database.User, error)
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	ClearUsers(ctx context.Context) error
	GetUsers(ctx context.Context) ([]string, error)
}

type State struct {
	Config  ConfigStore
	Queries QueryStore
}

type HandlerFunc func(*State, Command) error

type CommandMeta struct {
	Name        string
	Usage       string
	Description string
}

type Commands struct {
	cmdMap   map[string]HandlerFunc
	metaList []CommandMeta
}

func New() *Commands {
	return &Commands{
		cmdMap:   make(map[string]HandlerFunc),
		metaList: []CommandMeta{},
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, exists := c.cmdMap[cmd.Name]
	if !exists {
		return fmt.Errorf("Command not registered: %s", cmd.Name)
	}
	if err := f(s, cmd); err != nil {
		return fmt.Errorf("run command %q: %w", cmd.Name, err)
	}
	return nil
}

func (c *Commands) Register(name string, f HandlerFunc) {
	c.cmdMap[name] = f
}

func (c *Commands) RegisterWithMeta(meta CommandMeta, f HandlerFunc) {
	c.cmdMap[meta.Name] = f
	c.metaList = append(c.metaList, meta)
}

func (c *Commands) Meta() []CommandMeta {
	return c.metaList
}

func RegisterDefault(c *Commands) {
	allCommands := []struct {
		meta    CommandMeta
		handler HandlerFunc
	}{
		{
			meta: CommandMeta{Name: "help", Usage: "help", Description: "Show available commands"},
			handler: func(s *State, cmd Command) error {
				return handlerHelp(c, s, cmd)
			},
		},
		{
			meta:    CommandMeta{Name: "login", Usage: "login <name>", Description: "Set current user if user exists"},
			handler: handlerLogin,
		},
		{
			meta:    CommandMeta{Name: "register", Usage: "register <name>", Description: "Create a user and set as current"},
			handler: handlerRegister,
		},
		{
			meta:    CommandMeta{Name: "reset", Usage: "reset", Description: "Delete all users"},
			handler: handlerReset,
		},
		{
			meta:    CommandMeta{Name: "users", Usage: "users", Description: "List all users"},
			handler: handlerUsers,
		},
		{
			meta:    CommandMeta{Name: "agg", Usage: "agg", Description: "Aggregate feed"},
			handler: handlerAgg,
		},
	}

	for _, command := range allCommands {
		c.RegisterWithMeta(command.meta, command.handler)
	}
}
