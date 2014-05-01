package repository

import (
	"reflect"
	"testing"
)

func TestGetRecord(t *testing.T) {
	repo := &Repository{Name: "test"}

	id := 3
	repo.Records = []*Record{&Record{Id: 2}, &Record{Id: id}}

	r, pos, err := repo.GetRecord(id)

	if err != nil {
		t.Error("Unexpected error", err)
	}

	if pos != 1 {
		t.Error("Expected pos 1 received", pos)
	}

	if r == nil || r.Id != id {
		t.Error("Expected record id", id, "received", r.Id)
	}

	// get not existing id
	r, pos, err = repo.GetRecord(-1)

	if err == nil || pos != -1 || r != nil {
		t.Error("Expected not to found record, but found record", r, "pos", pos, "error", err)
	}
}

func TestListIds(t *testing.T) {
	repo := &Repository{Name: "test"}

	repo.Records = []*Record{&Record{Id: 2}, &Record{Id: 3}}

	ids := repo.ListIds()

	expected := []int{2, 3}
	if !reflect.DeepEqual(ids, expected) {
		t.Error("Expected", expected, "received", ids)
	}
}

func TestGetFreeId(t *testing.T) {
	repo := &Repository{Name: "test"}
	repo.Records = []*Record{&Record{Id: 2}, &Record{Id: 3}}

	id := repo.getFreeId()

	expected := 4
	if id != expected {
		t.Error("Expected id", expected, "received", id)
	}

	// get free id on empty repo
	repo.Records = []*Record{}

	id = repo.getFreeId()
	expected = 0
	if id != expected {
		t.Error("Expected id", expected, "received", id)
	}
}
