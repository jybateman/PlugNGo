package main

import (
	"fmt"
	"log"
	
	"github.com/paypal/gatt"
)

var svcUUID = gatt.MustParseUUID("fff0")
var cmdUUID = gatt.MustParseUUID("fff3")
var notifUUID = gatt.MustParseUUID("fff4")
var nameUUID = gatt.MustParseUUID("fff6")

var plugs map[string]*Plug

var DefaultClientOptions = []gatt.Option{
	gatt.LnxDeviceID(-1, true),
}

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
					}
					if err := p.SetNotifyValue(c, tmpPlug.handleNotification); err != nil {
						log.Println("ERROR:", err)
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
				}
			}
			tmpPlug.notifChan = make(chan []byte)
			plugs[p.ID()] = &tmpPlug
			go plugs[p.ID()].Handler()
			// plugs[p.ID()].On()
			plugs[p.ID()].Status()
			break
		}
	}
}

func Notifytest(c *gatt.Characteristic, b []byte, err error) {
	fmt.Println(b);
}

func main() {
	plugs = make(map[string]*Plug)
	
	d, err := gatt.NewDevice(DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		// gatt.PeripheralDisconnected(onPeriphDisconnected),
	)
	d.Init(onStateChanged)

	select {}	
}
