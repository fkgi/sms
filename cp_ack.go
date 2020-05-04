package sms

import (
	"bytes"
	"fmt"
	"io"
)

// CpAck is CP-ACK CPDU
type CpAck struct {
	TI byte `json:"ti"` // M / Transaction identifier
}

// MarshalCP output byte data of this CPDU
func (d CpAck) MarshalCP() []byte {
	b := make([]byte, 2)
	b[0] = (d.TI & 0x0f) << 4
	b[0] |= 0x09
	b[1] = 0x04
	return b
}

// UnmarshalAck decode Ack MT from bytes
func UnmarshalAck(b []byte) (a CpAck, e error) {
	e = a.UnmarshalCP(b)
	return
}

// UnmarshalCP get data of this CPDU
func (d *CpAck) UnmarshalCP(b []byte) (e error) {
	if len(b) < 2 {
		e = io.EOF
	} else if len(b) > 2 {
		e = ErrInvalidLength
	} else {
		d.TI, e = unmarshalCpHeader(0x04, b)
	}
	return
}

func (d CpAck) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "CP-Ack\n")
	fmt.Fprintf(w, "%sCP-TI: %s\n", Indent, cpTIStat(d.TI))

	return w.String()
}
