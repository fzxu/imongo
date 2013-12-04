package main

import (
	"image"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testUrl1 = "http://localhost/foo/sdf/sdflkj/abc.png"
var testUrl2 = "http://localhost/foo/bar/Picture2.jpg"

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

	img, _, err := image.Decode(w.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if img.Bounds().Dx() != 655 {
		t.Fail()
	}
}

func TestUploadJPG(t *testing.T) {
	picture2, err := os.Open(filepath.Join(filepath.Dir(ConfigFileUrl), "testdata", "picture2.jpg"))
	if err != nil {
		log.Panicln(err)
	}

	req, err := http.NewRequest("POST", testUrl2, picture2)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "image/jpeg")

	w := httptest.NewRecorder()

	hander := &ImgHandler{}
	hander.ServeHTTP(w, req)
}

func TestGetJPG(t *testing.T) {

	ext := strings.ToLower(filepath.Ext(testUrl2))
	basename := strings.TrimSuffix(testUrl2, ext)

	req, err := http.NewRequest("GET", basename+"__345x123"+ext, nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	hander := &ImgHandler{}
	hander.ServeHTTP(w, req)
	if w.Body.Len() == 0 {
		t.Fail()
	}

	img, _, err := image.Decode(w.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if img.Bounds().Dx() != 345 {
		t.Errorf("expected %v got %v", 345, img.Bounds().Dx())
	}

	if img.Bounds().Dy() != 123 {
		t.Errorf("expected %v got %v", 123, img.Bounds().Dy())
	}
}
