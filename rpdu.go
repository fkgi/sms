package sms

import (
	"fmt"
)

// RPDU represents a SMS RP PDU
type RPDU interface {
	EncodeMO() []byte
	EncodeMT() []byte
	DecodeMO([]byte) error
	DecodeMT([]byte) error
	fmt.Stringer
}
