package commands

import "fmt"

func handlerHelp(c *Commands, s *State, cmd Command) error {
	for _, meta := range c.Meta() {
		fmt.Printf("%s\n  usage: %s\n  %s\n", meta.Name, meta.Usage, meta.Description)
	}
	return nil
}
