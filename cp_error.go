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

// CpError is CP-ERROR CPDU
type CpError struct {
	TI byte `json:"ti"` // M / Transaction identifier
	CS byte `json:"cs"` // M /Cause
}

// MarshalCPMO returns binary data
func (d CpError) MarshalCPMO() []byte {
	return d.marshal()
}

// MarshalCPMT returns binary data
func (d CpError) MarshalCPMT() []byte {
	return d.marshal()
}

func (d CpError) marshal() []byte {
	b := make([]byte, 3)
	b[0] = (d.TI & 0x0f) << 4
	b[0] |= 0x09
	b[1] = 0x10
	b[2] = d.CS
	return b
}

// UnmarshalCpErrorMO decode Ack MT from bytes
func UnmarshalCpErrorMO(b []byte) (a CpError, e error) {
	e = a.UnmarshalCPMO(b)
	return
}

// UnmarshalCPMO reads binary data
func (d *CpError) UnmarshalCPMO(b []byte) error {
	return d.unmarshal(b)
}

// UnmarshalCpErrorMT decode Ack MT from bytes
func UnmarshalCpErrorMT(b []byte) (a CpError, e error) {
	e = a.UnmarshalCPMT(b)
	return
}

// UnmarshalCPMT reads binary data
func (d *CpError) UnmarshalCPMT(b []byte) error {
	return d.unmarshal(b)
}

func (d *CpError) unmarshal(b []byte) error {
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

func (d CpError) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: CP-Data\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sCP-CS:   %s\n", Indent, cpCauseStat(d.CS))

	return w.String()
}
