package repository

import (
	"errors"
	"mback/utils"
	"testing"
)

const USER = "tester"

func init() {
	utils.Conf = &utils.Config{User: USER}
}

var tests = []struct {
	path   string
	result string
}{
	{"/etc/data/1", "/etc/data/1"},
	{"/", "/"},
	{"/home/" + USER + "/data", "~/data"},
}

func TestSetPath(t *testing.T) {
	for _, pair := range tests {
		rec := &Record{}
		rec.SetPath(pair.path)

		if rec.Path != pair.result {
			t.Error(
				"For", pair.path,
				"expected", pair.result,
				"got", rec.Path,
			)
		}
	}
}

func TestGetPath(t *testing.T) {
	for _, pair := range tests {
		rec := &Record{}
		rec.SetPath(pair.path)

		path := rec.GetFile().GetPath()
		if path != pair.path {
			t.Error(
				"For", pair.path,
				"got", path,
			)
		}
	}
}

func TestRecordIsInstalled(t *testing.T) {

	rec := &Record{}
	rec.SetPath("/test/rec")

	repo := Repository{}
	repo.Records = []*Record{rec}
	rec.repository = &repo

	type fRes struct {
		result bool
		err    error
	}
	fErr := errors.New("")

	//install mock
	c_exists := -1
	c_isSymlink := -1
	c_isSameFile := -1
	utils.Fs = &utils.TestFS{
		OnExists: func(f *utils.File) bool {
			c_exists += 1
			return c_exists != 0
		},
		OnIsSymlink: func(f *utils.File) (bool, error) {
			c_isSymlink += 1
			data := []fRes{
				{false, nil},
				{false, fErr},
				{true, nil},
				{true, nil},
				{true, nil},
			}[c_isSymlink]

			return data.result, data.err
		},
		OnIsSameFile: func(first, second *utils.File) (bool, error) {
			c_isSameFile += 1
			data := []fRes{
				{false, nil},
				{false, fErr},
				{true, nil},
			}[c_isSameFile]

			return data.result, data.err
		},
	}
	defer utils.UninstallFs()

	results := []fRes{
		{false, nil},
		{false, nil},
		{false, fErr},
		{false, nil},
		{false, fErr},
		{true, nil},
	}

	for _, result := range results {
		installed, err := rec.IsInstalled()

		errNil := err == nil
		resErrNil := result.err == nil

		if errNil != resErrNil {
			t.Error("Expected", result.err, "actual", err)
			continue
		}
		if installed != result.result {
			t.Error("Expected", result.result, "actual", installed)
		}
	}
}
