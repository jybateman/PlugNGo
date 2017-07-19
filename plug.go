package main

import (
	"log"
	"bytes"
	"encoding/hex"
	
	"github.com/paypal/gatt"
)

type Plug struct {
	per gatt.Peripheral
	cmd *gatt.Characteristic
	name *gatt.Characteristic
	notif *gatt.Characteristic
	notifChan chan []byte
}

func (pl *Plug) handleNotification(c *gatt.Characteristic, b []byte, err error) {
	log.Println("Receive notification:", hex.EncodeToString(b))
	pl.notifChan <- b
}

func (pl *Plug) SendMessage(b []byte) {
	for i := 0; i < len(b); i += 20 {
		end := i + 20
		if end > len(b) {
			end = len(b)
		}
		err := pl.per.WriteCharacteristic(pl.cmd, b[i:end], true)
		if err != nil {
			log.Println("ERROR:", err)
		}
	}
}
	
func (pl *Plug) On() {
	b := CreateMessage([]byte{0x03, 0x00, 0x01, 0x00, 0x00})
	pl.SendMessage(b)
}

func (pl *Plug) Off() {
	b := CreateMessage([]byte{0x03, 0x00, 0x00, 0x00, 0x00})
	pl.SendMessage(b)
}

func (pl *Plug) Handler() {
	select {
	case data := <- pl.notifChan:
		switch {
		case bytes.HasPrefix(data, TimeConfirm):
			log.Println("Set Time Notification")
			break
		case bytes.HasPrefix(data, NameConfirm):
			log.Println("Set Name Notification")
			break
		case bytes.HasPrefix(data, StateConfirm):
			log.Println("State Change Notification")
			break
		case bytes.HasPrefix(data, StatePowerNotif):
			log.Println("State and Power Notification")
			break
		case bytes.HasPrefix(data, PowerDayHistory):
			log.Println("Power History for last 24h Notification")
			break
		case bytes.HasPrefix(data, PowerPerDay):
			log.Println("Power history kWh/day Notification")
			break
		case bytes.HasPrefix(data, ScheduleNotif):
			log.Println("Schedules Notification")
			break
		case bytes.HasPrefix(data, Notif):
			log.Println("Notification")
			break
		default:
			log.Println("Unknown Notification")
		}
	}
}
