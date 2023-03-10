# OVERVIEW
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/gpsinfo.svg)](https://pkg.go.dev/paepcke.de/gpsinfo) [![Go Report Card](https://goreportcard.com/badge/paepcke.de/gpsinfo)](https://goreportcard.com/report/paepcke.de/gpsinfo) [![Go Build](https://github.com/paepckehh/gpsinfo/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/gpsinfo/actions/workflows/golang.yml)

[paepcke.de/gpsinfo](https://paepcke.de/gpsinfo/)

- needs go1.20rc (sorry!)
- Show and decode nmea frames from an (usb) gps dongle.
- Focus on small embedded systems (debugging) on restricted resources.
- Focus on power saving parser (NO CLEAN idomatic go code for hot loop, no clean full state maschine, quick hack)
- 100 % pure go, stdlib only, no external dependencies, use as app or api (see api.go)

# INSTALL
```
go install paepcke.de/gpsinfo/cmd/gpsinfo@latest
```

### DOWNLOAD (prebuild)

[github.com/paepckehh/gpsinfo/releases](https://github.com/paepckehh/gpsinfo/releases)

# SHOWTIME

```Shell
gpsinfo /dev/gps0
[...]
```
# DOCS

[pkg.go.dev/paepcke.de/gpsinfo](https://pkg.go.dev/paepcke.de/gpsinfo)

# CONTRIBUTION

Yes, Please! PRs Welcome! 
