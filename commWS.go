package main

import (
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func ChangeState(req []string, ws *websocket.Conn) {
	if len(req) > 1 {
		if plugs[req[1]].State > 0 {
			go plugs[req[1]].Off()
		} else {
			go plugs[req[1]].On()
		}
		time.Sleep(time.Second * 2)
		ws.Write([]byte(implodeRequest([]string{"0", req[1], string(plugs[req[1]].State+48)})))
		log.Println("Message sent", implodeRequest([]string{"0", req[1], string(plugs[req[1]].State+48)}))
	}
}

func handleWS(ws *websocket.Conn) {
	buf := make([]byte, 512)
	if !isSession(ws.Request()) {
		return
	}
	for {
		_, err := ws.Read(buf)
		if err != nil {
			log.Println("ERROR:", err)
			ws.Close()
			return
		}
		arr := explodeRequest(string(buf))
		switch arr[0] {
		case "0":
			log.Println("WS received change state request")
			go ChangeState(arr, ws)
		}
	}
}
