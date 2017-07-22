package main

import (
	"log"
	"time"
	"bytes"
	"encoding/hex"
	"encoding/binary"

	// "github.com/paypal/gatt"
	"github.com/currantlabs/gatt"
)

type Plug struct {
	ID string
	Name string
	State byte
	per gatt.Peripheral
	cmd *gatt.Characteristic
	name *gatt.Characteristic
	notif *gatt.Characteristic
	notifChan chan []byte
	quit chan bool
}

var message []byte


func (pl *Plug) MonitorState() {
	for {
		select {
		case <- pl.quit:
			return
		case <- time.After(time.Second * 2):
			pl.Status()
		}
	}
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
		log.Println("Going to write message")
		err := pl.per.WriteCharacteristic(pl.cmd, b[i:end], true)
		if err != nil {
			log.Println("ERROR:", err)
		}
		log.Println("Wrote message")
	}
	pl.HandlerNotification()
}

func (pl *Plug) On() {
	log.Println("Creating request message")
	b := CreateMessage([]byte{0x03, 0x00, 0x01, 0x00, 0x00})
	log.Println("Sending On request")
	pl.SendMessage(b)
	log.Println("Sent On request")
	pl.State = 1
}

func (pl *Plug) Off() {
	log.Println("Creating request message")
	b := CreateMessage([]byte{0x03, 0x00, 0x00, 0x00, 0x00})
	log.Println("Sending Off request")
	pl.SendMessage(b)
	log.Println("Sent Off request")
	pl.State = 0
}

func (pl *Plug) Status() {
	log.Println("Creating request message")
	b := CreateMessage([]byte{0x04, 0x00, 0x00, 0x00, 0x00})
	log.Println("Sending Status request")
	pl.SendMessage(b)
	log.Println("Sent Status request")
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
	pl.State = data[4]
	power := binary.BigEndian.Uint32(data[6:])
	voltage := data[10]
	log.Println("Full notification:", data)
	log.Println("power:", power, "Voltage:", voltage, "State:", data[4])
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


func (pl *Plug) HandlerNotification() {
	defer log.Println("Good Night")
	log.Println("Good Morning")
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
	case <- time.After(time.Second * 10):
		log.Println("ERROR: Didn't receive notification")
	}
}
