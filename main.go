package main

import (
	"fmt"
	"os"

	"github.com/trey.copeland/bootdev_blogag/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("Read config: %w", err)
	}

	if err := cfg.SetUser("trey"); err != nil {
		return fmt.Errorf("Set user: %w", err)
	}

	cfg, err = config.Read()
	if err != nil {
		return fmt.Errorf("Read: %w", err)
	}

	fmt.Println(cfg)
	return nil
}
