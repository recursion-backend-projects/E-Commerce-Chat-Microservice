package websocket

import (
	"log"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write:", err)
			break
		} else {
			log.Printf("send: %s", message)
		}
	}
}