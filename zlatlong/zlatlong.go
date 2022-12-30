// package zlatlong ...
// [inspired|based|forked|enhanced|minimized] version of [github.com/dgryski/go-zlatlong]
// to handle 3D Corrdinates (2D compressed coordinates + altitude)
//
// 2D source [github.com/dgryski/go-zlatlong]
// Copyright (c) 2016 Damian Gryski <damian@gryski.com> - The MIT License (MIT)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package zlatlong

import (
	"math"
	"strconv"
	"strings"
)

//
// EXTERNAL INTERFACES
//

// Encode2D ...
func Encode2D(a, o float64) string { return Encode(a, o) }

// Encode3D ...
func Encode3D(a, o, l float64) string {
	return Encode(a, o) + "@" + strconv.FormatFloat(l, 'f', -1, 64)
}

// Encode ...
func Encode(a, o float64) string {
	var result []byte
	dy, dx := round(a*100000), round(o*100000)
	dy = (dy << 1) ^ (dy >> 31)
	dx = (dx << 1) ^ (dx >> 31)
	index := ((dy + dx) * (dy + dx + 1) / 2) + dy
	for index > 0 {
		rem := index & 31
		index = (index - rem) / 32
		if index > 0 {
			rem += 32
		}
		result = append(result, safeCharacters[rem])
	}
	return string(result)
}

// Decode2D ...
func Decode2D(in string) (lat, long float64) { return Decode(in) }

// Decode3D ...
func Decode3D(h string) (lat, long, hh float64) {
	s := strings.Split(h, "@")
	if len(s) != 2 {
		panic("unable to continue, invalid 3D gehash code [" + h + "]")
	}
	l, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		panic("unable to continue, invalid 3D gehash code [" + h + "]")
	}
	a, o := Decode(s[0])
	return a, o, l
}

// Decode ...
func Decode(in string) (lat, long float64) {
	value := []byte(in)
	l := len(value)
	index, max, n, k := 0, int64(4294967296), int64(0), uint(0)
	for {
		if index >= l {
			panic("[zlatlong] [decode] [no valid data points in inputstring] [" + in + "]")
		}
		b := int64(safeIdx[value[index]])
		index++
		if b == 255 {
			panic("[zlatlong] [decode] [invalid (255) character in inputstring] [" + in + "]")
		}
		tmp := (b & 31) * (1 << k)
		ht, lt, hn, ln := tmp/max, tmp%max, n/max, n%max
		nl := (lt | ln)
		n = (ht|hn)*max + nl
		k += 5
		if b < 32 {
			break
		}
	}
	if l > index {
		panic("[zlatlong] [decode] [more than one datapoint in inputstring] [" + in + "] [not supported]")
	}
	diagonal := int64((math.Sqrt(8*float64(n)+5) - 1) / 2)
	n -= diagonal * (diagonal + 1) / 2
	nx, ny := diagonal-n, n
	nx = (nx >> 1) ^ -(nx & 1)
	ny = (ny >> 1) ^ -(ny & 1)
	return float64(ny) * 0.00001, float64(nx) * 0.00001
}

//
// INTERNAL BACKEND
//

var (
	safeCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")
	safeIdx        [256]byte
)

func init() {
	for i := range safeIdx {
		safeIdx[i] = 255
	}

	for i, c := range safeCharacters {
		safeIdx[c] = byte(i)
	}
}

func round(f float64) int64 {
	if f < 0 {
		return int64(f - 0.5)
	}

	return int64(f + 0.5)
}
