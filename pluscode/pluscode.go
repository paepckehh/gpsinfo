// package pluscode ...
// This is a [code-minimal|feature-enhanced|optimized] fork of [github.com/google/open-location-code/go]
// This fork is not API or result compatible, please do not use outside this special use case. PLEASE USE ALWAYS THE ORIGINAL!
//
// [github.com/google/open-location-code/go]
// Copyright 2015 Tamás Gulácsi. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the 'License');
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an 'AS IS' BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package pluscode

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

//
// EXTERNAL INTERFACES
//

// Decode ...
func Decode(in string) (float64, float64) { return Decode2D(in) }

// Encode ...
func Encode(a, o float64) string { return encode(a, o, maxCodeLen) }

// Encode2D ...
func Encode2D(a, o float64) string { return encode(a, o, maxCodeLen) }

// Encode3D ...
func Encode3D(a, o, l float64) string {
	s := encode(a, o, maxCodeLen)
	return s + sep + strconv.FormatFloat(l, 'f', 2, 64)
}

// Decode2D ...
func Decode2D(in string) (float64, float64) {
	s, err := decode(in)
	if err != nil {
		panic("[pluscode] [unrecoverable-error] [invalid code] [decode] [" + err.Error() + "]")
	}
	return s.center()
}

// Decode3D ...
func Decode3D(in string) (float64, float64, float64) {
	s := strings.Split(in, sep)
	if len(s) != 2 {
		panic("[pluscode] [unrecoverable-error] [invalid 3D code] [split]")
	}
	c, err := decode(s[0])
	if err != nil {
		panic("[pluscode] [unrecoverable-error] [invalid 3D code] [decode] [" + err.Error() + "]")
	}
	a, o := c.center()
	l, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		panic("[pluscode] [unrecoverable-error] [invalid 3D code] [parse float] [" + err.Error() + "]")
	}
	return a, o, l
}

//
// INTERNAL BACKEND
//

type codeArea struct {
	latLo, longLo, latHi, longHi float64
	lengh                        int
}

const (
	minTrimmableCodeLen = 6
	sep                 = "@"
	separator           = '+'
	sepPos              = 8
	padding             = '0'
	alphabet            = "23456789CFGHJMPQRVWX"
	encBase             = len(alphabet)
	maxCodeLen          = 15
	pairCodeLen         = 10
	gridCodeLen         = maxCodeLen - pairCodeLen
	gridCols            = 4
	gridRows            = 5
	pairFPV             = 160000
	pairPrecision       = 8000
	gridLatFullValue    = 3125
	gridLongFullValue   = 1024
	gridLatFPV          = gridLatFullValue / gridRows
	gridLongFPV         = gridLongFullValue / gridCols
	finalLatPrecision   = pairPrecision * gridLatFullValue
	finalLongPrecision  = pairPrecision * gridLongFullValue
	latMax, longMax     = 90, 180
)

func normalizeLat(value float64) float64  { return normalize(value, latMax) }
func normalizeLong(value float64) float64 { return normalize(value, longMax) }
func decode(code string) (codeArea, error) {
	var area codeArea
	if err := checkFull(code); err != nil {
		return area, err
	}
	code = stripCode(code)
	l := len(code)
	if l < 2 {
		return area, errors.New("[code too short] [" + code + "]")
	}
	normalLat, normalLong := -latMax*pairPrecision, -longMax*pairPrecision
	extraLat, extraLong, digits, pv := 0, 0, pairCodeLen, pairFPV
	if l < digits {
		digits = l
	}
	for i := 0; i < digits-1; i += 2 {
		normalLat += strings.IndexByte(alphabet, code[i]) * pv
		normalLong += strings.IndexByte(alphabet, code[i+1]) * pv
		if i < digits-2 {
			pv /= encBase
		}
	}
	latPrecision, longPrecision := float64(pv)/pairPrecision, float64(pv)/pairPrecision
	if l > pairCodeLen {
		rowpv := gridLatFPV
		colpv := gridLongFPV
		digits = maxCodeLen
		if l < maxCodeLen {
			digits = l
		}
		for i := pairCodeLen; i < digits; i++ {
			dval := strings.IndexByte(alphabet, code[i])
			row, col := dval/gridCols, dval%gridCols
			extraLat += row * rowpv
			extraLong += col * colpv
			if i < digits-1 {
				rowpv /= gridRows
				colpv /= gridCols
			}
		}
		latPrecision, longPrecision = float64(rowpv)/finalLatPrecision, float64(colpv)/finalLongPrecision
	}
	lat := float64(normalLat)/pairPrecision + float64(extraLat)/finalLatPrecision
	long := float64(normalLong)/pairPrecision + float64(extraLong)/finalLongPrecision
	return codeArea{
		latLo:  math.Round(lat*1e14) / 1e14,
		longLo: math.Round(long*1e14) / 1e14,
		latHi:  math.Round((lat+latPrecision)*1e14) / 1e14,
		longHi: math.Round((long+longPrecision)*1e14) / 1e14,
		lengh:  l,
	}, nil
}

