package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	// "github.com/sunmyinf/workplacehub/handler"
)

const textMessageType int = 1

var upgrader = websocket.Upgrader{}

func main() {
	msgChan := make(chan []byte, 4)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		bufBody := new(bytes.Buffer)
		bufBody.ReadFrom(req.Body)

		msgChan <- bufBody.Bytes()
	})

	// websocket通信側のendpoint
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		defer c.Close()

		for msg := range msgChan {
			if err = c.WriteMessage(textMessageType, msg); err != nil {
				log.Println("write: %v", err)
				break
			}
		}
	})

	http.ListenAndServe("localhost:8010", nil)
}
