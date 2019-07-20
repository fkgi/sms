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
	MarshalTP() []byte
	json.Marshaler
	// json.Unmarshaler
	fmt.Stringer
}

// UnmarshalerTP is the interface implemented by types
// that can unmarshal a TPDU
type UnmarshalerTP interface {
	UnmarshalTP([]byte) error
}

// UnmarshalMOTP parse byte data to TPDU as SC.
func UnmarshalMOTP(b []byte) (TPDU, error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	switch b[0] & 0x03 {
	case 0x00:
		return UnmarshalDeliverReport(b)
	case 0x01:
		return UnmarshalSubmit(b)
	case 0x02:
		// return UnmarshalCommand(b)
		return nil, InvalidDataError{
			Name:  "reserved TPDU type",
			Bytes: b}
	}
	return nil, InvalidDataError{
		Name:  "reserved TPDU type",
		Bytes: b}
}

// UnmarshalMTTP parse byte data to TPDU as MS.
func UnmarshalMTTP(b []byte) (t TPDU, e error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	switch b[0] & 0x03 {
	case 0x00:
		return UnmarshalDeliver(b)
	case 0x01:
		return UnmarshalSubmitReport(b)
	case 0x02:
		return UnmarshalStatusReport(b)
	}
	return nil, InvalidDataError{
		Name:  "reserved TPDU type",
		Bytes: b}
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
