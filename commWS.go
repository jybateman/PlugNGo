package main

import (
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func ChangeState(req []string, ws *websocket.Conn) {
	pl, ok := plugs[req[1]]
	if len(req) > 1 && ok {
		if pl.State > 0 {
			go pl.Off()
		} else {
			go pl.On()
		}
		time.Sleep(time.Second * 2)
		ws.Write([]byte(implodeRequest([]string{"0", req[1], string(pl.State+48)})))
		log.Println("Message sent", implodeRequest([]string{"0", req[1], string(pl.State+48)}))
	}
}

func Status(req []string, ws *websocket.Conn, quit chan bool) {
	_, ok := plugs[req[1]]

	// if len(req) > 3 && ok {

	// } else
	if len(req) > 1 && ok {
		for {
			select {
			case <- quit:
				return
			case <- time.After(time.Second * 2):
				log.Println("Got:", getData(req[1]))
				ws.Write([]byte(implodeRequest([]string{"1", getData(req[1])})))
			}
		}
	}
}

func handleWS(ws *websocket.Conn) {
	quit := make(chan bool)
	buf := make([]byte, 512)
	if !isSession(ws.Request()) {
		return
	}
	for {
		_, err := ws.Read(buf)
		if err != nil {
			log.Println("ERROR:", err)
			ws.Close()
			select {
			case quit <- true:
			default:
			}
			return
		}
		arr := explodeRequest(string(buf))
		switch arr[0] {
		case "0":
			log.Println("WS received change state request")
			go ChangeState(arr, ws)
		case "1":
			log.Println("WS received status request", arr)
			go Status(arr, ws, quit)
		case "2":
			log.Println("WS received change name  request", arr)
			go plugs[arr[1]].SetName(arr[2])
		}
	}
}
