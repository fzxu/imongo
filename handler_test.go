package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandlePOST(t *testing.T) {
	picture1, err := os.Open(filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "arkxu", "imgongo", "testdata", "picture1.png"))
	if err != nil {
		log.Panicln(err)
	}

	req, err := http.NewRequest("POST", "http://example.com/foo/sdf/sdflkj/abc.jpg", picture1)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	hander := &ImgHandler{}
	hander.ServeHTTP(w, req)
}
