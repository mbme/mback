package repository

import (
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
