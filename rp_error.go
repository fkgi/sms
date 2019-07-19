package sms

import (
	"bytes"
	"fmt"
	"io"
)

// Error is RP-ERROR RPDU
type Error struct {
	MR   byte  // M / Message Reference
	CS   byte  // M / Cause
	Diag *byte // M / Diagnostics
	UD   TPDU  // O / User Data
}

// EncodeMO returns binary data
func (d Error) EncodeMO() []byte {
	return d.encode(4)
}

// EncodeMT returns binary data
func (d Error) EncodeMT() []byte {
	return d.encode(5)
}

func (d Error) encode(mti byte) []byte {
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
		b := d.UD.Encode()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

// DecodeMO reads binary data
func (d *Error) DecodeMO(b []byte) error {
	ud, e := d.decode(b, 4)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = DecodeAsSC(ud)
	}
	return e
}

// DecodeMT reads binary data
func (d *Error) DecodeMT(b []byte) error {
	ud, e := d.decode(b, 5)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = DecodeAsMS(ud)
	}
	return e
}

func (d *Error) decode(b []byte, mti byte) ([]byte, error) {
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
