package main

import (
	"github.com/codegangsta/martini"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"sync"
)

type Repository struct {
	Id          string
	RepoPath    string
	ShellCode   string
}

const confPath string = "config.conf"
var conf []Repository

var wg sync.WaitGroup


func main() {
	// log initialize
	os.MkdirAll("log", os.ModePerm)
	logFile, logErr  := os.Create("log/SimpleWebhook.log")
	defer logFile.Close()
	if logErr != nil {
		log.Fatalf("open file error: %v", logErr)
	}
	debugLog := log.New(logFile,"[Debug]",log.Llongfile)

	// config initialize
	_, confErr := os.Stat(confPath)
	if confErr != nil {
		if os.IsNotExist(confErr) {
			debugLog.Println("Config file not found. Create empty file automatically.")
		} else {
			debugLog.Fatalf("Something wrong with config.conf: \n%v\n", confErr)
		}
	}

	f, fErr := os.Open(confPath)
	defer f.Close()
	if fErr != nil {
		debugLog.Fatalf("Something wrong with config.conf: \n%v\n", fErr)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		debugLog.Fatalf("Something wrong with config.conf: \n%v\n", err)
	}
	f.Close()

	err = json.Unmarshal(content, &conf)
	if err != nil {
		debugLog.Fatalf("Syntax error in config.conf: \n%v\n", err)
	}

	repoCnt := 0
	for _, repo := range conf {

	}

	// martini initialize
	m := martini.Classic()
	m.Get("/", func() string {
		return "fuck"
	})
	m.RunOnAddr(":8081")
}