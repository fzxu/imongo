package main

import (
	"io"
	"log"
	"net/http"
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
	io.WriteString(w, "Hello POST\n")
}
