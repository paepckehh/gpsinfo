// package main
package main

// import
import (
	"os"

	"paepcke.de/gpsinfo"
	"paepcke.de/gpsinfo/gpsfeed"
)

// const defaults
const _defaultDevice = "/dev/gps0"

// main ...
func main() {
	port := _defaultDevice
	for i := 1; i < len(os.Args); i++ {
		port = gpsfeed.GetDeviceName(os.Args[i], port)
	}
	gpsinfo.Debug(port)
}
