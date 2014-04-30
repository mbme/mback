package utils

import (
	"mback/config"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	DIR_PERM  = 0755
	FILE_PERM = 0644
)

type File struct {
	path     string
	fileInfo os.FileInfo
}

type FileInfo struct {
	Owner string
}

func NewFile(path string) *File {
	f := &File{}
	f.path = simplifyPath(path)

	return f
}

func simplifyPath(file_path string) string {
	home_dir := filepath.Join("/home", config.GetConfig().User)

	if !strings.HasPrefix(file_path, home_dir) {
		return file_path
	}

	return strings.Replace(file_path, home_dir, "~", 1)
}

func (f *File) SimplifyPath() string {
	return simplifyPath(f.path)
}

func (f *File) GetPath() string {
	if strings.HasPrefix(f.path, "~/") {
		return filepath.Join("/home", config.GetConfig().User, f.path[2:])
	}

	return f.path
}

func (f *File) Mkdir() error {
	return os.Mkdir(f.GetPath(), DIR_PERM)
}

func (f *File) Remove() error {
	return os.Remove(f.GetPath())
}

func (f *File) RemoveAll() error {
	return os.RemoveAll(f.GetPath())
}

func (f *File) Symlink(link_path string) error {
	return os.Symlink(f.GetPath(), link_path)
}

func (f *File) GetInfo() (info *FileInfo, err error) {
	info = &FileInfo{}

	var file_info os.FileInfo
	file_info, err = f.stats()
	if err != nil {
		return
	}

	// FIXME
	uid := int(file_info.Sys().(*syscall.Stat_t).Uid)

	var u *user.User
	u, err = user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return
	}

	info.Owner = u.Username

	return
}

func (f *File) stats() (os.FileInfo, error) {
	if f.fileInfo != nil {
		return f.fileInfo, nil
	}

	return os.Stat(f.GetPath())
}

func (f *File) Exists() bool {
	_, err := f.stats()

	return !os.IsNotExist(err)
}

func (f *File) IsLink() bool {
	isLink, err := isSymlink(f.GetPath())
	if err != nil {
		return false
	}
	return isLink
}

func (f *File) IsDir() bool {
	stats, err := f.stats()
	if err != nil {
		return false
	}

	return stats.IsDir()
}

func (f *File) SameFile(other *File) bool {
	stats, err := f.stats()
	if err != nil {
		return false
	}

	otherStats, err := other.stats()
	if err != nil {
		return false
	}

	return os.SameFile(stats, otherStats)
}

func (f *File) CopyTo(dst string) error {
	if f.IsDir() {
		return copyDir(f.GetPath(), dst)
	} else {
		return copyFile(f.GetPath(), dst)
	}
}

func isSymlink(src string) (bool, error) {
	fi, err := os.Lstat(src)
	if err != nil {
		return false, err
	}

	return fi.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}
