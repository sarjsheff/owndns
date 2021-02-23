package main

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		res1 := sconfig.Json() //json.Marshal(config)
		w.Write(res1)
	}))

	mux.HandleFunc("/rejects/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUDArrayObject(w, req, "rejects", "RejectNames")
	}))
	mux.HandleFunc("/relaydns/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUDArrayObject(w, req, "relaydns", "RelayDNS")
	}))
	mux.HandleFunc("/rules/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUDArrayObject(w, req, "rules", "Rules")
	}))

	mux.HandleFunc("/users/", checkAuth(func(w http.ResponseWriter, req *http.Request) {
		HTTPCRUD(w, req, "users", func(res map[string]interface{}) {
			item := res["name"].(map[string]interface{})
			if sconfig.GetUser(item["Username"].(string)) != nil {
				fmt.Fprintf(w, `{"ok":false,"error":"User exist."}`)
			} else {
				sconfig.AddUser(item["Username"].(string), fmt.Sprintf("%x", sha256.Sum256([]byte(item["Password"].(string)))), item["IsAdmin"] != nil)
				fmt.Fprintf(w, `{"ok":true}`)
			}
		}, func(res map[string]interface{}) {
			sconfig.DelUser(res["id"].(string))
			fmt.Fprintf(w, `{"ok":true}`)
		}, func(res map[string]interface{}) {

			item := res["name"].(map[string]interface{})
			if sconfig.GetUser(res["id"].(string)) != nil {
				sconfig.SetUser(res["id"].(string), item["IsAdmin"].(bool) == true)
				fmt.Fprintf(w, `{"ok":true}`)
			} else {
				fmt.Fprintf(w, `{"ok":false,"error":"User not found."}`)
			}
		}, func(res map[string]interface{}) {
			res1, err := json.Marshal(sconfig.Get("Users"))
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

	http.ListenAndServe(sconfig.Get("HTTPListen").(string), mux)

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

	if err == nil {
		for key, value := range params {
			sconfig.Set(key, value)
		}
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

	//log.Println(res)

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
func HTTPCRUDArrayObject(w http.ResponseWriter, req *http.Request, name string, arr string) {
	HTTPCRUD(w, req, name, func(res map[string]interface{}) {
		sconfig.AddObject(arr, res["name"])
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		//ArrDelObject(arr, int(res["id"].(float64)))
		sconfig.DelObject(arr, int(res["id"].(float64)))
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		//ArrSaveObject(arr, int(res["id"].(float64)), res["name"])
		sconfig.SetObject(arr, int(res["id"].(float64)), res["name"])
		fmt.Fprintf(w, `{"ok":"true"}`)
	}, func(res map[string]interface{}) {
		res1, err := json.Marshal(sconfig.Get(arr))
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, `{"error":"true"}`)
		} else {
			w.Write(res1)
		}
	})
}

// HTTPCRUDArray - Universal HTTP handler for add, delete and set element of configuration array of string.
// func HTTPCRUDArray(w http.ResponseWriter, req *http.Request, name string, arr *[]string) {
// 	HTTPCRUD(w, req, name, func(res map[string]interface{}) {
// 		ArrAdd(arr, res["name"].(string))
// 		fmt.Fprintf(w, `{"ok":"true"}`)
// 	}, func(res map[string]interface{}) {
// 		ArrDel(arr, int(res["id"].(float64)))
// 		fmt.Fprintf(w, `{"ok":"true"}`)
// 	}, func(res map[string]interface{}) {
// 		ArrSave(arr, int(res["id"].(float64)), res["name"].(string))
// 		fmt.Fprintf(w, `{"ok":"true"}`)
// 	}, func(res map[string]interface{}) {
// 		res1, err := json.Marshal(*arr)
// 		if err != nil {
// 			log.Println(err)
// 			fmt.Fprintf(w, `{"error":"true"}`)
// 		} else {
// 			w.Write(res1)
// 		}
// 	})
// }
