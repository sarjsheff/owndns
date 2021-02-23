package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))

// HTTPUser - Object for session and json response.
type HTTPUser struct {
	Username string
	IsAdmin  bool
}

func checkAuth(handler func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "owndns")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		if session.Values["user"] != nil {
			handler(w, req)
		} else {
			http.Error(w, `{"error":true,"errortext":"Forbidden."}`, http.StatusForbidden)
		}
		err := session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func HTTPCurrentUser(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "owndns")
	if session.Values["user"] != nil {
		ret, _ := json.Marshal(map[string]interface{}{
			"ok":   true,
			"user": session.Values["user"],
		})
		w.Write(ret)
	} else {
		fmt.Fprintf(w, `{"ok":false}`)
	}
}

// HTTPLogout - Logout handler.
func HTTPLogout(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "owndns")
	delete(session.Values, "user")
	err := session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `{"ok":true}`)
}

// HTTPLogin - Login handler.
func HTTPLogin(w http.ResponseWriter, req *http.Request) {
	res := make(map[string]string)
	dec := json.NewDecoder(req.Body)

	err := dec.Decode(&res)

	if err != nil {
		fmt.Fprintf(w, `{"ok":false}`)
	} else {
		u := sconfig.GetUser(res["username"]) //config.Users[res["username"]]
		if u != nil && u["PasswordHash"].(string) == res["password"] {
			session, _ := store.Get(req, "owndns")
			user := HTTPUser{res["username"], u["IsAdmin"].(bool)}

			ret, _ := json.Marshal(map[string]interface{}{
				"ok":   true,
				"user": user,
			})
			session.Values["user"] = user

			err := session.Save(req, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(ret)
		} else {
			fmt.Fprintf(w, `{"ok":false}`)
		}

	}
}
