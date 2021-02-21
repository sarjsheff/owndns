package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var configpath = flag.String("c", "./owndns.config", "`Path` to config file.")

// // Rule - Filtering rule.
// type Rule struct {
// 	IP         string
// 	RejectName string `json:",omitempty"`
// 	AcceptOnly string `json:",omitempty"`
// }

// Config - Configuration.
type Config struct {
	RedirectIP  string   `json:",omitempty"`
	RelayDNS    []string `json:",omitempty"`
	RejectNames []string `json:",omitempty"`
	DNSListen   string   `json:",omitempty"`
	HTTPListen  string   `json:",omitempty"`
	Logging     bool     `json:",omitempty"`
	Rules       []interface{}
	Users       map[string]interface{}
}

var loglines int = 100
var logstate []Item
var logstore chan Item = make(chan Item, loglines)
var stat map[string]map[string]int

// var configpath string = "owndns.config"

var configwriter chan Config
var config Config

// Init - Initialize config.
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
	config = Config{DNSListen: ":53", HTTPListen: ":8081"}
	content, err := ioutil.ReadFile(*configpath)
	if err != nil {
		configwriter <- config
	} else {
		json.Unmarshal(content, &config)
	}
	// log.Println(config)
	if !isAdminExist() {
		if config.Users == nil {
			config.Users = map[string]interface{}{}
		}
		if config.Users["admin"] == nil {
			config.Users["admin"] = map[string]interface{}{
				"PasswordHash": fmt.Sprintf("%x", sha256.Sum256([]byte("password"))),
				"IsAdmin":      true,
			}
		}
		configwriter <- config
	}
}

func isAdminExist() bool {
	if config.Users == nil {
		return false
	} else {
		for _, v := range config.Users {
			if v.(map[string]interface{})["IsAdmin"] == true {
				return true
			}
		}
		return false
	}
}

// Validate - Validate DNS request.
func Validate(name string, ip string) bool {
	ret := true
	for _, s := range GetRules(ip) {
		if s["AcceptOnly"] == "*" {
			return true
		}
		if s["RejectName"] == "*" {
			ret = false
		} else {
			if s["RejectName"] != nil && len(s["RejectName"].(string)) > 2 {
				if strings.Contains(name, s["RejectName"].(string)) {
					ret = false
				}
			} else if s["AcceptOnly"] != nil && len(s["AcceptOnly"].(string)) > 2 {
				if strings.Contains(name, s["AcceptOnly"].(string)) {
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

// GetRedirectIP - Get redirect ip address.
func GetRedirectIP() (string, error) {
	if config.RedirectIP == "" {
		return "", errors.New("Empty")
	}
	return config.RedirectIP, nil

}

// GetRejectnames - Pattern list of dns name.
func GetRejectnames() []string {
	return config.RejectNames
}

// GetRelaydns - IP:PORT list of relay dns.
func GetRelaydns() []string {
	return config.RelayDNS
}

// Stat - Add DNS request to statistics.
func Stat(ip string, name string) {
	if _, ok := stat[ip]; !ok {
		stat[ip] = make(map[string]int)
		stat[ip][name] = 0
	}

	stat[ip][name] = stat[ip][name] + 1
}

// Log - Add log record.
func Log(item Item) {
	if config.Logging {
		logstore <- item
		if len(logstore) > loglines-1 {
			<-logstore
		}
	}
}

// GetLogStat - Get statistics map.
func GetLogStat() map[string]map[string]int {
	return stat
}

// GetLastLog - Array of last log records.
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

// ArrAddObject - Add item to array of objects.
func ArrAddObject(arr *[]interface{}, val interface{}) {
	*arr = append(*arr, &val)
	log.Println(config.Rules)
	configwriter <- config
}

// ArrAdd - Add item to array of string.
func ArrAdd(arr *[]string, val string) {
	*arr = append(*arr, val)
	configwriter <- config
}

// ArrDelObject - Del item from array of objects.
func ArrDelObject(arr *[]interface{}, id int) {
	copy((*arr)[id:], (*arr)[id+1:])
	(*arr)[len((*arr))-1] = ""
	*arr = (*arr)[:len((*arr))-1]
	configwriter <- config
}

// ArrDel - Del item from array of strings.
func ArrDel(arr *[]string, id int) {
	copy((*arr)[id:], (*arr)[id+1:])
	(*arr)[len((*arr))-1] = ""
	*arr = (*arr)[:len((*arr))-1]
	configwriter <- config
}

// ArrSaveObject - Set item in array of objects.
func ArrSaveObject(arr *[]interface{}, id int, item interface{}) {
	(*arr)[id] = item
	configwriter <- config
}

// ArrSave - Set item in array of strings.
func ArrSave(arr *[]string, id int, name string) {
	(*arr)[id] = name
	configwriter <- config
}

// GetRules - List filtering rules.
func GetRules(ip string) (ret []map[string]interface{}) {
	for _, s := range config.Rules {
		if s.(map[string]interface{})["IP"].(string) == ip {
			ret = append(ret, s.(map[string]interface{}))
		}
	}
	return ret
}
