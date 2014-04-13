package utils

import (
	"errors"
	"fmt"
	"io"
	conf "mback/config"
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

	// simple function to add file to set only if it does not already there
	addFile := func(file string) {
		//check if file exists and is not a directory
		if file_info, err := os.Stat(file); err != nil || file_info.IsDir() {
			return
		}

		if _, contains := files[file]; !contains {
			files[file] = true
			log.Debug("Adding file %v", file)
		}
	}

	// all other params should be file names
	for _, name := range args {
		file_path := filepath.Join(baseDir, name)

		// check if real file path specified
		if _, err := os.Stat(file_path); err != nil {
			addFile(file_path)
			continue
		}

		// else try to use globs
		paths, err := filepath.Glob(file_path)
		if err != nil {
			log.Fatal("Bad pattern: %v", file_path)
		}

		for _, file := range paths {
			addFile(file)
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

func CopyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
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

func SimplifyPath(file_path string) string {
	home_dir := filepath.Join("/home", conf.USER)

	if !strings.HasPrefix(file_path, home_dir) {
		return file_path
	}

	return strings.Replace(file_path, home_dir, "~", 1)
}
