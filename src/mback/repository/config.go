package repository

import (
	"errors"
	"mback/utils"
)

func (r *Repository) readConfig() error {
	if r.Records != nil {
		panic("config was already initialized earlier")
	}

	data, err := utils.Fs.Read(r.getConfigFile())
	if err != nil {
		return err
	}

	return utils.Decode(data, r)
}

func (r *Repository) writeConfig() error {
	if r.Records == nil {
		return errors.New("Repository config is nil")
	}

	data, err := utils.Encode(r)
	if err != nil {
		return err
	}

	return utils.Fs.Write(r.getConfigFile(), data)
}
