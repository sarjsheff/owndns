package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
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

type SafeConfig struct {
	mu sync.Mutex
	//config     Config
	config     map[string]interface{}
	configpath string
}

func NewSafeConfig(configpath string) *SafeConfig {
	ret := new(SafeConfig)
	ret.configpath = configpath
	cfg := map[string]interface{}{
		"DNSListen":  ":53",
		"HTTPListen": ":8081",
	}
	//Config{DNSListen: ":53", HTTPListen: ":8081"}

	ret.config = cfg

	content, err := ioutil.ReadFile(configpath)

	if err != nil {
		ret.flush()
	} else {
		json.Unmarshal(content, &cfg)
	}

	if !isAdminExist() {
		if cfg["Users"] == nil {
			cfg["Users"] = map[string]interface{}{}
		}
		if cfg["Users"].(map[string]interface{})["admin"] == nil {
			cfg["Users"].(map[string]interface{})["admin"] = map[string]interface{}{
				"PasswordHash": fmt.Sprintf("%x", sha256.Sum256([]byte("password"))),
				"IsAdmin":      true,
			}
		}
		ret.flush()
	}

	return ret
}

func (c *SafeConfig) flush() {
	//c.mu.Lock()
	j, err := json.MarshalIndent(c.config, "", "  ")
	if err == nil {
		ioutil.WriteFile(c.configpath, j, 0600)
	}
	//c.mu.Unlock()
}

func (c *SafeConfig) IsAdminExist() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config["Users"] == nil {
		return false
	}
	for _, v := range c.config["Users"].(map[string]interface{}) {
		if v.(map[string]interface{})["IsAdmin"] == true {
			return true
		}
	}
	return false
}

func (c *SafeConfig) FindReject(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return stringInSlice(name, c.config["RejectNames"].([]interface{}))
}

func (c *SafeConfig) GetRedirectIP() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config["RedirectIP"].(string) == "" {
		return "", errors.New("Empty")
	}
	return c.config["RedirectIP"].(string), nil
}

func (c *SafeConfig) GetRejectnames() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config["RejectNames"].([]string)
}

func (c *SafeConfig) GetRelaydns() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config["RelayDNS"].([]interface{})
}

func (c *SafeConfig) IsLogging() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config["Logging"].(bool)
}

func (c *SafeConfig) AddObject(name string, o interface{}) {
	c.mu.Lock()
	if c.config[name] == nil {
		c.config[name] = []interface{}{}
	}
	c.config[name] = append(c.config[name].([]interface{}), o)
	c.flush()
	c.mu.Unlock()
}

func (c *SafeConfig) DelObject(name string, id int) {
	c.mu.Lock()
	if c.config[name] != nil {
		copy(c.config[name].([]interface{})[id:], c.config[name].([]interface{})[id+1:])
		c.config[name].([]interface{})[len(c.config[name].([]interface{}))-1] = ""
		c.config[name] = c.config[name].([]interface{})[:len(c.config[name].([]interface{}))-1]
	}
	c.flush()
	c.mu.Unlock()
}

func (c *SafeConfig) SetObject(name string, id int, o interface{}) {
	c.mu.Lock()
	if c.config[name] != nil {
		c.config[name].([]interface{})[id] = o
	}
	c.flush()
	c.mu.Unlock()
}

func (c *SafeConfig) Get(name string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config[name]
}

func (c *SafeConfig) Set(name string, o interface{}) {
	c.mu.Lock()
	c.config[name] = o
	c.flush()
	c.mu.Unlock()
}

func (c *SafeConfig) GetRules(ip string) (ret []map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, s := range c.config["Rules"].([]interface{}) {
		if s.(map[string]interface{})["IP"].(string) == ip {
			ret = append(ret, s.(map[string]interface{}))
		}
	}
	return ret
}

func (c *SafeConfig) GetUser(name string) map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config["Users"] != nil {
		if c.config["Users"].(map[string]interface{})[name] != nil {
			return c.config["Users"].(map[string]interface{})[name].(map[string]interface{})
		}
	}
	return nil
}

