package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mmandelstrom/gatorcli/internal/config"
	"github.com/mmandelstrom/gatorcli/internal/database"
)

func main() {
	s := config.State{}
	content, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	s.Cfg = &content

	db, err := sql.Open("postgres", s.Cfg.DbURL)
	if err != nil {
		fmt.Printf("unable to open db")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.Db = dbQueries

	cmds := config.Commands{CmdNames: make(map[string]func(*config.State, config.Command) error)}
	cmds.Register("login", config.HandlerLogin)
	cmds.Register("register", config.RegisterHandler)

	if len(os.Args) < 2 {
		fmt.Printf("too few arguments\n")
		os.Exit(1)
	}
	cmd := config.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Printf("Error running command: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
