package main

import (
	conf "mback/config"
	"mback/log"
	repo "mback/repository"
	"mback/utils"
)

var err error

func Status(args ...string) {
	args_len := len(args)

	switch args_len {
	case 0:
		log.Info("status:")
		log.Info("root dir: %v", conf.GetRepositoryRoot())
	case 1:
		r := repo.New(args[0])

		log.Info("status for repository %v:", r.Name)

		exists := r.Exists()
		if !exists {
			log.Fatal("does not exists")
		}

		err = r.ReadConfig()
		handleErr("can't read config file %v:")

		log.Info("items in repository: %d", len(r.Records))
		for _, rec := range r.Records {
			log.Info("%d  %v", rec.Id, rec.GetRealPath())
		}
	default:
		log.Fatal("too many arguments for status: %v", args_len)
	}
}

func Add(args ...string) {
	args_len := len(args)
	if args_len < 1 {
		log.Fatal("wrong number of arguments for add: %v", args_len)
	}

	repo_name := args[0]
	if args_len == 1 {
		log.Info("creating repository '%s'", repo_name)
		createRepo(repo_name)
		return
	}

	log.Info("adding files to repository '%s'", repo_name)

	r := initRepo(repo_name)

	var current_dir string
	current_dir, err = utils.GetWorkingDir()
	handleErr("can't retrieve current directory: %v")

	var files []string
	files, err = utils.ListFiles(current_dir, args[1:]...)
	handleErr("can't add files: %v")

	for _, file := range files {
		log.Info("* %s", file)
	}
	if !utils.Confirmation("proceed?") {
		log.Fatal("cancelled")
	}

	err = r.AddFiles(files...)
	handleErr("can't add files: %v")

	err = r.WriteConfig()
	handleErr("can't write config: %v")
}

func Remove(args ...string) {
	args_len := len(args)

	if args_len < 1 {
		log.Fatal("wrong number of arguments for remove: %v", args_len)
	}

	repo_name := args[0]

	if args_len == 1 {
		log.Info("removing repository '%s'", repo_name)
		deleteRepo(repo_name)
		return
	}

	r := initRepo(repo_name)

	log.Info("removing files from '%s'", repo_name)

	var ids []int
	ids, err = parseIds(args[1:])
	handleErr("can't parse file ids: %v")

	for _, id := range ids {
		var rec *repo.Record
		rec, _, err = r.GetRecord(id)
		handleErr("can't find record: %v")
		log.Info("* %s", rec.GetRealPath())
	}

	deleteFiles(&r, ids)
}

func Install(args ...string) {
	args_len := len(args)

	if args_len < 1 {
		log.Fatal("wrong number of arguments for install: %v", args_len)
	}

	repo_name := args[0]

	r := initRepo(repo_name)

	var ids []int
	if args_len > 1 {
		// if we specified list of file ids to install
		ids, err = parseIds(args[1:])
		handleErr("Can't parse file ids: %v")
	} else {
		// if ids were not specified then install all repo files
		ids = r.ListIds()
	}

	// install all files
	for _, id := range ids {
		err = r.InstallFile(id)
		handleErr("Can't install file " + string(id) + ": %v")
	}
}

func Uninstall(args ...string) {
	args_len := len(args)

	if args_len < 1 {
		log.Fatal("Wrong number of arguments for uninstall: %v", args_len)
	}

	repo_name := args[0]

	r := initRepo(repo_name)

	var ids []int
	if args_len > 1 {
		// if we specified list of file ids to uninstall
		ids, err = parseIds(args[1:])
		handleErr("Can't parse file ids: %v")
	} else {
		// if ids were not specified then uninstall all repo files
		ids = r.ListIds()
	}

	for _, id := range ids {
		err = r.UninstallFile(id)
		if err != nil {
			log.Warn("Can't uninstall file %d: %v", id, err)
		}
	}
}
