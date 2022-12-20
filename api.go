// package gpsinfo ...
package gpsinfo

// import
import (
	"time"

	"paepcke.de/gpsinfo/gpsfeed"
)

//
// SIMPLE API
//

// Debug ...
func Debug(device string) { debug(&gpsfeed.GpsDevice{FileIO: device}) }

//
// GENERIC BACKEND
//

// DebugD ...
func DebugD(dev *gpsfeed.GpsDevice) { debug(dev) }

