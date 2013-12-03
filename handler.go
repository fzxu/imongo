package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ImgHandler struct {
}

func (h *ImgHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	switch req.Method {
	case "GET":
		h.handleGET(w, req)
	case "POST":
		h.handlePOST(w, req)
	}
}

func (h *ImgHandler) handleGET(w http.ResponseWriter, req *http.Request) {
	s := MgoSession.Copy()
	defer s.Close()

	name, path := h.convertPath(req.URL.Path)
	if name == "favicon.ico" {
		return
	}

	document, err := new(Document).Find(s, name, path)
	if err != nil {
		log.Panicln(err)
	}

	if document.ContentType != "" {
		w.Header().Set("Content-Type", document.ContentType)
	}

	var size *Size
	sizeReq := req.URL.Query().Get("size")
	if sizeReq == "" {
		size = Configuration.DefaultSize
	} else {
		size = Configuration.SizeMap[sizeReq]
	}

	origin, _, err := image.Decode(bytes.NewBuffer(document.Binary))
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	var img image.Image
	img = imaging.Resize(origin, size.Width, size.Height, imaging.CatmullRom)

	format := strings.ToLower(filepath.Ext(document.Name))
	if format != ".jpg" && format != ".jpeg" && format != ".png" {
		err = fmt.Errorf("unknown image format: %s", format)
		w.Write([]byte(err.Error()))
		return
	}

	switch format {
	case ".jpg", ".jpeg":
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: 95})
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: 95})
		}

	case ".png":
		err = png.Encode(w, img)
	}

	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (h *ImgHandler) handlePOST(w http.ResponseWriter, req *http.Request) {
	s := MgoSession.Copy()
	defer s.Close()

	name, path := h.convertPath(req.URL.Path)
	document, _ := new(Document).Find(s, name, path)
	document.Name = name
	document.Path = path

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Panicln(err)
	}
	document.Binary = data

	contentType := req.Header.Get("Content-Type")
	if contentType != "" {
		document.ContentType = contentType
	}

	err = document.Save(s)

	if err != nil {
		log.Panicln(err)
		io.WriteString(w, "error when storing \n")
	} else {
		io.WriteString(w, "stored\n")
	}

}

func (h *ImgHandler) convertPath(urlPath string) (string, string) {
	var path []string
	folders := strings.Split(urlPath, "/")
	for ind, folder := range folders {
		trimFolder := strings.Trim(folder, " ")
		if trimFolder != "" && ind != len(folders)-1 {
			path = append(path, strings.ToLower(folder))
		}
	}

	return folders[len(folders)-1], strings.Join(path, ",")
}
