package repository

import (
	"mback/log"
	"mback/utils"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func (x *Repository) addFile(file_path string) (err error) {
	if x.containsPath(file_path) {
		log.Info("Repository %v already contains file %v, skipping", x.Name, file_path)
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

	//fixme
	uid := int(file_info.Sys().(*syscall.Stat_t).Uid)

	var u *user.User
	u, err = user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return
	}

	newId := x.freeId()
	newName := buildRepoFileName(file_path, newId)

	// TODO add groups processing
	record := new(Record)
	record.Id = newId
	record.SetPath(file_path)
	record.User = u.Username

	// copy file to repository
	err = utils.CopyFile(file_path, x.getRepoFile(newName))
	if err != nil {
		return
	}

	x.addRecord(record)

	log.Info("Added file %v as %v", file_path, newName)
	return
}

func (r *Repository) removeRecord(id int) (err error) {
	rec, pos, err := r.GetRecord(id)
	if err != nil {
		return
	}

	err = os.Remove(r.getRepoFile(rec.GetRepoFileName()))
	if err != nil {
		return
	}

	// delete file from config
	r.Records = append(r.Records[:pos], r.Records[pos+1:]...)

	log.Debug("Removed record %v", rec)

	return
}
