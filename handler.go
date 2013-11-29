package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	// w.Header().Set("Content-Type", document.ContentType)
	w.Write(document.Binary)
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
