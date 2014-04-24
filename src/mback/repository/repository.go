package repository

import (
	"fmt"
	"mback/log"
	"mback/utils"
)

const (
	FILE_PERM        = 0644
	REPO_FILE_FORMAT = "%d %v"
)

type Repository struct {
	Name    string `json:"-"`
	Records []*Record
}

func Create(name string) (r Repository, err error) {
	if exists(name) {
		return r, fmt.Errorf("repository '%s' already exists", name)
	}

	r = Repository{Name: name}

	err = utils.NewFile(r.GetRootPath()).Mkdir()
	if err != nil {
		return
	}
	log.Debug("created directory for repository %v", r.Name)

	r.Records = make([]*Record, 0)

	err = r.writeConfig()
	if err != nil {
		return
	}

	log.Debug("created empty config file for repository %v", r.Name)
	return
}

func Open(name string) (r *Repository, err error) {
	r = &Repository{Name: name}

	if !exists(name) {
		return r, fmt.Errorf("repository '%s' doesn't exist", name)
	}

	err = r.readConfig()

	if err == nil {
		for _, rec := range r.Records {
			rec.repository = r
		}
	}

	return
}

func (r *Repository) Delete() (err error) {
	return utils.NewFile(r.GetRootPath()).RemoveAll()
}

func (r *Repository) GetRootPath() string {
	return getRepoRootPath(r.Name)
}

func (r *Repository) ListIds() (ids []int) {
	ids = make([]int, len(r.Records))
	for i, rec := range r.Records {
		ids[i] = rec.Id
	}

	return
}

func (r *Repository) GetRecord(id int) (record *Record, pos int, err error) {
	pos = -1
	r.eachPos(func(rec *Record, i int) {
		if rec.Id == id {
			record = rec
			pos = i
		}
	})

	if record == nil {
		err = fmt.Errorf("can't find record %d", id)
	}
	return
}

func (r *Repository) InstallFile(id int) (err error) {
	rec, _, err := r.GetRecord(id)
	if err != nil {
		return
	}

	file := rec.GetFile()

	// first backup file if it exists
	if file.Exists() {
		err = utils.Backup(file.GetPath())
		if err != nil {
			return err
		}
	}

	repoFile := rec.GetRepoFile()

	return repoFile.Symlink(file.GetPath())
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
		return fmt.Errorf("file %d is not installed", id)
	}

	file := rec.GetFile()
	err = file.Remove()
	if err != nil {
		return
	}

	return utils.RestoreBackup(file.GetPath())
}

func (x *Repository) addRecord(file_path string) (err error) {
	if x.containsPath(file_path) {
		log.Info("'%v' already contains file %v, skipping", x.Name, file_path)
		return
	}

	// get file access rights and owner
	file := utils.NewFile(file_path)

	newId := x.getFreeId()
	newName := buildRepoFileName(file_path, newId)

	var fileInfo *utils.FileInfo
	fileInfo, err = file.GetInfo()

	// TODO add groups processing
	record := new(Record)
	record.Id = newId
	record.SetPath(file_path)
	record.User = fileInfo.Owner
	record.repository = x

	// copy file to repository
	err = file.CopyTo(x.getRepoFile(newName).GetPath())
	if err != nil {
		return
	}

	x.Records = append(x.Records, record)

	err = x.writeConfig()
	if err != nil {
		return
	}

	log.Info("added file %v as %v", file_path, newName)
	return
}

func (r *Repository) AddRecords(files ...string) (err error) {
	for _, file_path := range files {
		err = r.addRecord(file_path)
		if err != nil {
			return
		}
	}

	return
}

func (r *Repository) removeRecord(id int) (err error) {
	rec, pos, err := r.GetRecord(id)
	if err != nil {
		return
	}

	err = rec.GetFile().Remove()
	if err != nil {
		return
	}

	// delete file from config
	r.Records = append(r.Records[:pos], r.Records[pos+1:]...)

	err = r.writeConfig()
	if err != nil {
		return
	}

	log.Debug("removed record %v", rec)

	return
}

func (r *Repository) RemoveRecords(ids ...int) (err error) {
	for _, id := range ids {
		err = r.removeRecord(id)
		if err != nil {
			return
		}
	}
	return
}
