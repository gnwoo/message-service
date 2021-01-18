package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)
var addr = flag.String("addr", "localhost:8080", "http service address")
func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/im"}
	log.Printf("connecting to %s", u.String())
	header:= http.Header{}
	header.Set("whoimi","24601")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	//defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()


	msgBody := `{"sender":"%v", "channel":"messgae", "session":"24602", "body":"this is body"}`
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Println(t.String())
			err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(msgBody,"24601")))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
