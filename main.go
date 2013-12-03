package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"labix.org/v2/mgo"
)

var (
	ConfigFileUrl string
	Configuration Config
	MgoSession    *mgo.Session
)

type Config struct {
	Host         string           `json:"host"`
	Port         string           `json:"port"`
	ReadTimeout  time.Duration    `json:"read_timeout"`
	WriteTimeout time.Duration    `json:"write_timeout"`
	DBHost       string           `json:"db_host"`
	DBName       string           `json:"db_name"`
	Collection   string           `json:"collection"`
	DefaultSize  *Size            `json:"default_size"`
	SizeMap      map[string]*Size `json:"size_map"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func main() {
	address := Configuration.Host + ":" + Configuration.Port
	log.Println(address)
	// start the server
	s := &http.Server{
		Addr:         address,
		Handler:      new(ImgHandler),
		ReadTimeout:  Configuration.ReadTimeout * time.Second,
		WriteTimeout: Configuration.WriteTimeout * time.Second,
	}
	log.Panicln(s.ListenAndServe())

}

func init() {

	var port string
	flag.StringVar(&port, "p", "9020", "The port which image server is running on")
	flag.StringVar(&ConfigFileUrl, "c", filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "arkxu", "imongo", "config.json"), "Specify configuration file")

	flag.Parse()

	// read config file
	configFile, err := os.Open(ConfigFileUrl)
	if err != nil {
		log.Panicln(err)
	}

	json.NewDecoder(configFile).Decode(&Configuration)
	Configuration.Port = port

	// Initialize mongo connection
	log.Println(Configuration.DBHost)
	MgoSession, err = mgo.Dial(Configuration.DBHost)
	if err != nil {
		log.Panicln(err)
	}
}
