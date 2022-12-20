// package nmeanano ...
// This is a [minimal|static|optimize] fork of the great [github.com/adrianmo/go-nmea] pkg!
// This fork is adapted and optimized for an specific use-case and *IS NOT* api/result compatible
// with the original, please DO NOT use outside very specific use cases!.PLEASE USE ALWAYS THE ORIGINAL!
//
// [github.com/adrianmo/go-nmea]
//
// Copyright (c) 2015 Adrian Moreno - The MIT License (MIT)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// (of this software and associated documentation files (the "Software"), to deal
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
package nmeanano

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

//
// External New Interfaces
//

// GetTimeStamp returns an validated go time.Time timestamp
func GetTimeStamp(x RMC) time.Time {
	return time.Date(2000+x.Date.YY,
		time.Month(x.Date.MM),
		x.Date.DD,
		x.Time.Hour,
		x.Time.Minute,
		x.Time.Second,
		x.Time.Millisecond*1000*1000,
		time.UTC)
}

//
// Internal Backend
//

const (
	TypeGGA = "GGA"
	Invalid = "0"
	GPS     = "1"
	DGPS    = "2"
	PPS     = "3"
	RTK     = "4"
	FRTK    = "5"
	EST     = "6"
)

// GGA ...
type GGA struct {
	BaseSentence
	Time          Time
	Latitude      float64
	Longitude     float64
	FixQuality    string
	NumSatellites int64
	HDOP          float64
	Altitude      float64
	Separation    float64
	DGPSAge       string
	DGPSId        string
}

func newGGA(s BaseSentence) (GGA, error) {
	p := NewParser(s)
	p.AssertType(TypeGGA)
	return GGA{
		BaseSentence:  s,
		Time:          p.Time(0, "time"),
		Latitude:      p.LatLong(1, 2, "latitude"),
		Longitude:     p.LatLong(3, 4, "longitude"),
		FixQuality:    p.EnumString(5, "fix quality", Invalid, GPS, DGPS, PPS, RTK, FRTK, EST),
		NumSatellites: p.Int64(6, "number of satellites"),
		HDOP:          p.Float64(7, "hdop"),
		Altitude:      p.Float64(8, "altitude"),
		Separation:    p.Float64(10, "separation"),
		DGPSAge:       p.String(12, "dgps age"),
		DGPSId:        p.String(13, "dgps id"),
	}, p.Err()
}

const (
	TypeGNS              = "GNS"
	NoFixGNS             = "N"
	AutonomousGNS        = "A"
	DifferentialGNS      = "D"
	PreciseGNS           = "P"
	RealTimeKinematicGNS = "R"
	FloatRTKGNS          = "F"
	EstimatedGNS         = "E"
	ManualGNS            = "M"
	SimulatorGNS         = "S"
)

// GNS ...
type GNS struct {
	BaseSentence
	Time       Time
	Latitude   float64
	Longitude  float64
	Mode       []string
	SVs        int64
	HDOP       float64
	Altitude   float64
	Separation float64
	Age        float64
	Station    int64
}

func newGNS(s BaseSentence) (GNS, error) {
	p := NewParser(s)
	p.AssertType(TypeGNS)
	m := GNS{
		BaseSentence: s,
		Time:         p.Time(0, "time"),
		Latitude:     p.LatLong(1, 2, "latitude"),
		Longitude:    p.LatLong(3, 4, "longitude"),
		Mode:         p.EnumChars(5, "mode", NoFixGNS, AutonomousGNS, DifferentialGNS, PreciseGNS, RealTimeKinematicGNS, FloatRTKGNS, EstimatedGNS, ManualGNS, SimulatorGNS),
		SVs:          p.Int64(6, "SVs"),
		HDOP:         p.Float64(7, "HDOP"),
		Altitude:     p.Float64(8, "altitude"),
		Separation:   p.Float64(9, "separation"),
		Age:          p.Float64(10, "age"),
		Station:      p.Int64(11, "station"),
	}
	return m, p.Err()
}

