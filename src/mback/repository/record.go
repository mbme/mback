package repository

import (
	"mback/config"
	conf "mback/config"
	"os"
	"path/filepath"
	"strings"
)

type Record struct {
	Id         int    `json:"id"`
	Path       string `json:"path"`
	User       string `json:"user"`
	repository *Repository
}

func (r *Record) GetRealPath() string {
	if strings.HasPrefix(r.Path, "~/") {
		return filepath.Join("/home", config.USER, r.Path[2:])
	} else {
		return r.Path
	}
}

func (r *Record) SetRealPath(path string) {
	r.Path = simplifyPath(path)
}

func (r *Record) GetRepoFileName() string {
	return buildRepoFileName(r.Path, r.Id)
}

func (r *Record) GetRepoPath() string {
	return r.repository.getRepoFilePath(r.GetRepoFileName())
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

	repoFileStat, err := os.Stat(r.GetRepoPath())

	if err != nil {
		return false
	}

	return os.SameFile(stat, repoFileStat)
}

func simplifyPath(file_path string) string {
	home_dir := filepath.Join("/home", conf.USER)

	if !strings.HasPrefix(file_path, home_dir) {
		return file_path
	}

	return strings.Replace(file_path, home_dir, "~", 1)
}

// func (r *Record) String() string {
// 	return r.GetRepoFileName()
// }
