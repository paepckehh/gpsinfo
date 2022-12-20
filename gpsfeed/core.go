// package gpsfeed ...
package gpsfeed

// import
import (
	"bufio"
	"strings"
	"sync"
	"time"
)

const (
	_deviceResponsiveTimeout = 5 // fail after, when nothing is emitted from device at all
	_deviceDataValidTimeout  = 6 // fail after, when no valid parsable nemea stamps are emitted
	_errDeviceUnresponsive   = _err + "[device is unresponsive]"
	_errDeviceInvalidData    = _err + "[device does not emit valid gps nmea data stamps]"
)

//
// Device
//

// openDev sets the global lock and creates the device buffered io scanner
func openDev(dev *GpsDevice) {
	var ok bool
	watchdogs.Wait() // expire all old global watchdog locks
	dev.Lock.Lock()
	for {
		if dev.Handle, ok = getHandle(dev); !ok {
			continue
		}
		dev.Feed = bufio.NewScanner(dev.Handle)
		dev.Dog.Store(true)
		go watchdog(dev)
		return
	}
}

// closeDev do close/realease all locks and resets the watchdog
func closeDev(dev *GpsDevice) {
	dev.Dog.Store(false)
	dev.ErrCount.Store(0)
	dev.Handle.Close()
	dev.Lock.Unlock()
}

//
// Watchdogs
//

// ...
var watchdogs sync.WaitGroup

// watchdog background tasks will timeout and close the device file handle and terminate bufio scanner if needed
func watchdog(dev *GpsDevice) {
	watchdogs.Add(2)
	go func() {
		dev.DataValid.Store(true)
		oldErr := dev.ErrCount.Load()
		for dev.DataValid.Load() {
			if dev.ErrCount.Load() > oldErr+8 {
				break
			}
			dev.DataValid.Store(false) // re-arm deviceDataValid watchdog
			time.Sleep(_deviceDataValidTimeout * time.Second)
		}
		if dev.Dog.Swap(false) { // we are fist to trigger ?
			dev.Handle.Close()
			if dev.CheckErrAdd() {
				out(_errDeviceInvalidData + " [" + dev.FileIO + "] [" + dev.GetErr() + "]")
			}
		}
		watchdogs.Done()
	}()
	go func() {
		dev.Responsive.Store(true)
		oldErr := dev.ErrCount.Load()
		for dev.Responsive.Load() {
			if dev.ErrCount.Load() > oldErr+8 {
				break
			}
			dev.Responsive.Store(false) // re-arm deviceResponsive watchdog
			time.Sleep(_deviceResponsiveTimeout * time.Second)
		}
		if dev.Dog.Swap(false) { // we are fist to trigger ?
			dev.Handle.Close()
			if dev.CheckErrAdd() {
				out(_errDeviceUnresponsive + " [" + dev.FileIO + "] [" + dev.GetErr() + "]")
			}
		}
		watchdogs.Done()
	}()
}

//
// Little Helper
//

const (
	_sep      = ","
	_checksep = "*"

	_checksum = " [checksum] "
	_ok       = __GREEN + "[ok]" + __OFF
	_fail     = __RED + "[fail]" + __OFF

	__OFF   = "\033[0m"
	__RED   = "\033[2;31m"
	__GREEN = "\033[2;32m"

	digits = "0123456789ABCDEFX"
)

// checkSumTag ...
func checkSumTag(sentence string) string {
	if checkSumValid(sentence) {
		return _ok
	}
	return _fail
}

// checkSumValid ...
func checkSumValid(sentence string) bool {
	t := strings.Split(sentence, _checksep)
	r := t[0][1:]
	if checkSum(r) == t[1] {
		return true
	}
	return false
}

// checkSum ...
func checkSum(s string) string {
	var sum uint8
	for i := 0; i < len(s); i++ {
		sum ^= s[i]
	}
	var r []byte
	r = append(r, digits[sum>>4])
	r = append(r, digits[sum&0xF])
	return string(r)
}