const (
	TypeGSA = "GSA"
	Auto    = "A"
	Manual  = "M"
	FixNone = "1"
	Fix2D   = "2"
	Fix3D   = "3"
)

// GSA ...
type GSA struct {
	BaseSentence
	Mode    string
	FixType string
	SV      []string
	PDOP    float64
	HDOP    float64
	VDOP    float64
}

func newGSA(s BaseSentence) (GSA, error) {
	p := NewParser(s)
	p.AssertType(TypeGSA)
	m := GSA{
		BaseSentence: s,
		Mode:         p.EnumString(0, "selection mode", Auto, Manual),
		FixType:      p.EnumString(1, "fix type", FixNone, Fix2D, Fix3D),
	}
	for i := 2; i < 14; i++ {
		if v := p.String(i, "satellite in view"); v != "" {
			m.SV = append(m.SV, v)
		}
	}
	m.PDOP = p.Float64(14, "pdop")
	m.HDOP = p.Float64(15, "hdop")
	m.VDOP = p.Float64(16, "vdop")
	return m, p.Err()
}

const (
	TypeGSV = "GSV"
)

// GSV ...
type GSV struct {
	BaseSentence
	TotalMessages   int64
	MessageNumber   int64
	NumberSVsInView int64
	Info            []GSVInfo
}

// GGSVInfo ...
type GSVInfo struct {
	SVPRNNumber int64
	Elevation   int64
	Azimuth     int64
	SNR         int64
}

func newGSV(s BaseSentence) (GSV, error) {
	p := NewParser(s)
	p.AssertType(TypeGSV)
	m := GSV{
		BaseSentence:    s,
		TotalMessages:   p.Int64(0, "total number of messages"),
		MessageNumber:   p.Int64(1, "message number"),
		NumberSVsInView: p.Int64(2, "number of SVs in view"),
	}
	for i := 0; i < 4; i++ {
		if 5*i+4 > len(m.Fields) {
			break
		}
		m.Info = append(m.Info, GSVInfo{
			SVPRNNumber: p.Int64(3+i*4, "SV prn number"),
			Elevation:   p.Int64(4+i*4, "elevation"),
			Azimuth:     p.Int64(5+i*4, "azimuth"),
			SNR:         p.Int64(6+i*4, "SNR"),
		})
	}
	return m, p.Err()
}

// Parser ...
type Parser struct {
	BaseSentence
	err error
}

// NewParser ...
func NewParser(s BaseSentence) *Parser {
	return &Parser{BaseSentence: s}
}

func (p *Parser) AssertType(typ string) {
	if p.Type != typ {
		p.SetErr("type", p.Type)
	}
}

// Err ...
func (p *Parser) Err() error {
	return p.err
}

// SetErr ...
func (p *Parser) SetErr(context, value string) {
	if p.err == nil {
		p.err = fmt.Errorf("nmea: %s invalid %s: %s", p.Prefix(), context, value)
	}
}

// String ...
func (p *Parser) String(i int, context string) string {
	if p.err != nil {
		return ""
	}
	if i < 0 || i >= len(p.Fields) {
		p.SetErr(context, "index out of range")
		return ""
	}
	return p.Fields[i]
}

// ListString ...
func (p *Parser) ListString(from int, context string) (list []string) {
	if p.err != nil {
		return []string{}
	}
	if from < 0 || from >= len(p.Fields) {
		p.SetErr(context, "index out of range")
		return []string{}
	}
	return append(list, p.Fields[from:]...)
}

// EnumString ...
func (p *Parser) EnumString(i int, context string, options ...string) string {
	s := p.String(i, context)
	if p.err != nil || s == "" {
		return ""
	}
	for _, o := range options {
		if o == s {
			return s
		}
	}
	p.SetErr(context, s)
	return ""
}

