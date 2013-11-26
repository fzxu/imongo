package main

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	s := MgoSession.Copy()
	defer s.Close()

	document := &Document{Name: "name", Path: "foo, bar, folder1", CreatedAt: time.Now()}
	document.Save(s)

	result, err := new(Document).Find(s, "name", "foo, bar, folder1")
	if err != nil {
		t.Fail()
	}
	if result.Name != "name" {
		t.Fail()
	}

	result2, err := new(Document).Find(s, "name", "foo, bar")
	if err == nil {
		t.Fail()
	}
	if result2.Name == "name" {
		t.Fail()
	}
}
