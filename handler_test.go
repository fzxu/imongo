package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var testUrl1 = "http://localhost/foo/sdf/sdflkj/abc.png"

func TestHandlePOST(t *testing.T) {
	picture1, err := os.Open(filepath.Join(filepath.Dir(ConfigFileUrl), "testdata", "picture1.png"))
	if err != nil {
		log.Panicln(err)
	}

	req, err := http.NewRequest("POST", testUrl1, picture1)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	hander := &ImgHandler{}
	hander.ServeHTTP(w, req)
}

func TestHandleGET(t *testing.T) {
	req, err := http.NewRequest("GET", testUrl1, nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	hander := &ImgHandler{}
	hander.ServeHTTP(w, req)
	if w.Body.Len() == 0 {
		t.Fail()
	}
}
