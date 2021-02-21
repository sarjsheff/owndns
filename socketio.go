package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Op struct {
	Op      string
	Payload interface{}
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var op *Op
		err = json.Unmarshal(message, &op)
		if err == nil {
			switch op.Op {
			case "Version":
				WSVersion(c)
				break
			default:
				WSVersion(c)
				break
			}
		}
	}
}

func WSVersion(c *websocket.Conn) error {
	r := []byte(fmt.Sprintf(`{"version":"%s"}`, BuildTime))
	err := c.WriteMessage(len(r), r)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}
