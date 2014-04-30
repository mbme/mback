package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"mback/log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	BaseDir  string `json:"base_dir"`
	LogLevel uint8  `json:"log_level"`
	User     string `json:"user"`
}

var Conf *Config

func LoadConfig() {
	configPath := filepath.Join(getHomeDir(), ".config", "mback")

	var err error
	Conf, err = readConfig(configPath)

	if err != nil {
		log.Fatal("can't load config: %v", err)
	}
}

func getCurrentUser() *user.User {
	user, err := user.Current()

	if err != nil {
		log.Fatal("can't get current user: %v", err)
	}

	return user
}

func getHomeDir() string {
	user := getCurrentUser()

	return user.HomeDir
}

func (c *Config) decode(reader io.Reader) (err error) {
	err = json.NewDecoder(reader).Decode(c)
	return
}

func readConfig(filePath string) (c *Config, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer file.Close()

	c = &Config{}

	err = c.decode(file)
	if err != nil {
		return
	}

	if c.User == "" {
		c.User = getCurrentUser().Username
	}

	return
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
