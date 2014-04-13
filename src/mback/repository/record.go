package repository

import (
	"mback/config"
	"mback/utils"
	"os"
	"path/filepath"
	"strings"
)

type Record struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
	User string `json:"user"`
}

func (r *Record) GetRealPath() string {
	if strings.HasPrefix(r.Path, "~/") {
		return filepath.Join("/home", config.USER, r.Path[2:])
	} else {
		return r.Path
	}
}

func (r *Record) IsInstalled(repo *Repository) bool {
	realPath := r.GetRealPath()

	stat, err := os.Stat(realPath)

	if err != nil {
		return false
	}

	if stat.Mode()&os.ModeSymlink != 0 {
		return false
	}

	repoFileStat, err := os.Stat(repo.getRepoFile(r.GetRepoFileName()))

	if err != nil {
		return false
	}

	return os.SameFile(stat, repoFileStat)
}

func (r *Record) SetPath(path string) {
	r.Path = utils.SimplifyPath(path)
}

func (r *Record) GetRepoFileName() string {
	return buildRepoFileName(r.Path, r.Id)
}

func (r *Record) String() string {
	return r.GetRepoFileName()
}
