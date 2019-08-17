package sms

import (
	"bytes"
	"fmt"
)

// CpAck is CP-ACK CPDU
type CpAck struct {
	TI byte `json:"ti"` // M / Transaction identifier
}

// MarshalCPMO returns binary data
func (d CpAck) MarshalCPMO() []byte {
	return d.marshal()
}

// MarshalCPMT returns binary data
func (d CpAck) MarshalCPMT() []byte {
	return d.marshal()
}

func (d CpAck) marshal() []byte {
	b := make([]byte, 2)
	b[0] = (d.TI & 0x0f) << 4
	b[0] |= 0x09
	b[1] = 0x04
	return b
}

// UnmarshalCpAckMO decode Ack MT from bytes
func UnmarshalCpAckMO(b []byte) (a CpAck, e error) {
	e = a.UnmarshalCPMO(b)
	return
}

// UnmarshalCPMO reads binary data
func (d *CpAck) UnmarshalCPMO(b []byte) error {
	return d.unmarshal(b)
}

// UnmarshalCpAckMT decode Ack MT from bytes
func UnmarshalCpAckMT(b []byte) (a CpAck, e error) {
	e = a.UnmarshalCPMT(b)
	return
}

// UnmarshalCPMT reads binary data
func (d *CpAck) UnmarshalCPMT(b []byte) error {
	return d.unmarshal(b)
}

func (d *CpAck) unmarshal(b []byte) error {
	if len(b) != 2 {
		return InvalidLengthError{}
	}

	if b[0]&0x0f != 0x09 {
		return UnexpectedMessageTypeError{
			Expected: 0x09, Actual: b[0] & 0x0f}
	}
	d.TI = b[0] >> 4
	d.TI &= 0x0f

	if b[1] != 0x04 {
		return UnexpectedMessageTypeError{
			Expected: 0x04, Actual: b[1]}
	}
	return nil
}

func (d CpAck) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: CP-Ack\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)

	return w.String()
}
