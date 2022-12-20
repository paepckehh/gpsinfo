// package gpsinfo
package gpsinfo

// import
import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"paepcke.de/airloctag"
	"paepcke.de/daylight/sun"
	"paepcke.de/gpsfeed"
	"paepcke.de/gpsinfo/geohash"
	"paepcke.de/gpsinfo/nmeanano"
	"paepcke.de/gpsinfo/pluscode"
	"paepcke.de/gpsinfo/zlatlong"
)

const _TS = "15:04:05" // time stamp layout [time.Parse]

// build ...
func build(dev *gpsfeed.GpsDevice, channelGpsFrames, channelOut chan string, displaySleep time.Duration) {
	var (
		counter, fix          int64
		fixQuality            string
		err                   error
		cDIFF, cXX            string
		NumberSVsInView       string
		indicator             string
		hash, airp, dist      string
		zl, pl, td, gh, at    string
		oldLat, oldLong       float64
		t, daylight, dtime    time.Duration
		sunrise, sunset, noon time.Time
		tsSys, tsGps          time.Time

		// shortcut objects for nmeanano structs
		s nmeanano.Sentence
		v nmeanano.VTG
		g nmeanano.GSV
		a nmeanano.GSA
		x nmeanano.GGA
		y nmeanano.GNS
		m nmeanano.RMC
		o nmeanano.VDMVDO
	)
	// set defaults
	m.BaseSentence.Raw = _defaults
	v.BaseSentence.Raw = _defaults
	g.BaseSentence.Raw = _defaults
	a.BaseSentence.Raw = _defaults
	x.BaseSentence.Raw = _defaults
	y.BaseSentence.Raw = _defaults
	o.BaseSentence.Raw = _defaults

	// main loop / sentence subscriber
	for line := range channelGpsFrames {
		tsSys = time.Now()
		var b strings.Builder
		if s, err = nmeanano.Parse(line); err != nil {
			continue
		}
		switch s.DataType() {
		case "RMC":
			m = s.(nmeanano.RMC)
			tsGps = nmeanano.GetTimeStamp(m)
		case "GSV":
			g = s.(nmeanano.GSV)
			if g.NumberSVsInView < 6 {
				cXX = _ALERT
			} else {
				cXX = _ALERT_G
			}
			indicator = cXX + strings.Repeat("|", int(g.NumberSVsInView*4)) + _OFF
			NumberSVsInView = fmt.Sprintf("%s %s[%v]%s", indicator, _CYAN, g.NumberSVsInView, _OFF)
		case "GSA":
			a = s.(nmeanano.GSA)
		case "GGA":
			x = s.(nmeanano.GGA)
			fix, _ = strconv.ParseInt(x.FixQuality, 10, 0)
			if int(fix) > 0 {
				fixQuality = "ok"
			} else {
				fixQuality = "invalid"
			}
		case "VTG":
			v = s.(nmeanano.VTG)
		case "GNS":
			y = s.(nmeanano.GNS)
		case "VDO":
			o = s.(nmeanano.VDMVDO)
		case "VDM":
			o = s.(nmeanano.VDMVDO)
		default:
			continue
		}
		if oldLat != m.Latitude && oldLong != m.Longitude {
			oldLat, oldLong = m.Latitude, m.Longitude
			dist = displayMCD(m.Latitude, m.Longitude, x.Altitude)
			airp = "IATA " + displayAirports(m.Latitude, m.Longitude, x.Altitude)
			sunrise, sunset, noon, daylight = sun.State(m.Latitude, m.Longitude, x.Altitude)
			hash, _, _ = airloctag.Encode(m.Latitude, m.Longitude, x.Altitude, "", 0)
			at = fmt.Sprintf("AirLocTag            : %s%v%s\n", _BLUE, hash, _OFF)
			gh = fmt.Sprintf("GeoHash              : %s%v%s\n", _BLUE, geohash.Encode2D(m.Latitude, m.Longitude), _OFF)
			pl = fmt.Sprintf("Pluscode [googlemap] : %s%v%s\n", _BLUE, pluscode.Encode2D(m.Latitude, m.Longitude), _OFF)
			zl = fmt.Sprintf("zLatLong [bingmap]   : %s%v%s\n", _BLUE, zlatlong.Encode2D(m.Latitude, m.Longitude), _OFF)
			td = fmt.Sprintf("Today [UTC]          : Sunrise %s%v%s      Sunset %s%v%s      Noon %s%v%s      Daylight %s%v%s\n", _CYAN, sunrise.Format(_TS), _OFF, _CYAN, sunset.Format(_TS), _OFF, _CYAN, noon.Format(_TS), _OFF, _CYAN, daylight, _OFF)
		}
		fmt.Fprintf(&b, _cleanNewline)
		fmt.Fprintf(&b, "\n\n\n\n\n\n")
		fmt.Fprintf(&b, _sectionLine)
		fmt.Fprintf(&b, "SENSOR DEVICE PORT   : %s%s%s\n", _BLUE, dev.FileIO, _OFF)
		fmt.Fprintf(&b, "RAW RMC STAMP        : %s%v%s\n", _GREY, m, _OFF)
		fmt.Fprintf(&b, "RAW GSV STAMP        : %s%v%s\n", _GREY, g, _OFF)
		fmt.Fprintf(&b, "RAW GSA STAMP        : %s%v%s\n", _GREY, a, _OFF)
		fmt.Fprintf(&b, "RAW GGA STAMP        : %s%v%s\n", _GREY, x, _OFF)
		fmt.Fprintf(&b, "RAW VTG STAMP        : %s%v%s\n", _GREY, v, _OFF)
		fmt.Fprintf(&b, "RAW GNS STAMP        : %s%v%s\n", _GREY, y, _OFF)
		fmt.Fprintf(&b, "RAW VDMVDO STAMP     : %s%v%s\n", _GREY, o, _OFF)
		fmt.Fprintf(&b, _sectionLine)
		fmt.Fprintf(&b, "AVDM MessageID       : %s%v%s    Channel: %s%v%s    Payload: %s%v%s\n", _BLUE, o.MessageID, _OFF, _BLUE, o.Channel, _OFF, _BLUE, o.Payload, _OFF)
		fmt.Fprintf(&b, "Orientation          : %s%.4f%s\n", _BLUE, m.Course, _OFF)
		fmt.Fprintf(&b, "Variation            : %s%.4f%s\n", _BLUE, m.Variation, _OFF)
		fmt.Fprintf(&b, "Speed                : %s%.4f%s [kmh] %s%.4f%s [knots]\n", _BLUE, (m.Speed * 1.852), _OFF, _BLUE, m.Speed, _OFF)
		fmt.Fprintf(&b, "Altitude             : %s%.1f%s [meter]\n", _BLUE, x.Altitude, _OFF)
		fmt.Fprintf(&b, "DMS Latitude         : %s%s%s\n", _CYAN, nmeanano.FormatDMS(m.Latitude), _OFF)
		fmt.Fprintf(&b, "DMS Longitude        : %s%s%s\n", _CYAN, nmeanano.FormatDMS(m.Longitude), _OFF)
		fmt.Fprintf(&b, "GPS Latitude         : %s%.9f%s\n", _BLUE, m.Latitude, _OFF)
		fmt.Fprintf(&b, "GPS Longitude        : %s%.9f%s\n", _BLUE, m.Longitude, _OFF)
		if m.Latitude != 0 && m.Longitude != 0 {
			fmt.Fprintf(&b, gh)
			fmt.Fprintf(&b, zl)
			fmt.Fprintf(&b, pl)
			fmt.Fprintf(&b, at)
			fmt.Fprintf(&b, _sectionLine)
			fmt.Fprintf(&b, airp)
			fmt.Fprintf(&b, _sectionLine)
			fmt.Fprintf(&b, dist)
			fmt.Fprintf(&b, _sectionLine)
			fmt.Fprintf(&b, td)
			fmt.Fprintf(&b, _sectionLine)
		}
		fmt.Fprintf(&b, "Sat's [visible]      : %s\n", NumberSVsInView)
		for i := range g.Info {
			fmt.Fprintf(&b, " + SV-PRN %s%2d%s SNR %s%2d%s Elevation %s%2d%s Azimuth %s%3d%s \n", _BLUE, g.Info[i].SVPRNNumber, _OFF, _BLUE, g.Info[i].SNR, _OFF, _BLUE, g.Info[i].Elevation, _OFF, _BLUE, g.Info[i].Azimuth, _OFF)
		}
		fmt.Fprintf(&b, _sectionLine)
		fmt.Fprintf(&b, "Fix Dilution         : Type %s%v%s Mode %s%v%s Precision Dilution %s%v%s ( Horizontal %s%v%s Vertical %s%v%s )\n", _BLUE, a.Mode, _OFF, _BLUE, a.Type, _OFF, _BLUE, a.PDOP, _OFF, _BLUE, a.HDOP, _OFF, _BLUE, a.VDOP, _OFF)
		fmt.Fprintf(&b, "Fix Quality          : %s[%s]%s\n", _ALERT_G, fixQuality, _OFF)
		fmt.Fprintf(&b, "Fix used Sat's       : %s%v%s\n", _BLUE, x.NumSatellites, _OFF)
		fmt.Fprintf(&b, "Fix Time             : %s%v%s\n", _BLUE, x.Time, _OFF)
		if m.Date.Valid && m.Time.Valid {
			fmt.Fprintf(&b, "Validity             : %s%s [verified]%s\n", _ALERT_G, m.Validity, _OFF)
		} else {
			fmt.Fprintf(&b, "Validity             : %s%s [faild]%s\n", _ALERT, m.Validity, _OFF)
		}
		t = tsGps.Sub(tsSys)
		if t < maxDiff {
			cDIFF = _ALERT_G
		} else {
			cDIFF = _ALERT
		}
		fmt.Fprintf(&b, "Time Difference      : %s%s%s\n", cDIFF, t, _OFF)
		fmt.Fprintf(&b, "GPS Time             : %s%s%s\n", _BLUE, tsGps, _OFF)
		fmt.Fprintf(&b, "Local System Time    : %s%s%s\n", _BLUE, tsSys, _OFF)
		fmt.Fprintf(&b, _sectionLine)
		counter++
		dtime = (time.Since(tsSys))
		fmt.Fprintf(&b, "%s###  Time needed to decode Frame: %v  ###  Total Frames: %v%s\n", _GREY, dtime, counter, _OFF)
		channelOut <- b.String()
		time.Sleep(displaySleep)
	}
}
