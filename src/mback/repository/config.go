package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func (r *Repository) decode(reader io.Reader) (err error) {
	err = json.NewDecoder(reader).Decode(r)
	return
}

func (r *Repository) encode(f *os.File) (err error) {
	data, err := json.MarshalIndent(r, "", "  ")

	if err != nil {
		return
	}

	_, err = f.Write(data)

	return
}

func (r *Repository) readConfig() (err error) {
	if r.Records != nil {
		panic("config was already initialized earlier")
	}

	file_path := r.getConfigFilePath()

	file, err := os.Open(file_path)
	if err != nil {
		return
	}

	defer file.Close()

	err = r.decode(file)
	if err != nil {
		return
	}

	return
}

func (r *Repository) writeConfig() (err error) {
	if r.Records == nil {
		err = errors.New("Repository config is nil")
		return
	}

	file_path := r.getConfigFilePath()

	file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, FILE_PERM)
	if err != nil {
		return
	}

	defer file.Close()

	err = r.encode(file)
	return
}
