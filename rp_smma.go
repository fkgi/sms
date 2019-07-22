package sms

import (
	"bytes"
	"fmt"
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
	r := bytes.NewReader(b)
	var e error

	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 6 {
		return UnexpectedMessageTypeError{Expected: 6, Actual: b[0]}
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return e
	}
	if r.Len() != 0 {
		return InvalidLengthError{}
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
