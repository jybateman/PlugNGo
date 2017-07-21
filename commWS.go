package main

import (
	"log"
	"golang.org/x/net/websocket"
)

func handleWS(ws *websocket.Conn) {
	buf := make([]byte, 512)
	if !isSession(ws.Request()) {
		return
	}
	for {
		_, err := ws.Read(buf)
		// plugs.
		if err != nil {
			log.Println("ERROR:", err)
			ws.Close()
			return
		}
		arr := explodeRequest(string(buf))
		if plugs[arr[1]].State > 0 {
			go plugs[arr[1]].Off()
		} else {
			go plugs[arr[1]].On()
		}
		ws.Write([]byte(implodeRequest([]string{"0", "0"})))
		log.Println("Message sent", implodeRequest([]string{"0", "0"}))
	}
}
