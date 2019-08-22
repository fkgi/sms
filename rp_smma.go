package sms

import (
	"bytes"
	"fmt"
)

// MemoryAvailable is RP-SMMA RPDU
type MemoryAvailable struct {
	TI  byte `json:"ti"`  // M / Transaction identifier
	RMR byte `json:"rmr"` // M / Message Reference
}

// MarshalRP returns binary data
func (d MemoryAvailable) MarshalRP() []byte {
	return []byte{6, d.RMR}
}

// MarshalCP output byte data of this CPDU
func (d MemoryAvailable) MarshalCP() []byte {
	return marshalCpDataWith(d.TI, d.MarshalRP())
}

// UnmarshalMemoryAvailable decode MemoryAvailable MO from bytes
func UnmarshalMemoryAvailable(b []byte) (a MemoryAvailable, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP reads binary data
func (d *MemoryAvailable) UnmarshalRP(b []byte) error {
	r := bytes.NewReader(b)
	var e error

	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 6 {
		return UnexpectedMessageTypeError{Expected: 6, Actual: b[0]}
	}
	if d.RMR, e = r.ReadByte(); e != nil {
		return e
	}
	if r.Len() != 0 {
		return InvalidLengthError{}
	}
	return nil
}

// UnmarshalCP get data of this CPDU
func (d *MemoryAvailable) UnmarshalCP(b []byte) (e error) {
	d.TI, b, e = unmarshalCpDataWith(b)
	if e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

func (d MemoryAvailable) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "RP-MemoryAvailable\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)

	return w.String()
}
