package utils

type FileSystem interface {
	Mkdir(file *File, withParents bool) error
	Symlink(src, dst *File) error
	Move(src, dst *File) error

	Read(src *File) ([]byte, error)
	Write(file *File, data []byte) error

	CopyFile(src, dst *File) error
	CopyDir(src, dst *File) error

	Remove(file *File, recursive bool) error

	GetInfo(file *File) (*FileInfo, error)
	Exists(file *File) bool
	IsDir(file *File) (bool, error)
	IsSymlink(file *File) (bool, error)
	IsSameFile(first, second *File) (bool, error)
}

var Fs FileSystem
