package sms

import (
	"fmt"
	"io"
)

// RPDU represents a SMS RP PDU
type RPDU interface {
	MarshalRPMO() []byte
	MarshalRPMT() []byte
	//DecodeMO([]byte) error
	//DecodeMT([]byte) error
	fmt.Stringer
	//json.Unmarshaler
	//json.Marshaler
}

// UnmarshalMORP parse byte data to TPDU as SC.
func UnmarshalMORP(b []byte) (t RPDU, e error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	switch b[0] & 0x03 {
	case 0x00:
		return UnmarshalDataMO(b)
	case 0x02:
		return UnmarshalAckMO(b)
	case 0x04:
		return UnmarshalErrorMO(b)
	case 0x06:
		return UnmarshalMemoryAvailableMO(b)
	}
	return nil, &InvalidDataError{
		Name:  "reserved RPDU type",
		Bytes: b}
}

// UnmarshalMTRP parse byte data to TPDU as MS.
func UnmarshalMTRP(b []byte) (t RPDU, e error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	switch b[0] & 0x03 {
	case 0x00:
		return UnmarshalDataMT(b)
	case 0x02:
		return UnmarshalAckMT(b)
	case 0x04:
		return UnmarshalErrorMT(b)
	case 0x06:
		return UnmarshalMemoryAvailableMT(b)
	}
	return nil, &InvalidDataError{
		Name:  "reserved RPDU type",
		Bytes: b}
}
