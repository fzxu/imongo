package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"labix.org/v2/mgo"
)

var (
	ConfigFileUrl string
	Configuration *Config
	MgoSession    *mgo.Session
)

type ConfigAll struct {
	Test *Config `json:"test"`
	Prod *Config `json:"prod"`
}

type Config struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	DBHost       string        `json:"db_host"`
	DBName       string        `json:"db_name"`
	Collection   string        `json:"collection"`
	StoredSize   *Size         `json:"stored_size"`
	CacheFolder  string        `json:"cache_folder"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Option int // 0 means Resize, 1 means Thumbnail
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
	var env string
	flag.StringVar(&port, "p", "9020", "The port which image server is running on")
	flag.StringVar(&env, "e", "test", "The running environment")
	flag.StringVar(&ConfigFileUrl, "c", filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "arkxu", "imongo", "config.json"), "Specify configuration file")

	flag.Parse()

	// read config file
	configFile, err := os.Open(ConfigFileUrl)
	if err != nil {
		log.Panicln(err)
	}

	var configAll ConfigAll
	json.NewDecoder(configFile).Decode(&configAll)
	switch strings.ToLower(env) {
	default:
		Configuration = configAll.Test
	case "prod":
		Configuration = configAll.Prod
	}
	Configuration.Port = port

	// Initialize mongo connection
	log.Println(Configuration.DBHost)
	MgoSession, err = mgo.Dial(Configuration.DBHost)
	if err != nil {
		log.Panicln(err)
	}
}
