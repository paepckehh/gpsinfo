package gpsfeed

import (
	"os"
	"strconv"
	"time"
)

//
// DISPLAY IO
//

const _err = "[gpsfeed] [error]"

// out handles output messages to stdout, adding an linefeed
func out(msg string) { outPlain(msg + "\n") }

// outPlain handles output messages to stdout
func outPlain(msg string) { os.Stdout.Write([]byte(msg)) }

//
// FILE IO
//

const (
	_modeDevice  uint32 = 1 << (32 - 1 - 5)
	_modeSymlink uint32 = 1 << (32 - 1 - 4)
)

// getHandle
func getHandle(dev *GpsDevice) (handle *os.File, ok bool) {
	var err error
	handle, err = os.Open(dev.FileIO)
	if err != nil {
		if dev.CheckErrAdd() {
			out("[error]" + "[" + dev.FileIO + "] [" + err.Error() + "] [" + dev.GetErr() + "]")
		}
		time.Sleep(DeviceTimeout * time.Second)
		return handle, false
	}
	return handle, true
}

// getDevice
func getDevice(device, old string) string {
	if isDevice(device) {
		out("[info] [device] [" + device + "]")
		return device
	}
	out("[error] [skip] [unknown device or option] [" + device + "]")
	return old
}

// isDevice ...
func isDevice(devicename string) bool {
	fi, err := os.Lstat(devicename)
	if err != nil {
		return false
	}
	if uint32(fi.Mode())&_modeDevice != 0 || uint32(fi.Mode())&_modeSymlink != 0 {
		// TODO validate symlink target for device mode
		return true
	}
	return false
}

//
// ERROR HANDLER
//

// getErr returns the number of errors assoc with this device
func getErr(dev *GpsDevice) string {
	return strconv.FormatUint(dev.ErrCount.Load(), 10)
}

// checkErrCount validates the current number of global errors
func checkErrCount(dev *GpsDevice) bool {
	dev.ErrCount.Add(1)
	e := dev.ErrCount.Load()
	return e < _errMax
}
