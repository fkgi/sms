package sms

import "fmt"

// RPDU represents a SMS RP PDU
type RPDU interface {
	Encode() []byte
	Decode([]byte) error
	fmt.Stringer
}
