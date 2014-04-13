package repository

import (
	"encoding/json"
	"fmt"
	"io"
	conf "mback/config"
	"os"
	"path/filepath"
	"strings"
)

func (r *Repository) addRecord(rec *Record) {
	r.Records = append(r.Records, *rec)
}

func (r *Repository) freeId() int {
	max := 0
	r.each(func(rec Record) {
		if rec.Id > max {
			max = rec.Id
		}
	})
	return max + 1
}

func (r *Repository) getDir() string {
	return filepath.Join(conf.GetRepositoryRoot(), r.Name)
}

func (r *Repository) getConfigFile() string {
	return r.getRepoFile(conf.CONF_FILE_NAME)
}

func (r *Repository) getRepoFile(fileName string) string {
	return filepath.Join(r.getDir(), fileName)
}

func buildRepoFileName(path string, id int) string {
	i := strings.LastIndex(path, "/")
	return fmt.Sprintf(REPO_FILE_FORMAT, id, path[i+1:])
}

func (r *Repository) containsPath(path string) bool {
	result := false
	r.each(func(rec Record) {
		if rec.GetRealPath() == path {
			result = true
		}
	})
	return result
}

func (r *Repository) decode(reader io.Reader) (err error) {
	err = json.NewDecoder(reader).Decode(r)
	return
}

func (r *Repository) encode(f *os.File) (err error) {
	data, err := json.MarshalIndent(r, "", "  ")

	if err != nil {
		return
	}

	_, err = f.Write(data)

	return
}

func (r *Repository) each(f func(Record)) {
	for _, item := range r.Records {
		f(item)
	}
}

func (r *Repository) eachPos(f func(Record, int)) {
	for i, item := range r.Records {
		f(item, i)
	}
}
