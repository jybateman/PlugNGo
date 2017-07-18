package main

import (
	"fmt"
	
	"github.com/paypal/gatt"
)

var (
	StartMessage = []byte{0x0f}
	EndMessage = []byte{0xff, 0xff}
)

type Plug struct {
	per gatt.Peripheral
	cmd *gatt.Characteristic
	name *gatt.Characteristic
	notif *gatt.Characteristic
}

func (pl *Plug) SendMessage(b []byte) {
	for i := 0; i < len(b); i += 20 {
		end := i + 20
		if end > len(b) {
			end = len(b)
		}
		err := pl.per.WriteCharacteristic(pl.cmd, b[i:end], true)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
	
func (pl *Plug) On() {
	b := CreateMessage([]byte{0x03, 0x00, 0x01, 0x00, 0x00})
	pl.SendMessage(b)
}
