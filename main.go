package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const textMessageType int = 1

var addr = flag.String("address", "localhost:8010", "setting of server address")
var chanNum = flag.Int("chanNum", 4, "number of buffer of channel between callback and echo")
var upgrader = websocket.Upgrader{}

func main() {
	flag.Parse()
	msgChan := make(chan []byte, *chanNum)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		bufBody := new(bytes.Buffer)
		if _, err := bufBody.ReadFrom(req.Body); err != nil {
			log.Printf("callback err: %v\n", err)
			return
		}

		msgChan <- bufBody.Bytes()
	})

	// websocket通信側のendpoint
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Printf("upgrade: %v\n", err)
			return
		}
		defer c.Close()

		for msg := range msgChan {
			if err = c.WriteMessage(textMessageType, msg); err != nil {
				log.Printf("write: %v\n", err)
				break
			}
		}
	})

	http.ListenAndServe(*addr, nil)
}
