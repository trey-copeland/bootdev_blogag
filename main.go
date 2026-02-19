package main

import (
	"fmt"

	"github.com/trey.copeland/bootdev_blogag/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := cfg.SetUser("trey"); err != nil {
		fmt.Println(err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}
