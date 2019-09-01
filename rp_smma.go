package sms

import (
	"bytes"
	"fmt"
)

// MemoryAvailable is RP-SMMA RPDU
type MemoryAvailable struct {
	cpData

	RMR byte `json:"rmr"` // M / Message Reference
}

// MarshalRP output byte data of this RPDU
func (d MemoryAvailable) MarshalRP() []byte {
	return []byte{6, d.RMR}
}

// MarshalCP output byte data of this CPDU
func (d MemoryAvailable) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalMemoryAvailable decode MemoryAvailable MO from bytes
func UnmarshalMemoryAvailable(b []byte) (a MemoryAvailable, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *MemoryAvailable) UnmarshalRP(b []byte) (e error) {
	r := bytes.NewReader(b)
	var tmp byte

	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 6 {
		e = UnexpectedMessageTypeError{
			Expected: 6, Actual: b[0]}
		return
	}
	if d.RMR, e = r.ReadByte(); e == nil && r.Len() != 0 {
		e = ExtraDataError{}
	}

	return
}

// UnmarshalCP get data of this CPDU
func (d *MemoryAvailable) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
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