// EnumChars ...
func (p *Parser) EnumChars(i int, context string, options ...string) []string {
	s := p.String(i, context)
	if p.err != nil || s == "" {
		return []string{}
	}
	strs := []string{}
	for _, r := range s {
		rs := string(r)
		for _, o := range options {
			if o == rs {
				strs = append(strs, o)
				break
			}
		}
	}
	if len(strs) != len(s) {
		p.SetErr(context, s)
		return []string{}
	}
	return strs
}

// Int64 ...
func (p *Parser) Int64(i int, context string) int64 {
	s := p.String(i, context)
	if p.err != nil {
		return 0
	}
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		p.SetErr(context, s)
	}
	return v
}

// Float64 ...
func (p *Parser) Float64(i int, context string) float64 {
	s := p.String(i, context)
	if p.err != nil {
		return 0
	}
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		p.SetErr(context, s)
	}
	return v
}

// Time ...
func (p *Parser) Time(i int, context string) Time {
	s := p.String(i, context)
	if p.err != nil {
		return Time{}
	}
	v, err := ParseTime(s)
	if err != nil {
		p.SetErr(context, s)
	}
	return v
}

// Date ...
func (p *Parser) Date(i int, context string) Date {
	s := p.String(i, context)
	if p.err != nil {
		return Date{}
	}
	v, err := ParseDate(s)
	if err != nil {
		p.SetErr(context, s)
	}
	return v
}

// LatLong ...
func (p *Parser) LatLong(i, j int, context string) float64 {
	a := p.String(i, context)
	b := p.String(j, context)
	if p.err != nil {
		return 0
	}
	s := fmt.Sprintf("%s %s", a, b)
	v, err := ParseLatLong(s)
	if err != nil {
		p.SetErr(context, err.Error())
	}
	if (b == North || b == South) && (v < -90.0 || 90.0 < v) {
		p.SetErr(context, "latitude is not in range (-90, 90)")
		return 0
	} else if (b == West || b == East) && (v < -180.0 || 180.0 < v) {
		p.SetErr(context, "longitude is not in range (-180, 180)")
		return 0
	}
	return v
}

// SixBitASCIIArmour ..
func (p *Parser) SixBitASCIIArmour(i, fillBits int, context string) []byte {
	if p.err != nil {
		return nil
	}
	if fillBits < 0 || fillBits >= 6 {
		p.SetErr(context, "fill bits")
		return nil
	}
	payload := []byte(p.String(i, "encoded payload"))
	numBits := len(payload)*6 - fillBits
	if numBits < 0 {
		p.SetErr(context, "num bits")
		return nil
	}
	result := make([]byte, numBits)
	resultIndex := 0
	for _, v := range payload {
		if v < 48 || v >= 120 {
			p.SetErr(context, "data byte")
			return nil
		}
		d := v - 48
		if d > 40 {
			d -= 8
		}
		for i := 5; i >= 0 && resultIndex < len(result); i-- {
			result[resultIndex] = (d >> uint(i)) & 1
			resultIndex++
		}
	}
	return result
}

const (
	TypeRMC    = "RMC"
	ValidRMC   = "A"
	InvalidRMC = "V"
)

// RMC ...
type RMC struct {
	BaseSentence
	Time      Time
	Validity  string
	Latitude  float64
	Longitude float64
	Speed     float64
	Course    float64
	Date      Date
	Variation float64
}

func newRMC(s BaseSentence) (RMC, error) {
	p := NewParser(s)
	p.AssertType(TypeRMC)
	m := RMC{
		BaseSentence: s,
		Time:         p.Time(0, "time"),
		Validity:     p.EnumString(1, "validity", ValidRMC, InvalidRMC),
		Latitude:     p.LatLong(2, 3, "latitude"),
		Longitude:    p.LatLong(4, 5, "longitude"),
		Speed:        p.Float64(6, "speed"),
		Course:       p.Float64(7, "course"),
		Date:         p.Date(8, "date"),
		Variation:    p.Float64(9, "variation"),
	}
	if p.EnumString(10, "direction", West, East) == West {
		m.Variation = 0 - m.Variation
	}
	return m, p.Err()
}