func (c *SafeConfig) AddUser(username string, passwordhash string, isadmin bool) {
	c.mu.Lock()
	if c.config["Users"] == nil {
		c.config["Users"] = map[string]interface{}{}
	}
	c.config["Users"].(map[string]interface{})[username] = map[string]interface{}{
		"PasswordHash": passwordhash, //fmt.Sprintf("%x", sha256.Sum256([]byte(item["Password"].(string)))),
		"IsAdmin":      isadmin,
	}
	c.flush()
	c.mu.Unlock()
}
func (c *SafeConfig) DelUser(username string) {
	c.mu.Lock()
	if c.config["Users"] != nil {
		delete(c.config["Users"].(map[string]interface{}), username)
		c.flush()
	}
	c.mu.Unlock()
}

func (c *SafeConfig) SetUser(username string, isadmin bool) {
	c.mu.Lock()
	if c.config["Users"] != nil && c.config["Users"].(map[string]interface{})[username] != nil {
		c.config["Users"].(map[string]interface{})[username].(map[string]interface{})["IsAdmin"] = isadmin
		c.flush()
	}
	c.mu.Unlock()

}

func (c *SafeConfig) Json() []byte {
	c.mu.Lock()
	ret, _ := json.Marshal(c.config)
	c.mu.Unlock()
	return ret
}

var loglines int = 100
var logstate []Item
var logstore chan Item = make(chan Item, loglines)
var stat map[string]map[string]int

// var configpath string = "owndns.config"

// var configwriter chan Config
// var config Config
var sconfig SafeConfig

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

	sconfig = *NewSafeConfig(*configpath)
}

// TODO: Replace to direct call
func isAdminExist() bool {
	return sconfig.IsAdminExist()
}

// Validate - Validate DNS request.
func Validate(name string, ip string) bool {
	ret := true
	for _, s := range sconfig.GetRules(ip) {
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

	if sconfig.FindReject(name) {
		ret = false
	} else if ret == true {
		ret = true
	}
	return ret
}

// GetRedirectIP - Get redirect ip address. TODO: Replace to direct call
func GetRedirectIP() (string, error) {
	return sconfig.GetRedirectIP()
}

// GetRejectnames - Pattern list of dns name. TODO: Replace to direct call
func GetRejectnames() []string {
	return sconfig.GetRejectnames()
	//return config.RejectNames
}

// GetRelaydns - IP:PORT list of relay dns. TODO: Replace to direct call
// func GetRelaydns() []string {
// 	return sconfig.GetRelaydns()
// }

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
	if sconfig.IsLogging() {
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
// func ArrAddObject(arr *[]interface{}, val interface{}) {
// 	*arr = append(*arr, &val)
// 	configwriter <- config
// }

// ArrAdd - Add item to array of string.
// func ArrAdd(arr *[]string, val string) {
// 	*arr = append(*arr, val)
// 	configwriter <- config
// }

// ArrDelObject - Del item from array of objects.
// func ArrDelObject(arr *[]interface{}, id int) {
// 	copy((*arr)[id:], (*arr)[id+1:])
// 	(*arr)[len((*arr))-1] = ""
// 	*arr = (*arr)[:len((*arr))-1]
// 	configwriter <- config
// }

// ArrDel - Del item from array of strings.
// func ArrDel(arr *[]string, id int) {
// 	copy((*arr)[id:], (*arr)[id+1:])
// 	(*arr)[len((*arr))-1] = ""
// 	*arr = (*arr)[:len((*arr))-1]
// 	configwriter <- config
// }

// ArrSaveObject - Set item in array of objects.
// func ArrSaveObject(arr *[]interface{}, id int, item interface{}) {
// 	(*arr)[id] = item
// 	configwriter <- config
// }

// ArrSave - Set item in array of strings.
// func ArrSave(arr *[]string, id int, name string) {
// 	(*arr)[id] = name
// 	configwriter <- config
// }

// GetRules - List filtering rules.
// func GetRules(ip string) (ret []map[string]interface{}) {
// 	for _, s := range config.Rules {
// 		if s.(map[string]interface{})["IP"].(string) == ip {
// 			ret = append(ret, s.(map[string]interface{}))
// 		}
// 	}
// 	return ret
// }
