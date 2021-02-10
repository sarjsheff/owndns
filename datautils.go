package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var configpath = flag.String("c", "./owndns.config", "`Path` to config file.")

type Rule struct {
	Ip         string
	RejectName string `json:",omitempty"`
	AcceptOnly string `json:",omitempty"`
}

type Config struct {
	RedirectIp  string   `json:",omitempty"`
	RelayDns    []string `json:",omitempty"`
	RejectNames []string `json:",omitempty"`
	DNSListen   string   `json:",omitempty"`
	HttpListen  string   `json:",omitempty"`
	Logging     bool     `json:",omitempty"`
	Rules       []Rule
}

var loglines int = 100
var logstate []Item
var logstore chan Item = make(chan Item, loglines)
var stat map[string]map[string]int

// var configpath string = "owndns.config"

var configwriter chan Config
var config Config

func Init() {

	flag.Parse()
	fmt.Printf("OwnDns v%s\n\n", BuildTime)
	if _, err := os.Stat(*configpath); os.IsNotExist(err) {

		fmt.Println("Config not found.")

		flag.PrintDefaults()
		fmt.Println()
		os.Exit(-1)
	}
	stat = make(map[string]map[string]int)
	configwriter = make(chan Config)

	go func() {
		for {
			x := <-configwriter
			j, err := json.MarshalIndent(x, "", "  ")
			if err == nil {
				ioutil.WriteFile(*configpath, j, 0600)
				config = x
			}

		}
	}()

	var err error

	config = Config{DNSListen: ":53", HttpListen: ":8081"}
	content, err := ioutil.ReadFile(*configpath)
	if err != nil {
		configwriter <- config
	} else {
		json.Unmarshal(content, &config)
	}
	// log.Println(config)
}

func Validate(name string, ip string) bool {
	ret := true
	for _, s := range GetRules(ip) {
		if s.AcceptOnly == "*" {
			return true
		}
		if s.RejectName == "*" {
			ret = false
		} else {
			if len(s.RejectName) > 2 {
				if strings.Contains(name, s.RejectName) {
					ret = false
				}
			} else if len(s.AcceptOnly) > 2 {
				if strings.Contains(name, s.AcceptOnly) {
					return true
				}
			}
		}
	}

	if stringInSlice(name, config.RejectNames) {
		ret = false
	} else if ret == true {
		ret = true
	}
	return ret
}

func GetRedirectIp() (string, error) {
	if config.RedirectIp == "" {
		return "", errors.New("Empty")
	} else {
		return config.RedirectIp, nil
	}
}

func GetRejectnames() []string {
	return config.RejectNames
}

func GetRelaydns() []string {
	return config.RelayDns
}

func Upsert(key string, value string) error {
	return nil
}

func Stat(ip string, name string) {
	if _, ok := stat[ip]; !ok {
		stat[ip] = make(map[string]int)
		stat[ip][name] = 0
	}

	stat[ip][name] = stat[ip][name] + 1
}

func GetLogStat() map[string]map[string]int {
	return stat
}

func GetLastLog() []Item {
	tmp := []Item{}
	notend := true
	for notend {
		select {
		case item := <-logstore:
			tmp = append(tmp, item)
			break
		default:
			notend = false
			break
		}
	}

	logstate = append(tmp, logstate...)
	if len(logstate) > loglines {
		logstate = logstate[:loglines]
	}

	return logstate
}

func Log(item Item) {
	if config.Logging {
		logstore <- item
		if len(logstore) > loglines-1 {
			<-logstore
		}
	}
}

func ArrAdd(arr *[]string, val string) {
	*arr = append(*arr, val)
	configwriter <- config
}

func ArrDel(arr *[]string, id int) {
	copy((*arr)[id:], (*arr)[id+1:])
	(*arr)[len((*arr))-1] = ""
	*arr = (*arr)[:len((*arr))-1]
	configwriter <- config
}

func ArrSave(arr *[]string, id int, name string) {
	(*arr)[id] = name
	configwriter <- config
}

func GetRules(ip string) (ret []Rule) {

	for _, s := range config.Rules {
		if s.Ip == ip {
			ret = append(ret, s)
		}
	}
	return ret
}
