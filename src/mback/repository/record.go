package repository

import (
	"fmt"
	"mback/utils"
)

type Record struct {
	Id         int    `json:"id"`
	Path       string `json:"path"`
	User       string `json:"user"`
	repository *Repository
}

func (r *Record) GetRepoFileName() string {
	return buildRepoFileName(r.Path, r.Id)
}

func (r *Record) GetRepoFile() *utils.File {
	return r.repository.getRepoFile(r.GetRepoFileName())
}

func (r *Record) GetFile() *utils.File {
	return utils.NewFile(r.Path)
}

func (r *Record) SetPath(path string) {
	r.Path = utils.NewFile(path).SimplifyPath()
}

func (r *Record) IsInstalled() bool {
	installedFile := r.GetFile()

	if !installedFile.Exists() {
		return false
	}

	if !installedFile.IsLink() {
		return false
	}

	repoFile := r.GetRepoFile()

	if !repoFile.Exists() {
		panic(fmt.Sprintf("Repository file %v does not exists", r))
	}

	return installedFile.SameFile(repoFile)
}
