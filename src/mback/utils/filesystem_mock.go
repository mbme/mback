package utils

type TestFS struct {
	OnMkdir   func(file *File, withParents bool) error
	OnSymlink func(src, dst *File) error
	OnMove    func(src, dst *File) error

	OnRead  func(src *File) ([]byte, error)
	OnWrite func(file *File, data []byte) error

	OnCopyFile func(src, dst *File) error
	OnCopyDir  func(src, dst *File) error

	OnRemove func(file *File, recursive bool) error

	OnGetInfo    func(file *File) (*FileInfo, error)
	OnExists     func(file *File) bool
	OnIsDir      func(file *File) (bool, error)
	OnIsSymlink  func(file *File) (bool, error)
	OnIsSameFile func(first, second *File) (bool, error)
}

func UninstallFs() {
	Fs = nil
}

func (fs *TestFS) Mkdir(file *File, withParents bool) error {
	return fs.OnMkdir(file, withParents)
}

func (fs *TestFS) Symlink(src, dst *File) error {
	return fs.OnSymlink(src, dst)
}
func (fs *TestFS) Move(src, dst *File) error {
	return fs.OnMove(src, dst)
}

func (fs *TestFS) Read(src *File) ([]byte, error) {
	return fs.OnRead(src)
}
func (fs *TestFS) Write(file *File, data []byte) error {
	return fs.OnWrite(file, data)
}

func (fs *TestFS) CopyFile(src, dst *File) error {
	return fs.OnCopyFile(src, dst)
}
func (fs *TestFS) CopyDir(src, dst *File) error {
	return fs.OnCopyDir(src, dst)
}

func (fs *TestFS) Remove(file *File, recursive bool) error {
	return fs.OnRemove(file, recursive)
}

func (fs *TestFS) GetInfo(file *File) (*FileInfo, error) {
	return fs.OnGetInfo(file)
}
func (fs *TestFS) Exists(file *File) bool {
	return fs.OnExists(file)
}
func (fs *TestFS) IsDir(file *File) (bool, error) {
	return fs.OnIsDir(file)
}
func (fs *TestFS) IsSymlink(file *File) (bool, error) {
	return fs.OnIsSymlink(file)
}
func (fs *TestFS) IsSameFile(first, second *File) (bool, error) {
	return fs.OnIsSameFile(first, second)
}
