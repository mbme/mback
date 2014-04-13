package repository

import (
	"errors"
	"fmt"
	conf "mback/config"
	"mback/log"
	"mback/utils"
	"os"
	path "path/filepath"
)

const (
	DIR_PERM         = 0755
	FILE_PERM        = 0644
	REPO_FILE_FORMAT = "%d %v"
)

type Repository struct {
	Name    string `json:"-"`
	Records []Record
}

func New(name string) (r Repository) {
	r = Repository{Name: name}
	return
}

func (r *Repository) GetPath() string {
	return path.Join(conf.GetRepositoryRoot(), r.Name)
}

func (r *Repository) Exists() bool {
	_, err := os.Stat(r.GetPath())
	return !os.IsNotExist(err)
}

func (r *Repository) Create() (err error) {
	err = os.Mkdir(r.GetPath(), DIR_PERM)
	if err != nil {
		return
	}
	log.Debug("Created directory for repository %v", r.Name)

	r.Records = make([]Record, 0)

	err = r.WriteConfig()
	if err != nil {
		return
	}

	log.Debug("Created empty config file for repository %v", r.Name)
	return
}

func (r *Repository) Delete() (err error) {
	return os.RemoveAll(r.GetPath())
}

func (r *Repository) AddFiles(files ...string) (err error) {
	err = r.ReadConfig()
	if err != nil {
		return
	}

	for _, file_path := range files {
		err = r.addFile(file_path)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) RemoveFiles(ids ...int) (err error) {
	for _, id := range ids {
		err = r.removeRecord(id)
		if err != nil {
			return
		}
	}
	return
}

func (r *Repository) ListIds() (ids []int) {
	ids = make([]int, len(r.Records))
	for i, rec := range r.Records {
		ids[i] = rec.Id
	}

	return
}

func (r *Repository) ReadConfig() (err error) {
	if r.Records != nil {
		log.Debug("Config was already read")
		return
	}

	file_path := r.getConfigFile()

	file, err := os.Open(file_path)
	if err != nil {
		return
	}

	defer file.Close()

	err = r.decode(file)
	if err != nil {
		return
	}

	return
}

func (r *Repository) WriteConfig() (err error) {
	if r.Records == nil {
		err = errors.New("Repository config is nil")
		return
	}

	file_path := r.getConfigFile()

	file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, FILE_PERM)
	if err != nil {
		return
	}

	defer file.Close()

	err = r.encode(file)
	return
}

func (r *Repository) GetRecord(id int) (record *Record, pos int, err error) {
	pos = -1
	r.eachPos(func(rec Record, i int) {
		if rec.Id == id {
			record = &rec
			pos = i
		}
	})

	if record == nil {
		err = fmt.Errorf("Can't find record %d", id)
	}
	return
}

func (r *Repository) InstallFile(id int) (err error) {
	rec, _, err := r.GetRecord(id)
	if err != nil {
		return
	}

	file_path := rec.GetRealPath()

	// first backup file if it exists
	if _, err := os.Stat(file_path); !os.IsNotExist(err) {
		err = utils.Backup(file_path)
		if err != nil {
			return err
		}
	}

	repo_file_path := r.getRepoFile(rec.GetRepoFileName())

	return os.Symlink(repo_file_path, file_path)
}

func (r *Repository) UninstallFile(id int) (err error) {
	rec, _, err := r.GetRecord(id)
	if err != nil {
		return
	}

	// check if installed (Symlink to repo file)
	// remove symlink
	// restore backup if exists

	if !rec.IsInstalled(r) {
		return fmt.Errorf("File %d is not installed", id)
	}

	file_path := rec.GetRealPath()

	err = os.Remove(file_path)
	if err != nil {
		return
	}

	return utils.RestoreBackup(file_path)
}
