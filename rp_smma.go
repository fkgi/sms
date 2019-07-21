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
	e = a.UnmarshalMORP(b)
	return
}

// UnmarshalMORP reads binary data
func (d *MemoryAvailable) UnmarshalMORP(b []byte) error {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0] != 6 {
		return &InvalidDataError{
			Name: "invalid MTI"}
	}
	r := bytes.NewReader(b[1:])
	var e error
	if d.MR, e = r.ReadByte(); e != nil {
		return e
	}
	if r.Len() != 0 {
		return &InvalidDataError{
			Name: "extra part"}
	}
	return nil
}

// UnmarshalMemoryAvailableMT decode MemoryAvailable MO from bytes
func UnmarshalMemoryAvailableMT(b []byte) (a MemoryAvailable, e error) {
	e = a.UnmarshalMTRP(b)
	return
}

// UnmarshalMTRP reads binary data
func (d *MemoryAvailable) UnmarshalMTRP(b []byte) error {
	return fmt.Errorf("invalid data")
}

func (d MemoryAvailable) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MemoryAvailable\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	return w.String()
}
