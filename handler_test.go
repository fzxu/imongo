package main

import (
	"image"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var testUrl1 = "http://localhost/foo/sdf/sdflkj/abc.png"
var testUrl2 = "http://localhost/foo/bar/Picture2.png"

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

	if img.Bounds().Dx() != Configuration.DefaultSize.Width {
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
	req, err := http.NewRequest("GET", testUrl2+"?size=l", nil)
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

	if img.Bounds().Dx() != Configuration.SizeMap["l"].Width {
		t.Errorf("expected %v got %v", Configuration.DefaultSize.Width, img.Bounds().Dx())
	}
}
