// package gpsinfo decodes gps nmea frames from your gpsdongle (debug/ingo)
package gpsinfo

import "paepcke.de/gpsinfo/gpsfeed"

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
