package utils

import (
	"fmt"
	"io"
	"mback/log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	DIR_PERM  = 0755
	FILE_PERM = 0644
)

type RealFS struct{}

func (fs *RealFS) stats(file *File) (os.FileInfo, error) {
	return os.Stat(file.GetPath())
}

func (fs *RealFS) Mkdir(file *File, withParents bool) error {
	if withParents {
		return os.MkdirAll(file.GetPath(), DIR_PERM)
	} else {
		return os.Mkdir(file.GetPath(), DIR_PERM)
	}
}

func (fs *RealFS) Remove(file *File, recursive bool) error {
	if recursive {
		return os.RemoveAll(file.GetPath())
	} else {
		return os.Remove(file.GetPath())
	}
}

func (fs *RealFS) Move(src, dst *File) error {
	return os.Rename(src.GetPath(), dst.GetPath())
}

func (fs *RealFS) Symlink(src, dst *File) error {
	return os.Symlink(src.GetPath(), dst.GetPath())
}

func (fs *RealFS) Exists(file *File) bool {
	_, err := fs.stats(file)

	return !os.IsNotExist(err)
}

func (fs *RealFS) IsDir(file *File) (bool, error) {
	stats, err := fs.stats(file)

	if err != nil {
		return false, err
	}

	return stats.IsDir(), nil
}

func (fs *RealFS) IsSymlink(file *File) (bool, error) {
	fi, err := os.Lstat(file.GetPath())
	if err != nil {
		return false, err
	}

	return fi.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

func (fs *RealFS) IsSameFile(first, second *File) (bool, error) {
	stats, err := fs.stats(first)
	if err != nil {
		return false, err
	}

	otherStats, err := fs.stats(second)
	if err != nil {
		return false, err
	}

	return os.SameFile(stats, otherStats), nil
}

func (fs *RealFS) GetInfo(file *File) (*FileInfo, error) {
	info := &FileInfo{}

	fileInfo, err := fs.stats(file)
	if err != nil {
		return nil, err
	}

	// FIXME
	uid := int(fileInfo.Sys().(*syscall.Stat_t).Uid)

	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return nil, err
	}

	info.Owner = u.Username

	return info, nil
}

func (fs *RealFS) CopyFile(src, dst *File) error {
	s, err := os.Open(src.GetPath())
	if err != nil {
		return err
	}

	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()

	d, err := os.Create(dst.GetPath())
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func (fs *RealFS) CopyDir(src, dst *File) error {

	fileInfo, err := fs.stats(src)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%v is not a directory", src)
	}

	isLink, err := fs.IsSymlink(src)
	if err != nil {
		return err
	}
	if isLink {
		log.Info("%v is symlink, not copying", src)
		return nil
	}

	err = filepath.Walk(src.GetPath(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pathFile := NewFile(path)
			isLink, err := fs.IsSymlink(pathFile)
			if err != nil {
				return err
			}

			// skip dirs and symlinks
			if info.IsDir() || isLink {
				return nil
			}

			// get relative path
			relPath, error := filepath.Rel(src.GetPath(), path)
			if error != nil {
				return error
			}

			// create new path using base and relative path
			newPath := filepath.Join(dst.GetPath(), relPath)

			// create all required directories
			error = fs.Mkdir(NewFile(filepath.Dir(newPath)), true)
			if error != nil {
				return error
			}

			return fs.CopyFile(pathFile, NewFile(newPath))
		})

	return err
}
