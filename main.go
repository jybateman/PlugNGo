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

var plugs map[string]Plug

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
						log.Println(err)
					}
					f := func(c *gatt.Characteristic, b []byte, err error) {
						fmt.Printf("notified: % X | %q\n", b, b)
					}
					if err := p.SetNotifyValue(c, f); err != nil {
						log.Println("ERROR:", err)
					}
					tmpPlug.notif = c
				}
				
				if (c.UUID().Equal(cmdUUID)) {
					log.Println("Command Characteristic Found")
					tmpPlug.cmd = c

				}
				if (c.UUID().Equal(nameUUID)) {
					log.Println("Name Characteristic Found")
					b, _ := p.ReadCharacteristic(c)
					log.Printf("Name: %q\n", b)
					tmpPlug.name = c
				}
			}
			plugs[p.ID()] = tmpPlug
			break
		}
	}
}

func Notifytest(c *gatt.Characteristic, b []byte, err error) {
	fmt.Println(b);
}

func main() {
	plugs = make(map[string]Plug)
	// hexMes, _ := hex.DecodeString("0300010000")
	// m := hex.EncodeToString(CreateMessage(hexMes))
	// fmt.Println("COMPARE MESSAGES:", m, "0f06030001000005ffff")
	// fmt.Println("COMPARE LENGTH:", len(m), len("0f06030001000005ffff"))
	// fmt.Println("IS EQUAL:", m == "0f06030001000005ffff")
	
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
