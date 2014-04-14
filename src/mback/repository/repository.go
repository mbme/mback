package repository

import (
	"fmt"
	"mback/log"
	"mback/utils"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

const (
	DIR_PERM         = 0755
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

	err = os.Mkdir(r.GetRootPath(), DIR_PERM)
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
	return os.RemoveAll(r.GetRootPath())
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

	file_path := rec.GetRealPath()

	// first backup file if it exists
	if _, err := os.Stat(file_path); !os.IsNotExist(err) {
		err = utils.Backup(file_path)
		if err != nil {
			return err
		}
	}

	repo_file_path := rec.GetRepoPath()

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
		return fmt.Errorf("file %d is not installed", id)
	}

	file_path := rec.GetRealPath()

	err = os.Remove(file_path)
	if err != nil {
		return
	}

	return utils.RestoreBackup(file_path)
}

func (x *Repository) addRecord(file_path string) (err error) {
	if x.containsPath(file_path) {
		log.Info("'%v' already contains file %v, skipping", x.Name, file_path)
		return
	}

	// get file access rights and owner
	var file *os.File
	file, err = os.Open(file_path)
	if err != nil {
		return
	}

	var file_info os.FileInfo
	file_info, err = file.Stat()
	if err != nil {
		return
	}

	//FIXME
	uid := int(file_info.Sys().(*syscall.Stat_t).Uid)

	var u *user.User
	u, err = user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return
	}

	newId := x.getFreeId()
	newName := buildRepoFileName(file_path, newId)

	// TODO add groups processing
	record := new(Record)
	record.Id = newId
	record.SetRealPath(file_path)
	record.User = u.Username
	record.repository = x

	// copy file to repository
	err = utils.CopyFile(file_path, x.getRepoFilePath(newName))
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

	err = os.Remove(rec.GetRepoPath())
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
