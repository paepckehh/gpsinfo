package main

import (
	"os"
	"sync"

	"paepcke.de/gpstime/gpsfeed"
)

const (
	_lf            = "\n"
	_defaultDevice = "/dev/gps0"
)

var (
	bg          sync.WaitGroup
	chanDisplay = make(chan string, 10)
)

// gpsfeed is an simple example app for gpsfeed [cat /dev/gps0, with some flowcontrol]
func main() {
	// read gps donge device from cli
	dev := &gpsfeed.GpsDevice{}
	dev.FileIO = _defaultDevice
	for i := 1; i < len(os.Args); i++ {
		dev.FileIO = gpsfeed.GetDeviceName(os.Args[i], dev.FileIO)
	}

	bg.Add(1)
	go display()
	go feeder(dev)
	bg.Wait()
}

// feeder reads line-by-line from gps device via bufio
func feeder(dev *gpsfeed.GpsDevice) {
	dev.Open()
	for dev.Feed.Scan() {
		s := dev.Feed.Text()
		dev.Responsive.Store(true)
		l := len(s)
		if l > 15 && l < 256 {
			if s[0] == '$' {
				dev.DataValid.Store(true)
				chanDisplay <- s + " --> checksum: " + gpsfeed.CheckSumTag(s)
			}
		}
	}
	chanDisplay <- "[gpsfeed] [device closed] [" + dev.FileIO + "]"
	close(chanDisplay)
}

// display engine
func display() {
	for s := range chanDisplay {
		os.Stdout.Write([]byte(s + _lf))
	}
	bg.Done()
}
