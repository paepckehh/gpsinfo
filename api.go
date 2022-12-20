// package gpsinfo ...
package gpsinfo

// import
import (
	"time"

	"paepcke.de/gpstime"
	"paepcke.de/gpstime/gpsfeed"
)

//
// SIMPLE API
//

// Debug ...
func Debug(device string) { debug(&gpsfeed.GpsDevice{FileIO: device}) }

// GetTime ...
func GetTime(device string) time.Time { return gpstime.GetTime(device) }

// GetLocation ...
func GetLocation(device string) gpstime.Coord {
	return gpstime.GetLocation(device)
}

//
// GENERIC BACKEND
//

// DebugD ...
func DebugD(dev *gpsfeed.GpsDevice) { debug(dev) }

// GetTimeD ...
func GetTimeD(dev *gpsfeed.GpsDevice) time.Time { return gpstime.GetTimeD(dev) }

// GetLocationD ...
func GetLocationD(dev *gpsfeed.GpsDevice) gpstime.Coord {
	return gpstime.GetLocationD(dev)
}