func encode(lat, long float64, codeLen int) string {
	switch {
	case codeLen <= 0:
		codeLen = pairCodeLen
	case codeLen < 2:
		codeLen = 2
	case codeLen < pairCodeLen && codeLen%2 == 1:
		codeLen++
	case codeLen > maxCodeLen:
		codeLen = maxCodeLen
	}
	lat, long = clipLatitude(lat), normalizeLong(long)
	if lat == latMax {
		lat = normalizeLat(lat - computeLatPrec(codeLen))
	}
	var code [15]byte
	var latVal int64 = int64(math.Round((lat+latMax)*finalLatPrecision*1e6) / 1e6)
	var longVal int64 = int64(math.Round((long+longMax)*finalLongPrecision*1e6) / 1e6)
	pos := maxCodeLen - 1
	switch {
	case codeLen > pairCodeLen:
		for i := 0; i < gridCodeLen; i++ {
			latDigit := latVal % int64(gridRows)
			longDigit := longVal % int64(gridCols)
			ndx := latDigit*gridCols + longDigit
			code[pos] = alphabet[ndx]
			pos--
			latVal /= int64(gridRows)
			longVal /= int64(gridCols)
		}
	default:
		latVal /= gridLatFullValue
		longVal /= gridLongFullValue
	}
	pos = pairCodeLen - 1
	for i := 0; i < pairCodeLen/2; i++ {
		latNdx := latVal % int64(encBase)
		longNdx := longVal % int64(encBase)
		code[pos] = alphabet[longNdx]
		pos--
		code[pos] = alphabet[latNdx]
		pos--
		latVal /= int64(encBase)
		longVal /= int64(encBase)
	}
	if codeLen >= sepPos {
		return string(code[:sepPos]) + string(separator) + string(code[sepPos:codeLen])
	}
	return string(code[:codeLen]) + strings.Repeat(string(padding), sepPos-codeLen) + string(separator)
}

func computeLatPrec(codeLen int) float64 {
	if codeLen <= pairCodeLen {
		return math.Pow(float64(encBase), float64(codeLen/-2+2))
	}
	return math.Pow(float64(encBase), -3) / math.Pow(float64(gridRows), float64(codeLen-pairCodeLen))
}

func (area codeArea) center() (lat, long float64) {
	return math.Min(area.latLo+(area.latHi-area.latLo)/2, latMax), math.Min(area.longLo+(area.longHi-area.longLo)/2, longMax)
}

func check(code string) error {
	n := len(code)
	if code == "" || n == 1 && code[0] == separator {
		return errors.New("[empty code]")
	}
	firstSep, firstPad := -1, -1
	for i, r := range code {
		if firstPad != -1 {
			switch r {
			case padding:
				continue
			case separator:
				if firstSep != -1 {
					return errors.New("[extraneous separator] [" + code + "]")
				}
				firstSep = i
				if n-1 == i {
					continue
				}
			}
			return errors.New("[after zero]")
		}
		if '2' <= r && r <= '9' {
			continue
		}
		switch r {
		case 'C', 'F', 'G', 'H', 'J', 'M', 'P', 'Q', 'R', 'V', 'W', 'X',
			'c', 'f', 'g', 'h', 'j', 'm', 'p', 'q', 'r', 'v', 'w', 'x':
			continue
		case separator:
			switch {
			case firstSep != -1:
				return errors.New("[extra separator seen] [" + code + "]")
			case i > sepPos || i%2 == 1:
				return errors.New("[separator in illegal position] [" + code + "]")
			}
			firstSep = i
		case padding:
			if i == 0 {
				return errors.New("[shouldn't start with padding character] [" + code + "]")
			}
			firstPad = i
		default:
			return errors.New("[invalid char] [" + code + "]")
		}
	}
	switch {
	case firstSep == -1:
		return errors.New("[missing separator] [" + code + "]")
	case n-firstSep-1 == 1:
		return errors.New("[only one char after separator] [" + code + "]")
	case firstPad != -1:
		switch {
		case firstSep < sepPos:
			return errors.New("[short codes cannot have padding] [" + code + "]")
		case n-firstPad-1%2 == 1:
			return errors.New("[odd number of padding chars] [" + code + "]")
		}
	}
	return nil
}

func checkFull(code string) error {
	if err := check(code); err != nil {
		return err
	}
	if i := strings.IndexByte(code, separator); i <= 0 && i > sepPos {
		return errors.New("[code too short] [" + code + "]")
	}
	switch {
	case strings.IndexByte(alphabet, upper(code[0]))*encBase >= latMax*2:
		return errors.New("[latitude outside range] [" + code + "]")
	case len(code) == 1:
		return nil
	case strings.IndexByte(alphabet, upper(code[1]))*encBase >= longMax*2:
		return errors.New("[longitude outside range] [" + code + "]")
	}
	return nil
}

func upper(b byte) byte {
	if 'c' <= b && b <= 'x' {
		return b + 'C' - 'c'
	}
	return b
}

func stripCode(code string) string {
	code = strings.Map(
		func(r rune) rune {
			if r == separator || r == padding {
				return -1
			}
			return rune(upper(byte(r)))
		},
		code)
	if len(code) > maxCodeLen {
		return code[:maxCodeLen]
	}
	return code
}

func normalize(value, max float64) float64 {
	switch {
	case value < -max:
		value += 2 * max
	case value >= max:
		value -= 2 * max
	}
	return value
}

func clipLatitude(lat float64) float64 {
	switch {
	case lat > latMax:
		return latMax
	case lat < -latMax:
		return -latMax
	default:
		return lat
	}
}
