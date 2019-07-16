package sms

import (
	"fmt"
	"io"
)

// RPDU represents a SMS RP PDU
type RPDU interface {
	EncodeMO() []byte
	EncodeMT() []byte
	DecodeMO([]byte) error
	DecodeMT([]byte) error
	fmt.Stringer
	//json.Unmarshaler
	//json.Marshaler
}

// DecodeMORP parse byte data to TPDU as SC.
func DecodeMORP(b []byte) (t RPDU, e error) {
	return decodeRP(b, true)
}

// DecodeMTRP parse byte data to TPDU as MS.
func DecodeMTRP(b []byte) (t RPDU, e error) {
	return decodeRP(b, false)
}

func decodeRP(b []byte, sc bool) (t RPDU, e error) {
	if len(b) == 0 {
		e = io.EOF
	} else {
		switch b[0] & 0x03 {
		case 0x00:
			t = &Data{}
		case 0x02:
			t = &Ack{}
		case 0x04:
			t = &Error{}
		case 0x06:
			t = &MemoryAvailable{}
		default:
			e = &InvalidDataError{
				Name:  "reserved RPDU type",
				Bytes: b}
		}
		if sc {
			e = t.DecodeMO(b)
		} else {
			e = t.DecodeMT(b)
		}
	}
	return
}

func readOptionalUD(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, nil
	}
	if len(b) < 3 {
		return nil, fmt.Errorf("invalid data")
	}
	if b[0] != 41 {
		return nil, fmt.Errorf("invalid data")
	}
	if len(b) != int(b[1]+2) {
		return nil, fmt.Errorf("invalid data")
	}
	return b[2:], nil
}
