package sms

import (
	"bytes"
	"fmt"
	"io"
)

// MemoryAvailable is RP-SMMA RPDU
type MemoryAvailable struct {
	MR byte
}

// Encode returns binary data
func (d *MemoryAvailable) Encode() []byte {
	if d == nil {
		return []byte{}
	}
	return []byte{6, d.MR}
}

// Decode reads binary data
func (d *MemoryAvailable) Decode(b []byte) (e error) {
	if d == nil {
		e = fmt.Errorf("nil data")
	} else if len(b) < 2 {
		e = io.EOF
	} else if b[0] != 6 {
		e = fmt.Errorf("invalid data")
	} else {
		d.MR = b[1]
	}
	return
}

func (d *MemoryAvailable) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MemoryAvailable\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	return w.String()
}
