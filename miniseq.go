package main

import (
	"fmt"
	"time"

	"github.com/rakyll/portmidi"
)

var buttons = map[int64]int64{}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func scan(out, instru *portmidi.Stream) {
	//col := []int{0, 16, 32, 48, 64, 80, 96, 112}
	notes := []int{0, 2, 4, 5, 7, 9, 11, 12}
	for {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				buttons[int64((j*16)+i)] += 27
				out.WriteShort(0x90, int64((j*16)+i), buttons[int64((j*16)+i)])
				if buttons[int64((j*16)+i)] == 127 {
					instru.WriteShort(0x90, int64(60+(notes[(7-j)])), 100)
				}
			}
			time.Sleep(300 * time.Millisecond)
			for j := 0; j < 8; j++ {
				buttons[int64((j*16)+i)] -= 27
				out.WriteShort(0x90, int64((j*16)+i), buttons[int64((j*16)+i)])
				if buttons[int64((j*16)+i)] == 100 {
					instru.WriteShort(0x80, int64(60+(notes[(7-j)])), 0)
				}
			}
		}
	}
}

func debug(debugport *portmidi.Stream) {
	ch := debugport.Listen()
	for {
		events := <-ch
		fmt.Println(events)
	}
}

func main() {
	portmidi.Initialize()
	out, err := portmidi.NewOutputStream(portmidi.DeviceID(32), 1024, 0)
	check(err)
	in, err := portmidi.NewInputStream(portmidi.DeviceID(33), 1024)
	check(err)
	instru, err := portmidi.NewOutputStream(portmidi.DeviceID(0), 1024, 0)
	check(err)
	instru.SetChannelMask(1)
	//debugport, err := portmidi.NewInputStream(portmidi.DeviceId(1), 1024)
	//check(err)*/
	for i := 0; i < portmidi.CountDevices(); i++ {
		fmt.Println(i, portmidi.Info(portmidi.DeviceID(i)))
	}
	//out.WriteSysExBytes(portmidi.Time(), []byte{0xB0, 0x00, 0x00})
	//out.WriteSysExBytes(portmidi.Time(), []byte{0x90, 0x60, 0x0F})
	go scan(out, instru)
	//go debug(debugport)
	ch := in.Listen()
	for {
		event := <-ch
		if event.Data2 == 127 {
			if buttons[event.Data1] == 0 {
				out.WriteShort(0x90, event.Data1, 100)
				buttons[event.Data1] = 100
			} else {
				out.WriteShort(0x90, event.Data1, 0)
				buttons[event.Data1] = 0
			}
		}
	}
	//in.Close()
}
