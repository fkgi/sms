package sms

import (
	"bytes"
	"fmt"
)

func cpCauseStat(c byte) string {
	switch c {
	case 17:
		return "Network failure"
	case 22:
		return "Congestion"
	case 81:
		return "Invalid Transaction Identifier value"
	case 95:
		return "Semantically incorrect message"
	case 96:
		return "Invalid mandatory information"
	case 97:
		return "Message type non existent or not implemented"
	case 98:
		return "Message not compatible with the short message protocol state"
	case 99:
		return "Information element non existent or not implemented"
	case 111:
		return "Protocol error, unspecified"
	}
	return fmt.Sprintf("Unspecified(%d)", c)
}

// Error is CP-ERROR CPDU
type Error struct {
	TI byte `json:"ti"` // M / Transaction identifier
	CS byte `json:"cs"` // M / Cause
}

// MarshalCP returns binary data
func (d Error) MarshalCP() []byte {
	b := make([]byte, 3)
	b[0] = (d.TI & 0x0f) << 4
	b[0] |= 0x09
	b[1] = 0x10
	b[2] = d.CS
	return b
}

// UnmarshalError decode Ack MT from bytes
func UnmarshalError(b []byte) (a Error, e error) {
	e = a.UnmarshalCP(b)
	return
}

// UnmarshalCP reads binary data
func (d *Error) UnmarshalCP(b []byte) error {
	if len(b) != 3 {
		return InvalidLengthError{}
	}

	if b[0]&0x0f != 0x09 {
		return UnexpectedMessageTypeError{
			Expected: 0x09, Actual: b[0] & 0x0f}
	}
	d.TI = b[0] >> 4
	d.TI &= 0x0f

	if b[1] != 0x10 {
		return UnexpectedMessageTypeError{
			Expected: 0x10, Actual: b[1]}
	}

	d.CS = b[2]
	return nil
}

func (d Error) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "CP-Error\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sCP-CS:   %s\n", Indent, cpCauseStat(d.CS))

	return w.String()
}
