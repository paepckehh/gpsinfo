// package gpsinfo ...
package gpsinfo

// import
import (
	"time"

	"paepcke.de/gpsfeed"
)

// const
const (
	ecoSleep     = time.Duration(90 * time.Millisecond)
	displaySleep = time.Duration(ecoSleep / 2)
	maxDiff      = time.Duration(2 * time.Second)
)

// debug ...
func debug(dev *gpsfeed.GpsDevice) {
	// setup
	channelOut := make(chan string, 10)
	channelGpsFrames := make(chan string, 50)

	for {
		// spin up background bufio sentence fetcher/filter process
		go func() {
			dev.Open()
			for dev.Feed.Scan() {
				sentence := dev.Feed.Text()
				dev.Responsive.Store(true)
				l := len(sentence)
				if l > 15 && l < 256 {
					if sentence[0] == '$' {
						dev.DataValid.Store(true)
						channelGpsFrames <- sentence
					}
				}
				time.Sleep(ecoSleep)
			}
			close(channelGpsFrames)
		}()

		// spin up background Display outout handler
		go func() {
			for s := range channelOut {
				out(s)
			}
		}()

		// builder loop
		build(dev, channelGpsFrames, channelOut, displaySleep)
		close(channelGpsFrames)
		close(channelOut)
		dev.Close()
	}
}
