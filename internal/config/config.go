package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	if err := write(cfg); err != nil {
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

func write(cfg Config) error {
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
