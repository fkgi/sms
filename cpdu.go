package sms

import (
	"fmt"
	"io"
)

// CPDU represents a SMS CP PDU
type CPDU interface {
	MarshalCPMO() []byte
	MarshalCPMT() []byte
	//Decode([]byte) error
	fmt.Stringer
	//json.Unmarshaler
	//json.Marshaler
}

// UnmarshalCPMO parse byte data to CPDU.
func UnmarshalCPMO(b []byte) (t CPDU, e error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	if b[0]&0x0f != 0x09 {
		return nil, UnexpectedMessageTypeError{
			Expected: 0x09, Actual: b[0] & 0x0f}
	}
	if len(b) < 2 {
		return nil, io.EOF
	}
	switch b[1] {
	case 0x01:
		return UnmarshalCpDataMO(b)
	case 0x04:
		return UnmarshalCpAckMO(b)
	case 0x10:
		return UnmarshalCpErrorMO(b)
	}
	return nil, UnexpectedMessageTypeError{Actual: b[1]}
}

// UnmarshalCPMT parse byte data to CPDU.
func UnmarshalCPMT(b []byte) (t CPDU, e error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	if b[0]&0x0f != 0x09 {
		return nil, UnexpectedMessageTypeError{
			Expected: 0x09, Actual: b[0] & 0x0f}
	}
	if len(b) < 2 {
		return nil, io.EOF
	}
	switch b[1] {
	case 0x01:
		return UnmarshalCpDataMT(b)
	case 0x04:
		return UnmarshalCpAckMT(b)
	case 0x10:
		return UnmarshalCpErrorMT(b)
	}
	return nil, UnexpectedMessageTypeError{Actual: b[1]}
}