const (
	SentenceStart             = "$"
	SentenceStartEncapsulated = "!"
	FieldSep                  = ","
	ChecksumSep               = "*"
)

var (
	customParsersMu = &sync.Mutex{}
	customParsers   = map[string]ParserFunc{}
)

type (
	ParserFunc func(BaseSentence) (Sentence, error)
	Sentence   interface {
		fmt.Stringer
		Prefix() string
		DataType() string
		TalkerID() string
	}
)

// BaseSentence ...
type BaseSentence struct {
	Talker   string
	Type     string
	Fields   []string
	Checksum string
	Raw      string
	TagBlock TagBlock
}

func (s BaseSentence) Prefix() string {
	return s.Talker + s.Type
}

func (s BaseSentence) DataType() string {
	return s.Type
}

func (s BaseSentence) TalkerID() string {
	return s.Talker
}
func (s BaseSentence) String() string { return s.Raw }
func parseSentence(raw string) (BaseSentence, error) {
	raw = strings.TrimSpace(raw)
	tagBlockParts := strings.SplitN(raw, `\`, 3)
	var (
		tagBlock TagBlock
		err      error
	)
	if len(tagBlockParts) == 3 {
		tags := tagBlockParts[1]
		raw = tagBlockParts[2]
		tagBlock, err = parseTagBlock(tags)
		if err != nil {
			return BaseSentence{}, err
		}
	}
	startIndex := strings.IndexAny(raw, SentenceStart+SentenceStartEncapsulated)
	if startIndex != 0 {
		return BaseSentence{}, fmt.Errorf("nmea: sentence does not start with a '$' or '!'")
	}
	sumSepIndex := strings.Index(raw, ChecksumSep)
	if sumSepIndex == -1 {
		return BaseSentence{}, fmt.Errorf("nmea: sentence does not contain checksum separator")
	}
	var (
		fieldsRaw   = raw[startIndex+1 : sumSepIndex]
		fields      = strings.Split(fieldsRaw, FieldSep)
		checksumRaw = strings.ToUpper(raw[sumSepIndex+1:])
		checksum    = Checksum(fieldsRaw)
	)
	if checksum != checksumRaw {
		return BaseSentence{}, fmt.Errorf(
			"nmea: sentence checksum mismatch [%s != %s]", checksum, checksumRaw)
	}
	talker, typ := parsePrefix(fields[0])
	return BaseSentence{
		Talker:   talker,
		Type:     typ,
		Fields:   fields[1:],
		Checksum: checksumRaw,
		Raw:      raw,
		TagBlock: tagBlock,
	}, nil
}

func parsePrefix(s string) (string, string) {
	if strings.HasPrefix(s, "PMTK") {
		return "PMTK", s[4:]
	}
	if strings.HasPrefix(s, "P") {
		return "P", s[1:]
	}
	if len(s) < 2 {
		return s, ""
	}
	return s[:2], s[2:]
}

func Checksum(s string) string {
	var checksum uint8
	for i := 0; i < len(s); i++ {
		checksum ^= s[i]
	}
	return fmt.Sprintf("%02X", checksum)
}

func MustRegisterParser(sentenceType string, parser ParserFunc) {
	if err := RegisterParser(sentenceType, parser); err != nil {
		panic(err)
	}
}

func RegisterParser(sentenceType string, parser ParserFunc) error {
	customParsersMu.Lock()
	defer customParsersMu.Unlock()
	if _, ok := customParsers[sentenceType]; ok {
		return fmt.Errorf("nmea: parser for sentence type '%q' already exists", sentenceType)
	}
	customParsers[sentenceType] = parser
	return nil
}

func Parse(raw string) (Sentence, error) {
	s, err := parseSentence(raw)
	if err != nil {
		return nil, err
	}
	if parser, ok := customParsers[s.Type]; ok {
		return parser(s)
	}
	if strings.HasPrefix(s.Raw, SentenceStart) {
		switch s.Type {
		case TypeRMC:
			return newRMC(s)
		case TypeGSV:
			return newGSV(s)
		case TypeGSA:
			return newGSA(s)
		case TypeGGA:
			return newGGA(s)
		case TypeVTG:
			return newVTG(s)
		case TypeGNS:
			return newGNS(s)
		}
	}
	return nil, fmt.Errorf("nmea: sentence prefix '%s' not supported", s.Prefix())
}

type TagBlock struct {
	Time         int64
	RelativeTime int64
	Destination  string
	Grouping     string
	LineCount    int64
	Source       string
	Text         string
}

func parseInt64(raw string) (int64, error) {
	i, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("nmea: tagblock unable to parse uint64 [%s]", raw)
	}
	return i, nil
}

func parseTagBlock(tags string) (TagBlock, error) {
	sumSepIndex := strings.Index(tags, ChecksumSep)
	if sumSepIndex == -1 {
		return TagBlock{}, fmt.Errorf("nmea: tagblock does not contain checksum separator")
	}
	var (
		fieldsRaw   = tags[0:sumSepIndex]
		checksumRaw = strings.ToUpper(tags[sumSepIndex+1:])
		checksum    = Checksum(fieldsRaw)
		tagBlock    TagBlock
		err         error
	)
	if checksum != checksumRaw {
		return TagBlock{}, fmt.Errorf("nmea: tagblock checksum mismatch [%s != %s]", checksum, checksumRaw)
	}
	items := strings.Split(tags[:sumSepIndex], ",")
	for _, item := range items {
		parts := strings.SplitN(item, ":", 2)
		if len(parts) != 2 {
			return TagBlock{},
				fmt.Errorf("nmea: tagblock field is malformed (should be <key>:<value>) [%s]", item)
		}
		key, value := parts[0], parts[1]
		switch key {
		case "c":
			tagBlock.Time, err = parseInt64(value)
			if err != nil {
				return TagBlock{}, err
			}
		case "d":
			tagBlock.Destination = value
		case "g":
			tagBlock.Grouping = value
		case "n":
			tagBlock.LineCount, err = parseInt64(value)
			if err != nil {
				return TagBlock{}, err
			}
		case "r":
			tagBlock.RelativeTime, err = parseInt64(value)
			if err != nil {
				return TagBlock{}, err
			}
		case "s":
			tagBlock.Source = value
		case "t":
			tagBlock.Text = value
		}
	}
	return tagBlock, nil
}

const (
	Degrees = '\u00B0'
	Minutes = '\''
	Seconds = '"'
	Point   = '.'
	North   = "N"
	South   = "S"
	East    = "E"
	West    = "W"
)

func ParseLatLong(s string) (float64, error) {
	var l float64
	if v, err := ParseDMS(s); err == nil {
		l = v
	} else if v, err := ParseGPS(s); err == nil {
		l = v
	} else if v, err := ParseDecimal(s); err == nil {
		l = v
	} else {
		return 0, fmt.Errorf("cannot parse [%s], unknown format", s)
	}
	return l, nil
}

func ParseGPS(s string) (float64, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format: %s", s)
	}
	dir := parts[1]
	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parse error: %s", err.Error())
	}
	degrees := math.Floor(value / 100)
	minutes := value - (degrees * 100)
	value = degrees + minutes/60
	if dir == North || dir == East {
		return value, nil
	} else if dir == South || dir == West {
		return 0 - value, nil
	} else {
		return 0, fmt.Errorf("invalid direction [%s]", dir)
	}
}

func FormatGPS(l float64) string {
	padding := ""
	degrees := math.Floor(math.Abs(l))
	fraction := (math.Abs(l) - degrees) * 60
	if fraction < 10 {
		padding = "0"
	}
	return fmt.Sprintf("%d%s%.4f", int(degrees), padding, fraction)
}

func ParseDecimal(s string) (float64, error) {
	l, err := strconv.ParseFloat(s, 64)
	if err != nil || s[0] != '-' && len(strings.Split(s, ".")[0]) > 3 {
		return 0.0, fmt.Errorf("parse error (not decimal coordinate)")
	}
	return l, nil
}

func ParseDMS(_ string) (float64, error) {
	return 0, fmt.Errorf("unsuported")
}

func FormatDMS(l float64) string {
	val := math.Abs(l)
	degrees := int(math.Floor(val))
	minutes := int(math.Floor(60 * (val - float64(degrees))))
	seconds := 3600 * (val - float64(degrees) - (float64(minutes) / 60))
	return fmt.Sprintf("%d\u00B0 %d' %f\"", degrees, minutes, seconds)
}

type Time struct {
	Valid       bool
	Hour        int
	Minute      int
	Second      int
	Millisecond int
}

func (t Time) String() string {
	seconds := float64(t.Second) + float64(t.Millisecond)/1000
	return fmt.Sprintf("%02d:%02d:%07.4f", t.Hour, t.Minute, seconds)
}

func ParseTime(s string) (Time, error) {
	if s == "" {
		return Time{}, nil
	}
	hour, _ := strconv.Atoi(s[:2])
	minute, _ := strconv.Atoi(s[2:4])
	second, _ := strconv.ParseFloat(s[4:], 64)
	whole, frac := math.Modf(second)
	return Time{true, hour, minute, int(whole), int(math.Round(frac * 1000))}, nil
}

type Date struct {
	Valid bool
	DD    int
	MM    int
	YY    int
}

func (d Date) String() string {
	return fmt.Sprintf("%02d/%02d/%02d", d.DD, d.MM, d.YY)
}

func ParseDate(ddmmyy string) (Date, error) {
	if ddmmyy == "" {
		return Date{}, nil
	}
	if len(ddmmyy) != 6 {
		return Date{}, fmt.Errorf("parse date: exptected ddmmyy format, got '%s'", ddmmyy)
	}
	dd, err := strconv.Atoi(ddmmyy[0:2])
	if err != nil {
		return Date{}, fmt.Errorf("%s", ddmmyy)
	}
	mm, err := strconv.Atoi(ddmmyy[2:4])
	if err != nil {
		return Date{}, fmt.Errorf("%s", ddmmyy)
	}
	yy, err := strconv.Atoi(ddmmyy[4:6])
	if err != nil {
		return Date{}, fmt.Errorf("%s", ddmmyy)
	}
	return Date{true, dd, mm, yy}, nil
}

func LatDir(l float64) string {
	if l < 0.0 {
		return South
	}
	return North
}

func LonDir(l float64) string {
	if l < 0.0 {
		return East
	}
	return West
}

const (
	TypeVDM = "VDM"
	TypeVDO = "VDO"
)

type VDMVDO struct {
	BaseSentence
	NumFragments   int64
	FragmentNumber int64
	MessageID      int64
	Channel        string
	Payload        []byte
}

func newVDMVDO(s BaseSentence) (VDMVDO, error) {
	p := NewParser(s)
	m := VDMVDO{
		BaseSentence:   s,
		NumFragments:   p.Int64(0, "number of fragments"),
		FragmentNumber: p.Int64(1, "fragment number"),
		MessageID:      p.Int64(2, "sequence number"),
		Channel:        p.String(3, "channel ID"),
		Payload:        p.SixBitASCIIArmour(4, int(p.Int64(5, "number of padding bits")), "payload"),
	}
	return m, p.Err()
}

const (
	TypeVTG = "VTG"
)

type VTG struct {
	BaseSentence
	TrueTrack        float64
	MagneticTrack    float64
	GroundSpeedKnots float64
	GroundSpeedKPH   float64
}

func newVTG(s BaseSentence) (VTG, error) {
	p := NewParser(s)
	p.AssertType(TypeVTG)
	return VTG{
		BaseSentence:     s,
		TrueTrack:        p.Float64(0, "true track"),
		MagneticTrack:    p.Float64(2, "magnetic track"),
		GroundSpeedKnots: p.Float64(4, "ground speed (knots)"),
		GroundSpeedKPH:   p.Float64(6, "ground speed (km/h)"),
	}, p.Err()
}
