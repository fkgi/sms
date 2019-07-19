package sms

import (
	"bytes"
	"fmt"
	"io"
)

// MemoryAvailable is RP-SMMA RPDU
type MemoryAvailable struct {
	MR byte // M / Message Reference
}

// EncodeMO returns binary data
func (d MemoryAvailable) EncodeMO() []byte {
	return []byte{6, d.MR}
}

// EncodeMT returns binary data
func (d *MemoryAvailable) EncodeMT() []byte {
	return []byte{}
}

// DecodeMO reads binary data
func (d *MemoryAvailable) DecodeMO(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	if len(b) < 2 {
		return io.EOF
	}
	if b[0] != 6 {
		return fmt.Errorf("invalid data")
	}
	d.MR = b[1]
	return nil
}

// DecodeMT reads binary data
func (d *MemoryAvailable) DecodeMT(b []byte) error {
	return fmt.Errorf("invalid data")
}

func (d MemoryAvailable) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MemoryAvailable\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	return w.String()
}
