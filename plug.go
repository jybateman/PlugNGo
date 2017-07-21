package main

import (
	"log"
	"bytes"
	"encoding/hex"
	"encoding/binary"
	
	"github.com/paypal/gatt"
)

type Plug struct {
	ID string
	Name string
	State bool
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

func (pl *Plug) Status() {
	b := CreateMessage([]byte{0x04, 0x00, 0x00, 0x00, 0x00})
	pl.SendMessage(b)
}

func (pl *Plug) SetName(name string) {
	
}

func (pl *Plug) Test() {
	b := CreateMessage([]byte{0x0a, 0x00, 0x00, 0x00, 0x00})
	pl.SendMessage(b)
}

// HANDLERS

func (pl *Plug) HandleStatus(data []byte) {
	if len(data) < 11 {
		log.Println("ERROR: Couldn't retrieve power usage")
		return
	}
	power := binary.BigEndian.Uint32(data[6:])
	voltage := data[10]
	log.Println("power:", power, "Voltage:", voltage)
}

func (pl *Plug) HandleDayHist(data []byte) {
	if len(data) < 4 {
		log.Println("ERROR: Couldn't retrieve power history")
		return
	}

	for i := 2; i + 2 <= len(data); i = i + 2 {
		log.Println(binary.BigEndian.Uint16(data[i:i+2]))
	}
}

func (pl *Plug) Handler() {
	defer log.Println("Good Night")
	for {
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
				pl.HandleStatus(data)
				break
			case bytes.HasPrefix(data, PowerDayHistory):
				log.Println("Power History for last 24h Notification")
				pl.HandleDayHist(data)
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
}
