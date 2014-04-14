package main

import (
	"mback/log"
	repo "mback/repository"
	"strconv"
)

func handleErr(msg string) {
	if err == nil {
		return
	}
	log.Fatal(msg, err)
}

func initRepo(name string) (r *repo.Repository) {
	r, err = repo.Open(name)

	handleErr("can't open repository: %v")

	return
}

func parseIds(ids_str []string) (res []int, err error) {
	ids_map := make(map[int]bool, len(ids_str))

	// convert to int and skip duplicates
	for _, id_str := range ids_str {
		var id int
		id, err = strconv.Atoi(id_str)
		if err != nil {
			return
		}

		if _, ok := ids_map[id]; !ok {
			ids_map[id] = true
		}
	}

	res = make([]int, len(ids_map))
	pos := 0
	for id, _ := range ids_map {
		res[pos] = id
	}

	return
}
