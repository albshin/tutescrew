package config

import (
	"encoding/json"
	"os"
)

// Configuration holds config values
type Configuration struct {
	Token   string    `json:"Token"`
	Prefix  string    `json:"Prefix"`
	CASInfo CASConfig `json:"CAS"`
}

// CASConfig holds config values related to CAS
type CASConfig struct {
	CASAuthURL     string `json:"CASAuthURL"`
	CASRedirectURL string `json:"CASRedirectURL"`
}

// LoadConfig loads configuration files from JSON into a struct
func LoadConfig(c *Configuration) error {
	// Read config file
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&c)
	if err != nil {
		return err
	}

	return nil
}
