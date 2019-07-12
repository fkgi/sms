package sms

import (
	"bytes"
	"fmt"
)

// MemoryAvailable is RP-SMMA RPDU
type MemoryAvailable struct {
	MR byte
}

func (d *MemoryAvailable) EncodeMO() []byte {
	if d == nil {
		return []byte{}
	}
	return []byte{6, d.MR}
}

func (d *MemoryAvailable) EncodeMT() []byte {
	return []byte{}
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
