package sms

import (
	"bytes"
	"fmt"
	"io"
)

// MemoryAvailable is RP-SMMA RPDU
type MemoryAvailable struct {
	MR byte `json:"mr"` // M / Message Reference
}

// MarshalRPMO returns binary data
func (d MemoryAvailable) MarshalRPMO() []byte {
	return []byte{6, d.MR}
}

// MarshalRPMT returns binary data
func (d MemoryAvailable) MarshalRPMT() []byte {
	return []byte{}
}

// UnmarshalMemoryAvailableMO decode MemoryAvailable MO from bytes
func UnmarshalMemoryAvailableMO(b []byte) (a MemoryAvailable, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *MemoryAvailable) UnmarshalRPMO(b []byte) error {
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

// UnmarshalMemoryAvailableMT decode MemoryAvailable MO from bytes
func UnmarshalMemoryAvailableMT(b []byte) (a MemoryAvailable, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *MemoryAvailable) UnmarshalRPMT(b []byte) error {
	return fmt.Errorf("invalid data")
}

func (d MemoryAvailable) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MemoryAvailable\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	return w.String()
}
