package utils

import (
	"fmt"
	"io"
	"mback/log"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	//TODO handle dirs here
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

func copyDir(src, dst string) error {

	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%v is not a directory", src)
	}

	isLink, err := isSymlink(src)
	if err != nil {
		return err
	}
	if isLink {
		log.Debug("%v is symlink, skipping", src)
		return nil
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isLink, err := isSymlink(path)
		if err != nil {
			return err
		}

		// skip dirs and symlinks
		if info.IsDir() || isLink {
			return nil
		}

		// get relative path
		relPath, error := filepath.Rel(src, path)
		if error != nil {
			return error
		}

		// create new path using base and relative path
		newPath := filepath.Join(dst, relPath)

		// create all required directories
		error = os.MkdirAll(filepath.Dir(newPath), DIR_PERM)
		if error != nil {
			return error
		}

		return copyFile(path, newPath)
	})

	return err
}
