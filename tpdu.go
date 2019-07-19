package sms

import (
	"bytes"
	"encoding/json"
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
	json.Unmarshaler
	json.Marshaler
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
			e = InvalidDataError{
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
			e = InvalidDataError{
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
