package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sunmyinf/go-workplace/decode"
	"github.com/sunmyinf/go-workplace/webhook"
)

const textMessageType int = 1

var upgrader = websocket.Upgrader{}

func main() {
	port := flag.String("port", "8010", "port number")
	chanNum := flag.Int("chanNum", 4, "number of buffer of channel between callback and echo")
	flag.Parse()

	msgChan := make(chan string, *chanNum)

	appSecret := os.Getenv("WORKPLACE_APP_SECRET")
	accessToken := os.Getenv("WORKPLACE_ACCESS_TOKEN")
	verifyToken := os.Getenv("WORKPLACE_VERIFY_TOKEN")

	ws := webhook.NewServer(appSecret, accessToken, verifyToken)
	ws.HandleObjectFunc(decode.ObjectGroup, func(rb decode.RequestBody) error {
		// send post's or comment's message
		msgChan <- string(rb.Data[0].Object)
		return nil
	})

	// endpoint for web socket connection
	ws.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Printf("upgrade: %v\n", err)
			return
		}
		defer c.Close()

		for msg := range msgChan {
			if err = c.WriteMessage(textMessageType, []byte(msg)); err != nil {
				log.Printf("write: %v\n", err)
				break
			}
		}
	})

	if err := ws.ListenAndServe(":" + *port); err != nil {
		log.Panicf("failed to launch server: %v", err)
	}
}
