// package gpsfeed
package gpsfeed

// import
import (
	"bufio"
	"os"
	"sync"
	"sync/atomic"
)

//
// Device
//

// DeviceTimeout ...
const DeviceTimeout = _deviceDataValidTimeout + 1 // wait before try to access device again

// GpsDevice holds all the feed, locks and handles for a specifc gps input device dongle
type GpsDevice struct {
	FileIO     string         // file name, eg. /dev/gps0
	Feed       *bufio.Scanner // line feed
	ErrCount   atomic.Uint64  // global err counter
	InitDone   atomic.Bool    // initial time sync done
	Responsive atomic.Bool    // feed state [responsive/unresponsive]
	DataValid  atomic.Bool    // feed state [datavalid/invalidutput]]
	Dog        atomic.Bool    // watchdog state
	Lock       sync.Mutex     // global device lock
	Handle     *os.File       // file handle
	Global     sync.WaitGroup // global state
}

// Open ...
func (dev *GpsDevice) Open() { openDev(dev) }

// Close ...
func (dev *GpsDevice) Close() { closeDev(dev) }

//
// Error Handling
//

const _errMax = 10

// GetErr returns the number of errors assoc with this device
func (dev *GpsDevice) GetErr() string { return getErr(dev) }

// CheckErrAdd validates the current number of global errors
func (dev *GpsDevice) CheckErrAdd() bool { return checkErrCount(dev) }

//
// Little Helper
//

// GetDeviceName ...
func GetDeviceName(arg, old string) string { return getDevice(arg, old) }

// CheckSumValid validates an NMEA sentence checksum
func CheckSumValid(sentence string) bool { return checkSumValid(sentence) }

// CheckSumTag calculates the NMEA sentence CheckSum and returns an tag string [ok|fail]
func CheckSumTag(sentence string) string { return checkSumTag(sentence) }
