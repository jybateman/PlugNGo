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
		arr := explodeString(string(buf))
		go plugs[arr[1]].Status()
	}
}
