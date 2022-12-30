// package geohash ..
// [inspired|based|forked|enhanced|minimized] version of [github.com/mmcloughlin/geohash]
// to handle 3D Corrdinates (2D compressed geohash coordinates + altitude)
//
// 2D geohash source: [github.com/mmcloughlin/geohash]
// Copyright (c) 2015 Michael McLoughlin - The MIT License (MIT)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package geohash

// import
import (
	"math"
	"strconv"
	"strings"
)

//
// EXTERNAL INTERFACES
//

//
// Encode
//

// Encode ...
func Encode(a, o float64) string { return b32enc.enc(encodeIntWithP(a, o, 5*12))[:] }

// Encode2D ...
func Encode2D(a, o float64) string { return b32enc.enc(encodeIntWithP(a, o, 5*12))[:] }

// Encode3D ...
func Encode3D(a, o, l float64) string {
	return b32enc.enc(encodeIntWithP(a, o, 5*12))[:] + "@" + strconv.FormatFloat(l, 'f', -1, 64)
}

//
//  Decode
//

// Decode ...
func Decode(h string) (lat, long float64) {
	return boundingBoxIntWithP(b32enc.dec(h), uint(5*len(h))).round()
}

// Decode2D ...
func Decode2D(h string) (lat, long float64) {
	return boundingBoxIntWithP(b32enc.dec(h), uint(5*len(h))).round()
}

// Decode3D ...
func Decode3D(h string) (lat, long, h float64) {
	s := strings.Split(h, "@")
	if len(s) != 2 {
		panic("unable to continue, invalid 3D gehash code [" + h + "]")
	}
	l, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		panic("unable to continue, invalid 3D gehash code [" + h + "]")
	}
	a, o := boundingBoxIntWithP(b32enc.dec(h), uint(5*len(h))).round()
	return a, o, l
}

//
// INTERNAL LEGACY BACKEND
//

type (
	box      struct{ minLat, maxLat, minLong, maxLong float64 }
	encoding struct {
		encode string
		decode [256]byte
	}
)

var b32enc = newEncoding("0123456789bcdefghjkmnpqrstuvwxyz")

func encRange(x, r float64) uint32                       { return uint32((x + r) / (2 * r) * math.Exp2(32)) }
func decRange(x uint32, r float64) float64               { return 2*r*float64(x)/math.Exp2(32) - r }
func encInt(lat, long float64) uint64                    { return interleave(encRange(lat, 90), encRange(long, 180)) }
func interleave(x, y uint32) uint64                      { return spread(x) | (spread(y) << 1) }
func deinterleave(x uint64) (a,b uint32)                 { return squash(x), squash(x >> 1) }
func maxDecimalPower(r float64) float64                  { return math.Pow10(int(math.Floor(math.Log10(r)))) }
func encodeIntWithP(lat, long float64, bits uint) uint64 { return encInt(lat, long) >> (64 - bits) }
func errorWithP(bits uint) (a,b float64) {
	return math.Ldexp(180.0, -(int(bits) / 2)), math.Ldexp(360.0, -(int(bits) - int(bits)/2))
}

func (b box) round() (a,b float64) {
	x, y := maxDecimalPower(b.maxLat-b.minLat), maxDecimalPower(b.maxLong-b.minLong)
	return math.Ceil(b.minLat/x) * x, math.Ceil(b.minLong/y) * y
}

func boundingBoxIntWithP(hash uint64, bits uint) box {
	latInt, longInt := deinterleave(hash << (64 - bits))
	lat, long := decRange(latInt, 90), decRange(longInt, 180)
	latErr, longErr := errorWithP(bits)
	return box{minLat: lat, maxLat: lat + latErr, minLong: long, maxLong: long + longErr}
}

func spread(in uint32) uint64 {
	x := uint64(in)
	x = (x | (x << 16)) & 0x0000ffff0000ffff
	x = (x | (x << 8)) & 0x00ff00ff00ff00ff
	x = (x | (x << 4)) & 0x0f0f0f0f0f0f0f0f
	x = (x | (x << 2)) & 0x3333333333333333
	x = (x | (x << 1)) & 0x5555555555555555
	return x
}

func squash(x uint64) uint32 {
	x &= 0x5555555555555555
	x = (x | (x >> 1)) & 0x3333333333333333
	x = (x | (x >> 2)) & 0x0f0f0f0f0f0f0f0f
	x = (x | (x >> 4)) & 0x00ff00ff00ff00ff
	x = (x | (x >> 8)) & 0x0000ffff0000ffff
	x = (x | (x >> 16)) & 0x00000000ffffffff
	return uint32(x)
}

func newEncoding(in string) *encoding {
	e := new(encoding)
	e.encode = in
	for i := range e.decode {
		e.decode[i] = 0xff
	}
	for i := range in {
		e.decode[in[i]] = byte(i)
	}
	return e
}

func (e *encoding) dec(s string) uint64 {
	x := uint64(0)
	for i := range s {
		x = (x << 5) | uint64(e.decode[s[i]])
	}
	return x
}

func (e *encoding) enc(x uint64) string {
	b := [12]byte{}
	for i := range b {
		b[11-i] = e.encode[x&0x1f]
		x >>= 5
	}
	return string(b[:])
}
