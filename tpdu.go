package sms

import (
	"bytes"
	"io"
	"time"
)

var (
	msgRef byte
	// Indent for String() output for each TPDU
	Indent = " | "
)

func init() {
	msgRef = byte(time.Now().Nanosecond())
}

// TPDU represents a SMS PDU
type TPDU interface {
	Encode() []byte
	Decode([]byte) error
	String() string
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

func readUD(r *bytes.Reader, d dcs, h bool) ([]byte, []udh, error) {
	p, e := r.ReadByte()
	if e != nil {
		return nil, nil, e
	}
	l := d.unitSize()
	l *= int(p)
	if l%8 != 0 {
		l += 8 - l%8
	}

	ud := make([]byte, l/8)
	if r.Len() < len(ud) {
		return nil, nil, io.EOF
	}
	r.Read(ud)

	if h {
		return ud[ud[0]+1:], decodeUDH(ud[0 : ud[0]+1]), nil
	}
	return ud, nil, nil
}

func writeUD(w *bytes.Buffer, ud []byte, h []udh, d dcs) {
	udh := encodeUDH(h)

	u := d.unitSize()
	l := len(udh) + len(ud)
	l = ((l * 8) - (l * 8 % u)) / u

	w.WriteByte(byte(l))
	w.Write(udh)
	w.Write(ud)
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
		time.FixedZone("unknown", l*15*60))
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
