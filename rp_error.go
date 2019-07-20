package sms

import (
	"bytes"
	"fmt"
	"io"
)

// Error is RP-ERROR RPDU
type Error struct {
	MR   byte  `json:"mr"`             // M / Message Reference
	CS   byte  `json:"cs"`             // M / Cause
	Diag *byte `json:"diag,omitempty"` // O / Diagnostics
	UD   TPDU  `json:"ud,omitempty"`   // O / User Data
}

// MarshalRPMO returns binary data
func (d Error) MarshalRPMO() []byte {
	return d.marshal(4)
}

// MarshalRPMT returns binary data
func (d Error) MarshalRPMT() []byte {
	return d.marshal(5)
}

func (d Error) marshal(mti byte) []byte {
	w := new(bytes.Buffer)
	w.WriteByte(mti)
	w.WriteByte(d.MR)
	if d.Diag != nil {
		w.WriteByte(2)
		w.WriteByte(d.CS)
		w.WriteByte(*d.Diag)
	} else {
		w.WriteByte(1)
		w.WriteByte(d.CS)
	}
	if d.UD != nil {
		b := d.UD.MarshalTP()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

// UnmarshalErrorMO decode Error MO from bytes
func UnmarshalErrorMO(b []byte) (a Error, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *Error) UnmarshalRPMO(b []byte) error {
	ud, e := d.unmarshal(b, 4)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = UnmarshalMOTP(ud)
	}
	return e
}

// UnmarshalErrorMT decode Error MO from bytes
func UnmarshalErrorMT(b []byte) (a Error, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *Error) UnmarshalRPMT(b []byte) error {
	ud, e := d.unmarshal(b, 5)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = UnmarshalMTTP(ud)
	}
	return e
}

func (d *Error) unmarshal(b []byte, mti byte) ([]byte, error) {
	if d == nil {
		return nil, fmt.Errorf("nil data")
	}
	if len(b) < 4 {
		return nil, io.EOF
	}
	if b[0] != mti {
		return nil, fmt.Errorf("invalid data")
	}
	d.MR = b[1]
	d.CS = b[3]
	if b[2] == 2 {
		tmp := b[4]
		d.Diag = &tmp
		return readOptionalUD(b[5:])
	}
	return readOptionalUD(b[4:])
}

func (d Error) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Error\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-CS:   cause=%d, diagnostic=%d\n", Indent, d.CS, d.Diag)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}
