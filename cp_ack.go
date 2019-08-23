package sms

import (
	"bytes"
	"fmt"
)

// Ack is CP-ACK CPDU
type Ack struct {
	TI byte `json:"ti"` // M / Transaction identifier
}

// MarshalCP returns binary data
func (d Ack) MarshalCP() []byte {
	b := make([]byte, 2)
	b[0] = (d.TI & 0x0f) << 4
	b[0] |= 0x09
	b[1] = 0x04
	return b
}

// UnmarshalAck decode Ack MT from bytes
func UnmarshalAck(b []byte) (a Ack, e error) {
	e = a.UnmarshalCP(b)
	return
}

// UnmarshalCP reads binary data
func (d *Ack) UnmarshalCP(b []byte) (e error) {
	if len(b) != 2 {
		e = InvalidLengthError{}
	} else {
		d.TI, e = unmarshalCpHeader(0x04, b)
	}
	return
}

func (d Ack) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "CP-Ack\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)

	return w.String()
}
