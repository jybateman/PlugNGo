package main

import (
	"log"
	"time"
	"bytes"
	"strconv"
	"encoding/hex"
	"encoding/binary"

	// "github.com/paypal/gatt"
	"github.com/currantlabs/gatt"
)

type Schedule struct {
	Name string
	StartHour, StartMinute string
	EndHour, EndMinute string
	Flag byte
}

type Plug struct {
	ID string
	Name string
	State byte
	Schedules []Schedule
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
	if bytes.HasPrefix(b, StartMessage) {
		message = b
	} else {
		message = append(message, b...)
	}
	if bytes.HasSuffix(b, EndMessage) {
		log.Println("Sent notification to handler", message)
		pl.notifChan <- message
		message = []byte{0x00}
	}
}

func (pl *Plug) SendMessage(b []byte) {
	log.Println("Going to write message:", hex.EncodeToString(b))

	// err := pl.per.WriteCharacteristic(pl.cmd, b, true)
	// log.Println("Wrote block of message length:", len(b))
	// if err != nil {
	//	log.Println("ERROR:", err)
	// }
	// pl.HandlerNotification()

	for i, end := 0, 0; i < len(b); i = end {
		end = i + 20
		if end >= len(b) {
			end = len(b)
		}
		err := pl.per.WriteCharacteristic(pl.cmd, b[i:end], true)
		log.Println("Wrote block of message length:", len(b[i:end]), hex.EncodeToString(b[i:end]))
		if err != nil {
			log.Println("ERROR:", err)
		}
	}
	log.Println("Wrote all message length:", len(b))
	pl.HandlerNotification()
}

// Plug COMMANDS

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
	var bName []byte

	log.Println("Creating request message")
	mes := []byte{0x02, 0x00}
	if len(name) > 20 {
		bName = []byte(name[:20])
	} else {
		bName = []byte(name)
		for len(bName) < 20 {
			bName = append(bName, 0x00)
		}
	}
	mes = append(mes, bName...)
	b := CreateMessage(mes)
	log.Println("Sending SetName request")
	pl.SendMessage(b)
	log.Println("Sent SetName request")
	b, _ = pl.per.ReadCharacteristic(pl.name)
	pl.Name = string(b)
}

func (pl *Plug) GetSchedule() {
	log.Println("Creating request message")
	b := CreateMessage([]byte{0x07, 0x00, 0x00, 0x00, 0x00})
	log.Println("Sending GetSchedule request")
	pl.SendMessage(b)
	log.Println("Sent GetSchedule request")
}

func (pl *Plug) SetSchedule(sched Schedule) {
	var bName []byte

	pl.GetSchedule()
	log.Println("Creating request message")
	mes := []byte{0x06, 0x00, 0x01}
	if len(sched.Name) > 16 {
		bName = []byte(sched.Name[:16])
	} else {
		bName = []byte(sched.Name)
		for len(bName) < 16 {
			bName = append(bName, 0x00)
		}
	}
	mes = append(mes, bName...)
	mes = append(mes, 0x80)
	bStartHour := -1
	bStartMinute := -1
	// TODO: if Atoi fail send message to client that SetSchedule has failed and don't send notification
	if len(sched.StartHour) > 0 {
		bStartHour, _ = strconv.Atoi(sched.StartHour)
		bStartMinute, _ = strconv.Atoi(sched.StartMinute)
	}

	bEndHour := -1
	bEndMinute := -1
	// TODO: if Atoi fail send message to client that SetSchedule has failed and don't send notification
	if len(sched.EndHour) > 0 {
		bEndHour, _ = strconv.Atoi(sched.EndHour)
		bEndMinute, _ = strconv.Atoi(sched.EndMinute)
	}

	mes = append(mes, byte(bStartHour), byte(bStartMinute), byte(bEndHour), byte(bEndMinute))
	for len(mes) < 112 {
		mes = append(mes, 0x00)
	}
	b := CreateMessage(mes)
	log.Println("Sending SetSchedule request", b)
	pl.SendMessage(b)
	log.Println("Sent SetSchedule request")
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
	storeDatum(pl.ID, power, voltage)
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

func (pl *Plug) HandleSchedule(data []byte) {
	if len(data) < 4 {
		log.Println("ERROR: Couldn't retrieve schedule")
		return
	}
	for offset := 4; offset + 21 < len(data); offset += 22 {
		var tmpSched Schedule
		tmpSched.Name = string(data[offset+1:offset+17])
		tmpSched.StartHour = strconv.Itoa(int(data[offset+18]))
		tmpSched.StartMinute = strconv.Itoa(int(data[offset+19]))
		tmpSched.EndHour = strconv.Itoa(int(data[offset+20]))
		tmpSched.EndMinute = strconv.Itoa(int(data[offset+21]))
		tmpSched.Flag = data[offset+17]
		pl.Schedules = append(pl.Schedules, tmpSched)
	}
	log.Println("Received schedule:", pl.Schedules)
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
			pl.HandleSchedule(data)
			break
		case bytes.HasPrefix(data, Notif):
			log.Println("Notification")
			break
		default:
			log.Println("Unknown Notification")
		}
		break
	case <- time.After(time.Second * 10):
		log.Println("ERROR: Didn't receive notification")
	}
}
