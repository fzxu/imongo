package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Address      string        `json:"address"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

func main() {
	// read config file
	configFile, err := os.Open(filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "arkxu", "imgongo", "config.json"))
	if err != nil {
		log.Fatal(err)
	}

	var configuration Config
	json.NewDecoder(configFile).Decode(&configuration)

	// start the server
	s := &http.Server{
		Addr:         configuration.Address,
		Handler:      new(ImgHandler),
		ReadTimeout:  configuration.ReadTimeout * time.Second,
		WriteTimeout: configuration.WriteTimeout * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
