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

func (pl *Plug) On() {
	b := CreateMessage([]byte{0x03, 0x00, 0x01, 0x00, 0x00})
}
