# OVERVIEW

[paepche.de/gpsinfo](https://paepcke.de/gpsinfo)

- Show and decode nmea frames from an (usb) gps dongle.
- Focus on small embedded systems (debugging) on restricted resources.
- Focus on power saving parser (NO CLEAN idomatic go code for hot loop, no clean full state maschine, quick hack)
- 100 % pure go, stdlib only, no external dependencies, use as app or api (see api.go)

# INSTALL
```
go install paepcke.de/gpsinfo/cmd/gpsinfo@latest
```

# DOWNLOAD (prebuild)

[github.com/paepckehh/gpsinfo/releases](https://github.com/paepckehh/gpsinfo/releases)

# SHOWTIME

```Shell
gpsinfo /dev/gps0
[...]
```

# CONTRIBUTION

Yes, Please! PRs Welcome! 
