package main

import (
	"mback/log"
	"mback/utils"
	"os"
)

var commands = map[string]func(...string){
	"status":    Status,
	"add":       Add,
	"remove":    Remove,
	"install":   Install,
	"uninstall": Uninstall,
}

func main() {
	utils.LoadConfig()
	utils.Fs = &utils.RealFS{}
	log.LogLevel = utils.Conf.LogLevel

	log.Debug("Args: %v", os.Args)

	args_len := len(os.Args)
	if args_len < 2 {
		log.Fatal("Should be at least 1 argument")
	}

	command := os.Args[1]

	action, ok := commands[command]
	if !ok {
		log.Fatal("Unsupported action: %v", command)
	}

	//applying action
	action(os.Args[2:]...)
}
