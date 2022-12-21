# Overview

[paepche.de/gpsinfo](https://paepcke.de/gpsinfo)

- Show and decode information from an (usb) gps dongle.
- Focus on small embedded systems (debugging) on restricted resources.
- Focus onpower saving parser (NO CLEAN idomatic go code for hot loop, no clean full state maschine, quick hack)
- 100 % pure go, stdlib only, no external dependencies, use as app or api (see api.go)

## Install 
```
go install paepcke.de/gpsinfo/cmd/gpsinfo@latest
```

# Showtime 

```Shell
gpsinfo /dev/gps0
```
