package main

import (
	"fmt"
	"mback/log"
	repo "mback/repository"
	"mback/utils"
	"strconv"
)

func handleErr(msg string) {
	if err == nil {
		return
	}
	log.Fatal(msg, err)
}

func createRepo(name string) {
	r := repo.New(name)
	if r.Exists() {
		log.Fatal("Repository '%v' already exists!", r.Name)
	}

	err = r.Create()
	handleErr("Can't create repository: %v")
	log.Info("Created repository '%v'", r.Name)
}

func initRepo(name string) repo.Repository {
	r := repo.New(name)

	exists := r.Exists()
	if !exists {
		log.Fatal("does not exists")
	}

	err = r.ReadConfig()
	handleErr("can't read config file %v:")

	return r
}

func deleteRepo(name string) {
	log.Info("Deleting repository %v", name)

	r := initRepo(name)

	msg := fmt.Sprintf("Repository %s contains %d items.\nDo you really want to delete it?", r.Name, len(r.Records))
	if !utils.Confirmation(msg) {
		log.Fatal("Cancelled")
	}

	err = r.Delete()
	handleErr("Can't delete repository " + name + ": %v")
}

func deleteFiles(r *repo.Repository, ids []int) {

	records := make([]*repo.Record, len(ids))
	for i, id := range ids {
		var rec *repo.Record
		rec, _, err = r.GetRecord(id)
		handleErr("Can't find files: %v")

		records[i] = rec
	}

	log.Info("Removing files:")
	for _, rec := range records {
		log.Info(rec.GetRealPath())
	}

	if !utils.Confirmation("Proceed?") {
		log.Fatal("Cancelled")
	}

	err = r.RemoveFiles(ids...)
	handleErr("Can't delete files: %v")

	err = r.WriteConfig()
	handleErr("Can't write config: %v")
}

func parseIds(ids_str []string) (ids []int, err error) {
	ids = make([]int, len(ids_str))
	for pos, str := range ids_str {
		ids[pos], err = strconv.Atoi(str)
		if err != nil {
			return
		}
	}

	ids = skipDuplicates(ids)
	return
}

func skipDuplicates(ids []int) []int {
	ids_map := make(map[int]bool, len(ids))

	for _, id := range ids {
		if _, ok := ids_map[id]; !ok {
			ids_map[id] = true
		}
	}

	res := make([]int, len(ids_map))
	pos := 0
	for id, _ := range ids_map {
		res[pos] = id
	}

	return res
}
