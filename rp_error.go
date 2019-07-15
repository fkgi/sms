package sms

import (
	"bytes"
	"fmt"
	"io"
)

// ErrorMO is RP-ERROR RPDU
type ErrorMO struct {
	MR   byte
	CS   byte
	Diag byte
	UD   TPDU
}

// Encode returns binary data
func (d *ErrorMO) Encode() []byte {
	if d == nil {
		return []byte{}
	}
	return encodeError(4, d.MR, d.CS, d.Diag, d.UD)
}

// Decode reads binary data
func (d *ErrorMO) Decode(b []byte) (e error) {
	if d == nil {
		e = fmt.Errorf("nil data")
	} else {
		d.MR, d.CS, d.Diag, d.UD, e = decodeError(b, 4)
	}
	return
}

func (d *ErrorMO) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MO Error\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-CS:   cause=%d, diagnostic=%d\n", Indent, d.CS, d.Diag)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}

// ErrorMT is RP-ERROR RPDU
type ErrorMT struct {
	MR   byte
	CS   byte
	Diag byte
	UD   TPDU
}

// Encode returns binary data
func (d *ErrorMT) Encode() []byte {
	if d == nil {
		return []byte{}
	}
	return encodeError(5, d.MR, d.CS, d.Diag, d.UD)
}

// Decode reads binary data
func (d *ErrorMT) Decode(b []byte) (e error) {
	if d == nil {
		e = fmt.Errorf("nil data")
	} else {
		d.MR, d.CS, d.Diag, d.UD, e = decodeError(b, 5)
	}
	return
}

func (d *ErrorMT) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MT Error\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-CS:   cause=%d, diagnostic=%d\n", Indent, d.CS, d.Diag)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}

func encodeError(mti, mr, cs, diag byte, ud TPDU) []byte {
	w := new(bytes.Buffer)
	w.WriteByte(mti)
	w.WriteByte(mr)
	w.WriteByte(2)
	w.WriteByte(cs)
	w.WriteByte(diag)
	if ud != nil {
		b := ud.Encode()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

func decodeError(b []byte, mti byte) (mr, cs, diag byte, ud TPDU, e error) {
	if len(b) < 3 {
		e = io.EOF
	} else if b[0] != mti {
		e = fmt.Errorf("invalid data")
	} else {
		mr = b[1]
		cs = b[2]
		diag = b[3]
		ud, e = DecodeAsSC(b[4:])
	}
	return
}
