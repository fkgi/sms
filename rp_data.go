package sms

import (
	"bytes"
	"fmt"
	"io"
)

// RpData is RP-DATA RPDU
type RpData struct {
	MR byte    `json:"mr"`           // M / Message Reference
	OA Address `json:"oa,omitempty"` // C / Originator Address
	DA Address `json:"da,omitempty"` // C / Recipient Address
	UD TPDU    `json:"ud,omitempty"` // M / User Data
}

// MarshalRPMO returns binary data
func (d RpData) MarshalRPMO() []byte {
	w := new(bytes.Buffer)

	w.WriteByte(0) // MTI
	w.WriteByte(d.MR)
	w.WriteByte(0) // OA
	_, a := d.DA.Marshal()
	w.WriteByte(byte(len(a)))
	w.Write(a)
	b := d.UD.MarshalTP()
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// MarshalRPMT returns binary data
func (d RpData) MarshalRPMT() []byte {
	w := new(bytes.Buffer)

	w.WriteByte(1) // MTI
	w.WriteByte(d.MR)
	_, a := d.OA.Marshal()
	w.WriteByte(byte(len(a)))
	w.Write(a)
	w.WriteByte(0) // DA
	b := d.UD.MarshalTP()
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// UnmarshalDataMO decode Data MO from bytes
func UnmarshalDataMO(b []byte) (a RpData, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *RpData) UnmarshalRPMO(b []byte) error {
	ud, e := d.unmarshal(b, 0)
	if e != nil {
		return e
	}
	d.UD, e = UnmarshalTPMO(ud)
	return e
}

// UnmarshalDataMT decode Data MO from bytes
func UnmarshalDataMT(b []byte) (a RpData, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *RpData) UnmarshalRPMT(b []byte) error {
	ud, e := d.unmarshal(b, 1)
	if e != nil {
		return e
	}
	d.UD, e = UnmarshalTPMT(ud)
	return e
}

func (d *RpData) unmarshal(b []byte, mti byte) ([]byte, error) {
	r := bytes.NewReader(b)
	var e error

	if tmp, e := r.ReadByte(); e != nil {
		return nil, e
	} else if tmp != mti {
		return nil, UnexpectedMessageTypeError{Expected: mti, Actual: tmp}
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return nil, e
	}
	if d.OA, e = readRPAddr(r); e != nil {
		return nil, e
	}
	if d.DA, e = readRPAddr(r); e != nil {
		return nil, e
	}
	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return nil, e
	}
	if n, e := r.Read(b); e != nil {
		return nil, e
	} else if n != len(b) {
		return nil, io.EOF
	}
	if r.Len() != 0 {
		return nil, InvalidLengthError{}
	}
	return b, nil
}

func (d RpData) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: RP-Data\n")
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-OA: %s\n", Indent, d.OA)
	fmt.Fprintf(w, "%sRP-DA: %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sRP-UD: %s\n", Indent, d.UD)

	return w.String()
}
