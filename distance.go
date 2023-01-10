package gpsinfo

import (
	"fmt"
	"math"
	"sort"

	"paepcke.de/airloctag/airports"
)

// places
var places = map[string]coord{
	"Base                 ": {a: 053.564236667, o: 009.957351667, l: 26.5},
	"Hamburg              ": {a: 053.551086, o: 009.993682, l: 20},
	"Hamburg, Airport HAM ": {a: 053.624830, o: 009.987996, l: 53},
	"Hamburg, Airport XFW ": {a: 053.534831, o: 009.835496, l: 23},
	"Luebeck              ": {a: 053.869720, o: 010.686389, l: 50},
	"Copenhagen           ": {a: 055.676098, o: 012.568337, l: 150},
	"Berlin, Teufelsberg  ": {a: 052.497222, o: 013.241111, l: 150},
	"London               ": {a: 051.509865, o: -00.118092, l: 500},
	"Munich               ": {a: 048.137154, o: 011.576124, l: 1500},
	"Kiev                 ": {a: 050.450100, o: 030.523400, l: 550},
	"Tel Aviv             ": {a: 032.109333, o: 034.855499, l: 190},
	"Toulouse             ": {a: 043.604652, o: 001.444209, l: 500},
	"New York             ": {a: 040.730610, o: -73.935242, l: 20},
	"Reykjavik            ": {a: 064.128288, o: -21.827774, l: 0},
	"Sydney               ": {a: -33.865143, o: 151.209900, l: 20},
	"Pyongyang            ": {a: 039.019444, o: 125.738052, l: 110},
	"[+] North Pole       ": {a: 090.000000, o: 000.000000, l: 0},
	"[-] South Pole       ": {a: -90.000000, o: 000.000000, l: 0},
}

// coord ..
type coord struct {
	a float64 // Latitude
	o float64 // Longitude
	l float64 // Elevation
}

// dist ...
func dist(xa, xo, xl, ya, yo, yl float64) float64 {
	xa, xo, ya, yo = xa*(math.Pi/180), xo*(math.Pi/180), ya*(math.Pi/180), yo*(math.Pi/180)
	h := math.Pow(math.Sin((ya-xa)/2), 2) + math.Cos(xa)*math.Cos(ya)*math.Pow(math.Sin((yo-xo)/2), 2)
	return 2*6378100*math.Asin(math.Sqrt(h)) + math.Abs(xl-yl)
}

// displayAirports ...
func displayAirports(a, o, l float64) string {
	var dlist []float64
	var display string
	myDistance := make(map[float64]string, len(airports.Airports))
	for place, co := range airports.Airports {
		d := dist(a, o, l, co.A, co.O, co.L)
		myDistance[d] = place
		dlist = append(dlist, d)
	}
	sort.Float64s(dlist)
	for i := 0; i < 9; i++ {
		display += padS(fmt.Sprintf("%v: %v%v km%v ", myDistance[dlist[i]], _BLUE, int(dlist[i])/1000, _OFF))
	}
	display += "\n"
	return display
}

// displayMCD ...
func displayMCD(a, o, l float64) string {
	var dlist []float64
	var display string
	myDistance := make(map[float64]string, len(places))
	for place, co := range places {
		d := dist(a, o, l, co.a, co.o, co.l)
		myDistance[d] = place
		dlist = append(dlist, d)
	}
	sort.Float64s(dlist)
	loops := len(dlist)
	for i := 0; i < loops; {
		for c := 0; c < 3; c++ {
			s := ""
			switch {
			case dlist[i] < 3*1000:
				s = pad(fmt.Sprintf("%v: %v%v meter%v", myDistance[dlist[i]], _CYAN, int(dlist[i]), _OFF))
			case dlist[i] > 1000*1000:
				s = pad(fmt.Sprintf("%v: %v%v km%v", myDistance[dlist[i]], _GREY, int(dlist[i])/1000, _OFF))
			default:
				s = pad(fmt.Sprintf("%v: %v%v km%v", myDistance[dlist[i]], _BLUE, int(dlist[i])/1000, _OFF))
			}
			i++
			display += s
		}
		display += "\n"
	}
	return display
}
