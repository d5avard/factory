package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
)

var configFilename = "config.ini"
var configPath = ".config/factory"

// Define struct for config
type Config struct {
	APIKey   string `toml:"api_key"`
	Database struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		Name     string `toml:"name"`
	} `toml:"database"`
}

func GetDebugVar() bool {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		return false
	}

	if debug == "true" {
		value, err := strconv.ParseBool(debug)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}

		if value {
			return value
		}
	}

	return false
}

func GetConfigFilename(debug bool) (string, error) {
	var err error

	if debug {
		var pwd string
		if pwd, err = os.Getwd(); err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
		return filepath.Join(pwd, configFilename), nil
	} else {
		var home string
		home, err = os.UserHomeDir()

		if err != nil {
			return "", errors.New("HOME environment variable not set")
		}
		return filepath.Join(home, configPath, configFilename), nil
	}
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(filename, &config)
	return config, err
}
