package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mmandelstrom/gatorcli/internal/database"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Cfg *Config
	Db  *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CmdNames map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if _, exists := c.CmdNames[name]; exists {
		fmt.Printf("command: %s is already registered\n", name)
		return
	}
	c.CmdNames[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	existingCmd, ok := c.CmdNames[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	if err := existingCmd(s, cmd); err != nil {
		return err
	}
	return nil
}

func (cfg Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	if err := writeToConfigFile(cfg); err != nil {
		return fmt.Errorf("unable to write to config, error: %s", err)
	}
	return nil
}

func ReadConfig() (Config, error) {
	cfg := Config{}
	path, err := getConfigFilePath()
	if err != nil {
		return cfg, fmt.Errorf("unable to get path, error: %s", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("error: %s \nunable to read from %s", err, path)
	}

	if err := json.Unmarshal(content, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory, error: %s", err)
	}
	fullPath := homeDir + "/" + configFileName
	return fullPath, nil
}

func writeToConfigFile(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to get path, error: %s", err)
	}

	content, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("unable to marshal struct to json, error: %s", err)
	}

	if err := os.WriteFile(path, content, 0666); err != nil {
		return fmt.Errorf("unable to write to %s\nerror: %s", path, err)
	}
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("user does not exist\n")
		os.Exit(1)
	}
	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("unable to set user with state pointer")
	}
	fmt.Println("User has been updated!")
	return nil
}

func HandlerDelUsers(s *State, cmd Command) error {
	err := s.Db.DelUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Users database has been reset")
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	usrSlice, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, usr := range usrSlice {
		if usr.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", usr.Name)
		} else {
			fmt.Printf("* %s\n", usr.Name)
		}

	}
	return nil
}

func RegisterHandler(s *State, cmd Command) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("no name was passed as argument")
	}
	name := os.Args[2]
	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		fmt.Printf("User: %s already exists\n", name)
		os.Exit(1)
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name}

	s.Db.CreateUser(context.Background(), userParams)
	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("user: %s was added\n", name)

	return nil
}
