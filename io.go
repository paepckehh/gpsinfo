// package gpsinfo ...
package gpsinfo

// import
import (
	"os"
)

//
// DISPLAY IO
//

// out handles output messages to stdout, adding an linefeed
func out(msg string) { outPlain(msg + "\n") }

// outPlain handles output messages to stdout
func outPlain(msg string) { os.Stdout.Write([]byte(msg)) }
