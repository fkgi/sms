package sms

import (
	"bytes"
	"fmt"
)

// Data is RP-DATA RPDU
type Data struct {
	MR byte    `json:"mr"`           // M / Message Reference
	OA Address `json:"oa,omitempty"` // C / Originator Address
	DA Address `json:"da,omitempty"` // C / Recipient Address
	UD TPDU    `json:"ud,omitempty"` // M / User Data
}

// MarshalRPMO returns binary data
func (d Data) MarshalRPMO() []byte {
	w := new(bytes.Buffer)
	w.WriteByte(0) // MTI
	w.WriteByte(d.MR)
	w.WriteByte(0) // OA
	l, a := d.DA.Encode()
	w.WriteByte(l)
	w.Write(a)
	b := d.UD.MarshalTP()
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// MarshalRPMT returns binary data
func (d Data) MarshalRPMT() []byte {
	w := new(bytes.Buffer)
	w.WriteByte(1) // MTI
	w.WriteByte(d.MR)
	l, a := d.OA.Encode()
	w.WriteByte(l)
	w.Write(a)
	w.WriteByte(0) // DA
	b := d.UD.MarshalTP()
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// UnmarshalDataMO decode Data MO from bytes
func UnmarshalDataMO(b []byte) (a Data, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *Data) UnmarshalRPMO(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	r := bytes.NewReader(b)
	var e error
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return e
	}
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if d.DA, e = readAddr(r); e != nil {
		return e
	}
	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return e
	}
	if _, e := r.Read(b); e != nil {
		return e
	}
	d.UD, e = UnmarshalMOTP(b)
	return e
}

// UnmarshalDataMT decode Data MO from bytes
func UnmarshalDataMT(b []byte) (a Data, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *Data) UnmarshalRPMT(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	r := bytes.NewReader(b)
	var e error
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 1 {
		return fmt.Errorf("invalid data")
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return e
	}
	if d.OA, e = readAddr(r); e != nil {
		return e
	}
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return e
	}
	if _, e = r.Read(b); e != nil {
		return e
	}
	d.UD, e = UnmarshalMTTP(b)
	return e
}

func (d Data) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Data\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-OA:   %s\n", Indent, d.OA)
	fmt.Fprintf(w, "%sRP-DA:   %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	return w.String()
}
