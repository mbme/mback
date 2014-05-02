package utils

import (
	"mback/log"
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

	config := NewFile(configPath)

	data, err := Fs.Read(config)
	if err != nil {
		log.Fatal("can't load config: %v", err)
	}

	Conf = &Config{}
	err = Decode(data, Conf)
	if err != nil {
		log.Fatal("can't load config: %v", err)
	}

	log.LogLevel = Conf.LogLevel
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
