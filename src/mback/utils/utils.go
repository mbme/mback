package utils

import (
	"errors"
	"fmt"
	"mback/log"
	"os"
	"path/filepath"
	"strings"
)

const BACKUP_EXT = ".mback"

func Confirmation(msg string) (resp bool) {
	log.Info(msg + " (y/n): ")

	var str string
	done := false
	for !done {
		fmt.Scanf("%s", &str)

		switch str {
		case "y":
			done = true
			resp = true
		case "n":
			done = true
			resp = false
		default:
			log.Info("wrong answer, should be y/n")
		}
	}
	return
}

func ListFiles(baseDir string, args ...string) (result []string, err error) {

	files := make(map[string]bool, 10)

	// all other params should be file names
	for _, name := range args {
		file_path := filepath.Join(baseDir, name)

		file := NewFile(file_path)

		if !file.Exists() {
			continue
		}

		// adding file only if not exist yet
		if _, contains := files[file_path]; !contains {
			files[file_path] = true
			if file.IsDir() {
				log.Debug("Adding dir %v", file_path)
			} else {
				log.Debug("Adding file %v", file_path)
			}
		}
	}

	result = make([]string, len(files))
	i := 0
	for key := range files {
		result[i] = key
		i++
	}

	return
}

// Create file backup
func Backup(file_path string) error {
	return os.Rename(file_path, getBackupFile(file_path))
}

func RestoreBackup(file_path string) error {
	return os.Rename(getBackupFile(file_path), file_path)
}

func getBackupFile(file_path string) (backup_file string) {
	base, file := filepath.Split(file_path)

	// check if file is hidden, and if not then hide it
	if file[0] == '.' {
		backup_file = file + BACKUP_EXT
	} else {
		backup_file = "." + file + BACKUP_EXT
	}

	return filepath.Join(base, backup_file)
}

func GetWorkingDir() (wd string, err error) {
	// list of env vars in format name=var
	env_vars := os.Environ()

	pair := func(env_var string) (key, val string) {
		pos := strings.Index(env_var, "=")
		if pos == -1 {
			key = env_var
			return
		}

		key = env_var[:pos]
		val = env_var[pos+1:]

		return
	}

	for _, env_var := range env_vars {
		key, val := pair(env_var)
		if key == "PWD" {
			wd = val
			return
		}
	}

	err = errors.New("can't find env variable PWD")
	return
}
