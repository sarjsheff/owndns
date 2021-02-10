package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func httpserver() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(AssetFile()))

	mux.HandleFunc("/version.json", HTTP_Version)
	mux.HandleFunc("/stat.json", HTTP_Stat)
	mux.HandleFunc("/log.json", HTTP_Log)

	mux.HandleFunc("/config/", func(w http.ResponseWriter, req *http.Request) {

		def := func() {
			res1, err := json.Marshal(config)
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, `{"error":"true"}`)
			} else {
				w.Write(res1)
			}
		}

		postconfig := Config{DNSListen: ":53", HttpListen: ":8081"}
		dec := json.NewDecoder(req.Body)
		err := dec.Decode(&postconfig)

		if err != nil {
			def()
		} else {
			switch req.URL.Path {
			case "/config/save":
				configwriter <- postconfig
				fmt.Fprintf(w, `{"ok":"true"}`)
				break
			default:
				def()
				break
			}
		}
	})

	mux.HandleFunc("/rejects/", func(w http.ResponseWriter, req *http.Request) {
		HTTP_CRUD_Array(w, req, "rejects", &config.RejectNames)
	})
	mux.HandleFunc("/relaydns/", func(w http.ResponseWriter, req *http.Request) {
		HTTP_CRUD_Array(w, req, "relaydns", &config.RelayDns)
	})
	//req.URL.Path

	http.ListenAndServe(config.HttpListen, mux)

}

func HTTP_Version(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf(`{"version":"%s"}`, BuildTime))
}

func HTTP_Log(w http.ResponseWriter, req *http.Request) {

	res, err := json.Marshal(GetLastLog())
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"error":"true"}`)
	} else {
		w.Write(res)
		//fmt.Fprintf(w, fmt.Sprintf(`{"version":"%s"}`, BuildTime))
	}
}

func HTTP_Stat(w http.ResponseWriter, req *http.Request) {

	res, err := json.Marshal(GetLogStat())
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"error":"true"}`)
	} else {
		w.Write(res)
		//fmt.Fprintf(w, fmt.Sprintf(`{"version":"%s"}`, BuildTime))
	}
}

type CRUDFn func(map[string]interface{})

func HTTP_CRUD(w http.ResponseWriter, req *http.Request, ctxname string, add CRUDFn, del CRUDFn, save CRUDFn, def CRUDFn) {
	//fmt.Fprintf(w, fmt.Sprintf(`%s`, req.URL.Path))

	res := make(map[string]interface{})
	dec := json.NewDecoder(req.Body)

	err := dec.Decode(&res)

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

func HTTP_CRUD_Array(w http.ResponseWriter, req *http.Request, name string, arr *[]string) {
	HTTP_CRUD(w, req, name, func(res map[string]interface{}) {
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
