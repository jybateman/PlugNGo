package main

import (
	"fmt"
	"log"
	// "time"
	"net/http"

	"golang.org/x/net/websocket"

	// "github.com/paypal/gatt"
	"github.com/currantlabs/gatt"
)

var svcUUID = gatt.MustParseUUID("fff0")
var cmdUUID = gatt.MustParseUUID("fff3")
var notifUUID = gatt.MustParseUUID("fff4")
var nameUUID = gatt.MustParseUUID("fff6")

var plugs map[string]*Plug

var DefaultClientOptions = []gatt.Option{
	gatt.LnxDeviceID(-1, true),
}

var device gatt.Device

func onStateChanged(d gatt.Device, s gatt.State) {
	log.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		log.Println("scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	log.Printf("Peripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	log.Println("  Local Name        =", a.LocalName)
	log.Println("  TX Power Level    =", a.TxPowerLevel)
	log.Println("  Manufacturer Data =", a.ManufacturerData)
	log.Println("  Service Data      =", a.ServiceData)
	p.Device().Connect(p)
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	log.Printf("Peripheral connected\n")

	services, err := p.DiscoverServices(nil)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	for _, service := range services {
		if (service.UUID().Equal(svcUUID)) {
			log.Printf("Service Found %s\n", service.Name())

			cs, _ := p.DiscoverCharacteristics(nil, service)

			var tmpPlug Plug

			tmpPlug.per = p
			for _, c := range cs {

				if (c.UUID().Equal(notifUUID)) {
					log.Println("Notif Characteristic Found")
					_, err := p.DiscoverDescriptors(nil, c)
					if err != nil {
						log.Println("ERROR:", err)
						return
					}
					if err := p.SetNotifyValue(c, tmpPlug.handleNotification); err != nil {
						log.Println("ERROR:", err)
						return
					}
					tmpPlug.notif = c
				} else if (c.UUID().Equal(cmdUUID)) {
					log.Println("Command Characteristic Found")
					tmpPlug.cmd = c
				} else if (c.UUID().Equal(nameUUID)) {
					log.Println("Name Characteristic Found")
					b, _ := p.ReadCharacteristic(c)
					log.Printf("Name: %q\n", b)
					tmpPlug.name = c
					tmpPlug.Name = string(b)
				}
			}
			if tmpPlug.cmd == nil || tmpPlug.name == nil || tmpPlug.notif == nil {
				return
			}
			tmpPlug.notifChan = make(chan []byte)
			tmpPlug.quit = make(chan bool)
			tmpPlug.ID = p.ID()
			plugs[p.ID()] = &tmpPlug

			// UNCOMMENT MonitorState
			// go plugs[p.ID()].MonitorState()
			// DON'T UNCOMMENT Handler
			// go plugs[p.ID()].Handler()

			// TEST
			// plugs[p.ID()].SetName("HelloWorld")
			// var sch Schedule
			// sch.Name = "hello"
			// sch.StartHour = "1"
			// sch.StartMinute = "30"
			// sch.EndHour = "10"
			// sch.EndMinute = "25"
			// plugs[p.ID()].SetSchedule(sch)
			// time.Sleep(5*time.Second)
			// go plugs[p.ID()].GetSchedule()
			// END TEST

			break
		}
	}
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	_, ok := plugs[p.ID()]
	if ok {
		plugs[p.ID()].quit <- true
		delete(plugs, p.ID())
		p.Device().Connect(p)
	}
	log.Println("Disconnected", p.Name(), err)
}

func Notifytest(c *gatt.Characteristic, b []byte, err error) {
	fmt.Println(b);
}

func initDevice() {
	var err error
	plugs = make(map[string]*Plug)

	device, err = gatt.NewDevice(DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	// Register handlers.
	device.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)
	device.Init(onStateChanged)
}

func main() {
	initDevice()

	initSQL()
	// Handle http
	http.HandleFunc("/", checkSession)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/home", home)
	http.HandleFunc("/plug/", plug)

	http.Handle("/ws", websocket.Handler(handleWS))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.ListenAndServe(":4243", nil)
}
