package main

import (
	"bytes"
	"flag"
	"json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sunmyinf/workplacehub/decode"
)

const textMessageType int = 1

var upgrader = websocket.Upgrader{}

func main() {
	port := flag.String("port", "8010", "port number")
	chanNum := flag.Int("chanNum", 4, "number of buffer of channel between callback and echo")
	flag.Parse()

	msgChan := make(chan []byte, *chanNum)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		bufBody := new(bytes.Buffer)
		if _, err := bufBody.ReadFrom(req.Body); err != nil {
			log.Printf("callback err: %v\n", err)
			return
		}

		group := decode.Group{}
		if err := json.Unmarshal(bufBody, &group); err != nil {
			log.Printf("json unmarshal error: %v")
			return
		}

		// send post's or comment's message
		msgChan <- bufBody.Bytes()
	})

	// endpoint for web socket connection
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

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Panicf("Failed to launch server: %v", err)
	}
}
