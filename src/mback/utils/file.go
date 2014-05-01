package utils

import (
	"path/filepath"
	"strings"
)

type File struct {
	path string
}

type FileInfo struct {
	Owner string
}

func NewFile(path string) *File {
	f := &File{}
	f.path = SimplifyPath(path)

	return f
}

func SimplifyPath(file_path string) string {
	home_dir := filepath.Join("/home", Conf.User)

	if !strings.HasPrefix(file_path, home_dir) {
		return file_path
	}

	return strings.Replace(file_path, home_dir, "~", 1)
}

func (f *File) GetPath() string {
	if strings.HasPrefix(f.path, "~/") {
		return filepath.Join("/home", Conf.User, f.path[2:])
	}

	return f.path
}

func (f *File) Mkdir() error {
	return Fs.Mkdir(f, false)
}

func (f *File) Remove(recursive bool) error {
	return Fs.Remove(f, recursive)
}

func (f *File) Symlink(link_path string) error {
	return Fs.Symlink(f, NewFile(link_path))
}

func (f *File) GetInfo() (info *FileInfo, err error) {
	return Fs.GetInfo(f)
}

func (f *File) Exists() bool {
	return Fs.Exists(f)
}

func (f *File) IsLink() (bool, error) {
	return Fs.IsSymlink(f)
}

func (f *File) IsDir() (bool, error) {
	return Fs.IsDir(f)
}

func (f *File) SameFile(other *File) (bool, error) {
	return Fs.IsSameFile(f, other)
}

func (f *File) CopyTo(dst string) error {
	isDir, err := f.IsDir()

	if err == nil && isDir {
		return Fs.CopyDir(f, NewFile(dst))
	} else {
		return Fs.CopyFile(f, NewFile(dst))
	}
}
