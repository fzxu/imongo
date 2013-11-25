package main

import (
	"io"
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
	io.WriteString(w, "Hello GET\n")
}

func (h *ImgHandler) handlePOST(w http.ResponseWriter, req *http.Request) {
	h.createDocument(req)
	io.WriteString(w, "Hello POST\n")
}

func (h *ImgHandler) createDocument(req *http.Request) {
	s := MgoSession.Copy()
	defer s.Close()

	name, path := h.convertPath(req.URL.Path)
	document := &Document{Name: name, Path: path}
	err := document.Save(s)
	if err != nil {
		log.Fatalln(err)
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
