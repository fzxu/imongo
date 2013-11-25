package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"labix.org/v2/mgo"
)

var (
	Configuration Config
	MgoSession    *mgo.Session
)

type Config struct {
	Address      string        `json:"address"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	DBHost       string        `json:"db_host"`
	DBName       string        `json:"db_name"`
	Collection   string        `json:"collection"`
}

func main() {
	// read config file
	configFile, err := os.Open(filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "arkxu", "imgongo", "config.json"))
	if err != nil {
		log.Panicln(err)
	}

	json.NewDecoder(configFile).Decode(&Configuration)

	// Initialize mongo connection
	log.Println(Configuration.DBHost)
	MgoSession, err = mgo.Dial(Configuration.DBHost)
	if err != nil {
		log.Panicln(err)
	}

	// start the server
	s := &http.Server{
		Addr:         Configuration.Address,
		Handler:      new(ImgHandler),
		ReadTimeout:  Configuration.ReadTimeout * time.Second,
		WriteTimeout: Configuration.WriteTimeout * time.Second,
	}
	log.Panicln(s.ListenAndServe())

}
