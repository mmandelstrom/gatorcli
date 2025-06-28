package main

import (
	"fmt"

	"github.com/mmandelstrom/gatorcli/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		print("error: %s", err)
	}

	if err := cfg.SetUser("marcus"); err != nil {
		fmt.Printf("error: %s", err)
	}

	content, err := config.ReadConfig()
	if err != nil {
		print("error: %s", err)
	}
	fmt.Printf("DbURL: %s, CurrentUsername: %s\n", content.DbURL, content.CurrentUserName)
}
