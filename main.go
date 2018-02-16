package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sunmyinf/workplacehub/handler"
)

const textMessageType int = 1

var upgrader = websocket.Upgrader{}

func main() {
	msgChan := make(chan string, 4)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Reqest) {

	})

	// websocket通信側のAPI
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		defer c.Close()

		for msg := range msgChan {
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write: %v", err)
				break
			}
		}
	})
}
