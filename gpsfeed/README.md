# Overview 

- library to handle gps nmea emitting devices via bufio.Scanner interface
- automatic recovery (watchdog), handling cecksums,
- can detect if the devices is unresponsive, emitts defective frames, disconnects, missbehaves ...
- 100 % pure go, stdlib only, no external dependencies 
- see api.go for more details, cmd/gpsfeed for an example app
