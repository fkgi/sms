package sms

import (
	"fmt"
	"io"
)

// RPDU represents a SMS RP PDU
type RPDU interface {
	serializer
	fmt.Stringer
	//json.Unmarshaler
	//json.Marshaler
}

// serializer provids binary encode and decode
type serializer interface {
	Encode() []byte
	Decode([]byte) error
}

// DecodeRPasSC parse byte data to TPDU as SC.
func DecodeRPasSC(b []byte) (t RPDU, e error) {
	return decodeRP(b, true)
}

// DecodeRPasMS parse byte data to TPDU as MS.
func DecodeRPasMS(b []byte) (t RPDU, e error) {
	return decodeRP(b, false)
}

func decodeRP(b []byte, sc bool) (t RPDU, e error) {
	if len(b) == 0 {
		e = io.EOF
	} else if sc {
		switch b[0] & 0x03 {
		case 0x00:
			t = &DataMO{}
		case 0x02:
			t = &AckMO{}
		case 0x04:
			t = &ErrorMO{}
		case 0x06:
			t = &MemoryAvailable{}
		default:
			e = &InvalidDataError{
				Name:  "reserved RPDU type",
				Bytes: b}
		}
	} else {
		switch b[0] & 0x03 {
		case 0x01:
			t = &DataMT{}
		case 0x03:
			t = &AckMT{}
		case 0x05:
			t = &ErrorMT{}
		default:
			e = &InvalidDataError{
				Name:  "reserved RPDU type",
				Bytes: b}
		}
	}

	if e == nil {
		e = t.Decode(b)
	}
	return
}
