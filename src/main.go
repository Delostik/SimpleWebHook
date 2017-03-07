// TODO config module separation; refresh periodic.
// TODO config file for parameters, use toml maybe
// TODO set workspace
// TODO martini log
// TODO daemon controller
// TODO automatic test cases

package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

type Hooks struct {
	Id        string
	ShellCode string
}

var hookMap map[string]string // [Id] => ShellCode

const confPath string = "config.conf"

var conf []Hooks

var wg sync.WaitGroup
var debugLog *log.Logger

func killInstance() {
	defer func() {
		if err := recover(); err != nil {
			debugLog.Fatalf("error while cleaning up exist instance: %v", err)
		}
	}()
	// TODO kill exist instance (ps -ef |grep xx |awk '{print $2}' |kill -9)
}

func main() {
	// log initialize
	os.MkdirAll("log", os.ModePerm)
	logFile, logErr := os.Create("log/SimpleWebhook.log")
	defer logFile.Close()
	if logErr != nil {
		log.Fatalf("open file error: %v", logErr)
	}
	debugLog = log.New(logFile, "[Debug]", log.Llongfile)
	hookMap = make(map[string]string)

	// kill exist instance
	killInstance()

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
	for _, hook := range conf {
		if len(hook.Id) > 0 && len(hook.ShellCode) > 0 {
			repoCnt++
			hookMap[hook.Id] = hook.ShellCode
			debugLog.Printf("New hook assigned: %v", hook.Id)
		}
	}
	debugLog.Printf("%v hooks registered", repoCnt)

	// martini initialize
	m := martini.Classic()

	// TODO multi request(duplicated) handle
	// TODO separate from main package
	// TODO return value structure
	m.Get("/:id", func(params martini.Params) string {
		if id, existId := params["id"]; existId {
			if shellCode, existHook := hookMap[id]; existHook {
				debugLog.Printf("Receive hook call: %v", id)
				output, err := exec.Command("/bin/sh", "-c", shellCode).Output()
				if err != nil {
					debugLog.Printf("Running shellcode error %v: %v", id, err)
					return "error"
				}
				debugLog.Printf("Running result for %v: %v", id, output)
				return "done"
			}
			debugLog.Printf("unknow id %v", id)
			return "not found"
		}
		return "invalid request"
	})
	m.RunOnAddr(":8009")
}
