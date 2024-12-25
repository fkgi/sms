package sms

import (
	"bytes"
	"fmt"
	"io"
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
	TI byte `json:"cp-ti"` // M / Transaction identifier
	CS byte `json:"cp-cs"` // M / Cause
}

// MarshalCP output byte data of this CPDU
func (d CpError) MarshalCP() []byte {
	b := make([]byte, 3)
	b[0] = (d.TI & 0x0f) << 4
	b[0] |= 0x09
	b[1] = 0x10
	b[2] = d.CS
	return b
}

// UnmarshalError decode Ack MT from bytes
func UnmarshalError(b []byte) (a CpError, e error) {
	e = a.UnmarshalCP(b)
	return
}

// UnmarshalCP get data of this CPDU
func (d *CpError) UnmarshalCP(b []byte) (e error) {
	if len(b) < 3 {
		e = io.EOF
	} else if len(b) > 3 {
		e = ErrInvalidLength
	} else {
		d.TI, e = unmarshalCpHeader(0x10, b)
		d.CS = b[2]
	}
	return
}

func (d CpError) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "CP-Error\n")
	fmt.Fprintf(w, "%sCP-TI: %s\n", Indent, cpTIStat(d.TI))
	fmt.Fprintf(w, "%sCP-CS: %s\n", Indent, cpCauseStat(d.CS))

	return w.String()
}

func (d CpError) Error() string {
	return "CP error, cause: " + cpCauseStat(d.CS)
}
