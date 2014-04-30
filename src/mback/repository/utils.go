package repository

import (
	"fmt"
	"mback/utils"
	"path/filepath"
	path "path/filepath"
	"strings"
)

const CONF_FILE_NAME = ".mback"

func (r *Repository) getFreeId() int {
	max := 0
	r.each(func(rec *Record) {
		if rec.Id > max {
			max = rec.Id
		}
	})
	return max + 1
}

func (r *Repository) getConfigFile() *utils.File {
	return r.getRepoFile(CONF_FILE_NAME)
}

func (r *Repository) getRepoFile(fileName string) *utils.File {
	return utils.NewFile(filepath.Join(r.GetRootPath(), fileName))
}

func (r *Repository) containsPath(path string) bool {
	result := false
	r.each(func(rec *Record) {
		if rec.GetFile().GetPath() == path {
			result = true
		}
	})
	return result
}

func (r *Repository) each(f func(*Record)) {
	for _, item := range r.Records {
		f(item)
	}
}

func (r *Repository) eachPos(f func(*Record, int)) {
	for i, item := range r.Records {
		f(item, i)
	}
}

func exists(name string) bool {
	return utils.NewFile(getRepoRootPath(name)).Exists()
}

func getRepoRootPath(name string) string {
	return path.Join(utils.Conf.BaseDir, name)
}

func buildRepoFileName(path string, id int) string {
	i := strings.LastIndex(path, "/")
	return fmt.Sprintf(REPO_FILE_FORMAT, id, path[i+1:])
}
