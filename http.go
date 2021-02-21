package main

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func httpserver() {
	gob.Register(HTTPUser{})

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(AssetFile()))

	mux.HandleFunc("/login", HTTPLogin)
	mux.HandleFunc("/logout", HTTPLogout)
	mux.HandleFunc("/user.json", HTTPCurrentUser)

	// mux.HandleFunc("/ses", func(w http.ResponseWriter, req *http.Request) {
	// 	session, err := store.Get(req, "owndns")
	// 	if err == nil {
	// 		log.Println(session)
	// 		log.Println(session.Values["test"])
	// 		if session.Values["test"] != nil {
	// 			session.Values["test"] = session.Values["test"].(int) + 1
	// 		} else {
	// 			session.Values["test"] = 1
	// 		}
	// 	} else {
	// 		log.Println(err)
	// 	}

	// 	err = session.Save(req, w)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// })

	mux.HandleFunc("/version.json", checkAuth(HTTPVersion))
	mux.HandleFunc("/stat.json", checkAuth(HTTPStat))
	mux.HandleFunc("/log.json", checkAuth(HTTPLog))

	mux.HandleFunc("/setvalue", checkAuth(HTTPSetValue))

	mux.HandleFunc("/config/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		res1, err := json.Marshal(config)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"error":"true"}`)
		} else {
			w.Write(res1)
		}
	}))

	mux.HandleFunc("/rejects/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUDArray(w, req, "rejects", &config.RejectNames)
	}))
	mux.HandleFunc("/relaydns/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUDArray(w, req, "relaydns", &config.RelayDNS)
	}))
	mux.HandleFunc("/rules/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUDArrayObject(w, req, "rules", &config.Rules)
	}))

	mux.HandleFunc("/users/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUD(w, req, "users", func(res map[string]interface{}) {

			log.Printf("%T\n", res["name"].(map[string]interface{})["Username"])
			item := res["name"].(map[string]interface{})
			if _, ok := config.Users[item["Username"].(string)]; ok {
				fmt.Fprintf(w, `{"ok":false,"error":"User exist."}`)
			} else {
				if config.Users == nil {
					config.Users = map[string]interface{}{}
				}
				config.Users[item["Username"].(string)] = map[string]interface{}{
					"PasswordHash": fmt.Sprintf("%x", sha256.Sum256([]byte(item["Password"].(string)))),
					"IsAdmin":      item["IsAdmin"] != nil,
				}
				configwriter <- config
				fmt.Fprintf(w, `{"ok":true}`)
			}

		}, func(res map[string]interface{}) {
			delete(config.Users, res["id"].(string))
			configwriter <- config
			fmt.Fprintf(w, `{"ok":true}`)
		}, func(res map[string]interface{}) {
			//ArrSaveObject(arr, int(res["id"].(float64)), res["name"])

			item := res["name"].(map[string]interface{})
			log.Println(item["IsAdmin"])
			if _, ok := config.Users[res["id"].(string)]; ok {
				config.Users[res["id"].(string)] = map[string]interface{}{
					"PasswordHash": config.Users[res["id"].(string)].(map[string]interface{})["PasswordHash"],
					"IsAdmin":      item["IsAdmin"].(bool) == true,
				}
				configwriter <- config
				fmt.Fprintf(w, `{"ok":true}`)
			} else {
				fmt.Fprintf(w, `{"ok":false,"error":"User not found."}`)
			}
		}, func(res map[string]interface{}) {
			res1, err := json.Marshal(config.Users)
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, `{"error":true}`)
			} else {
				w.Write(res1)
			}
		})
	}))

	//req.URL.Path

	mux.HandleFunc("/ws", WSHandler)

	http.ListenAndServe(config.HTTPListen, mux)

}

// HTTPVersion - Version HTTP handler.
func HTTPVersion(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf(`{"version":"%s"}`, BuildTime))
}

// HTTPLog - Last log HTTP handler.
func HTTPLog(w http.ResponseWriter, req *http.Request) {

	res, err := json.Marshal(GetLastLog())
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"error":"true"}`)
	} else {
		w.Write(res)
		//fmt.Fprintf(w, fmt.Sprintf(`{"version":"%s"}`, BuildTime))
	}
}

// HTTPStat - Statistics handler
func HTTPStat(w http.ResponseWriter, req *http.Request) {

	res, err := json.Marshal(GetLogStat())
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"error":"true"}`)
	} else {
		w.Write(res)
		//fmt.Fprintf(w, fmt.Sprintf(`{"version":"%s"}`, BuildTime))
	}
}

// HTTPSetValue - Handler for setting simple value of config.
func HTTPSetValue(w http.ResponseWriter, req *http.Request) {
	params := make(map[string]interface{})
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&params)

	postconfig := config

	if err == nil {
		rr := reflect.ValueOf(&postconfig).Elem()
		for key, value := range params {
			v := rr.FieldByName(key)
			if v.IsValid() {
				v.Set(reflect.ValueOf(value))
			}
		}

		configwriter <- postconfig
		fmt.Fprintf(w, `{"ok":"true"}`)
	} else {
		fmt.Fprintf(w, `{"error":"true"}`)
	}

}

// CRUDFn - Stub function.
type CRUDFn func(map[string]interface{})

// HTTPCRUD - Universal HTTP handler for add, delete and set.
func HTTPCRUD(w http.ResponseWriter, req *http.Request, ctxname string, add CRUDFn, del CRUDFn, save CRUDFn, def CRUDFn) {
	//fmt.Fprintf(w, fmt.Sprintf(`%s`, req.URL.Path))

	res := make(map[string]interface{})
	dec := json.NewDecoder(req.Body)

	err := dec.Decode(&res)

	log.Println(res)

	if err != nil {
		def(nil)
	} else {

		switch req.URL.Path {
		case fmt.Sprintf("/%s/add", ctxname):
			if add != nil {
				add(res)
			}
			break
		case fmt.Sprintf("/%s/del", ctxname):
			if del != nil {
				del(res)
			}
			break
		case fmt.Sprintf("/%s/save", ctxname):
			if save != nil {
				save(res)
			}
			break
		default:
			log.Println("Default handler", req.URL.Path, ctxname)
			def(res)
			break
		}
	}
}

// HTTPCRUDArrayObject - Universal HTTP handler for add, delete and set element of configuration array of objects.
func HTTPCRUDArrayObject(w http.ResponseWriter, req *http.Request, name string, arr *[]interface{}) {
	HTTPCRUD(w, req, name, func(res map[string]interface{}) {
		ArrAddObject(arr, res["name"])
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		ArrDelObject(arr, int(res["id"].(float64)))
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		ArrSaveObject(arr, int(res["id"].(float64)), res["name"])
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		res1, err := json.Marshal(arr)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"error":"true"}`)
		} else {
			w.Write(res1)
		}
	})
}

// HTTPCRUDArray - Universal HTTP handler for add, delete and set element of configuration array of string.
func HTTPCRUDArray(w http.ResponseWriter, req *http.Request, name string, arr *[]string) {
	HTTPCRUD(w, req, name, func(res map[string]interface{}) {
		ArrAdd(arr, res["name"].(string))
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		ArrDel(arr, int(res["id"].(float64)))
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		ArrSave(arr, int(res["id"].(float64)), res["name"].(string))
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		res1, err := json.Marshal(*arr)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"error":"true"}`)
		} else {
			w.Write(res1)
		}
	})
}
