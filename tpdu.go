package sms

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

var (
	msgRef chan byte
	// Indent for String() output for each TPDU
	Indent = " | "
)

func init() {
	msgRef = make(chan byte, 1)
	msgRef <- byte(time.Now().Nanosecond())
}

// NextMsgReference make Message Reference ID
func NextMsgReference() byte {
	ret := <-msgRef
	msgRef <- ret + 1
	return ret
}

// TPDU represents a SMS TP PDU
type TPDU interface {
	Encode() []byte
	Decode([]byte) error
	fmt.Stringer
}

// DecodeAsSC parse byte data to TPDU as SC.
func DecodeAsSC(b []byte) (t TPDU, e error) {
	return decode(b, true)
}

// DecodeAsMS parse byte data to TPDU as MS.
func DecodeAsMS(b []byte) (t TPDU, e error) {
	return decode(b, false)
}

func decode(b []byte, sc bool) (t TPDU, e error) {
	if len(b) == 0 {
		e = io.EOF
	} else if sc {
		switch b[0] & 0x03 {
		case 0x00:
			t = &DeliverReport{}
		case 0x01:
			t = &Submit{}
		case 0x02:
			// t = &Command{}
		case 0x03:
			e = &InvalidDataError{
				Name:  "reserved TPDU type",
				Bytes: b}
		}
	} else {
		switch b[0] & 0x03 {
		case 0x00:
			t = &Deliver{}
		case 0x01:
			t = &SubmitReport{}
		case 0x02:
			t = &StatusReport{}
		case 0x03:
			e = &InvalidDataError{
				Name:  "reserved TPDU type",
				Bytes: b}
		}
	}

	if e == nil {
		e = t.Decode(b)
	}
	return
}

func read7Bytes(r *bytes.Reader) ([7]byte, error) {
	if r.Len() < 7 {
		return [7]byte{}, io.EOF
	}
	b := make([]byte, 7)
	r.Read(b)
	return [7]byte{
		b[0], b[1], b[2], b[3], b[4], b[5], b[6]}, nil
}

func encodeSCTimeStamp(t time.Time) (r []byte) {
	r = make([]byte, 7)
	r[0] = int2SemiOctet(t.Year())
	r[1] = int2SemiOctet(int(t.Month()))
	r[2] = int2SemiOctet(t.Day())
	r[3] = int2SemiOctet(t.Hour())
	r[4] = int2SemiOctet(t.Minute())
	r[5] = int2SemiOctet(t.Second())

	_, z := t.Zone()
	z /= 900
	r[6] = byte((z % 10) & 0x0f)
	r[6] = (r[6] << 4) | byte(((z/10)%10)&0x0f)
	if z < 0 {
		r[6] = r[6] | 0x08
	}
	return
}

func decodeSCTimeStamp(t [7]byte) time.Time {
	d := [6]int{}
	for i := range d {
		d[i] = semiOctet2Int(t[i])
	}
	l := semiOctet2Int(t[6] & 0x7f)
	if t[6]&0x80 == 0x80 {
		l = -l
	}
	return time.Date(2000+d[0],
		time.Month(d[1]), d[2], d[3], d[4], d[5], 0,
		time.FixedZone("", l*15*60)).Local()
}

func int2SemiOctet(i int) (b byte) {
	b = byte(i % 10)
	b = (b << 4) | byte((i/10)%10)
	return
}

func semiOctet2Int(b byte) (i int) {
	i = int(b & 0x0f)
	i = (i * 10) + int((b&0xf0)>>4)
	return
}
