package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	SMTPHost       string `json:"SMTPHost"`
	SMTPPort       int    `json:"SMTPPort"`
	SenderEmail    string `json:"SenderEmail"`
	SenderPassword string `json:"SenderPassword"`
	DatabaseHost   string `json:"DatabaseHost"`
	DatabasePort   int    `json:"DatabasePort"`
	DatabaseUser   string `json:"DatabaseUser"`
	DatabaseName   string `json:"DatabaseName"`
	AdminEmail     string `json:"AdminEmail"`
	AdminPassword  string `json:"AdminPassword"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
