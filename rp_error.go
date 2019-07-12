package sms

import (
	"bytes"
	"fmt"
)

// Error is RPDU nack answer message
type Error struct {
	MR   byte
	CS   byte
	Diag byte
	UD   TPDU
}

func (d *Error) encode(mo bool) []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)
	if mo {
		w.WriteByte(4)
	} else {
		w.WriteByte(5)
	}
	w.WriteByte(d.MR)

	// w.WriteByte(IEI)
	w.WriteByte(2)
	w.WriteByte(d.CS)
	w.WriteByte(d.Diag)

	if d.UD != nil {
		b := d.UD.Encode()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

func (d *Error) EncodeMO() []byte {
	return d.encode(true)
}

func (d *Error) EncodeMT() []byte {
	return d.encode(false)
}

func (d *Error) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Error\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-CS:   cause=%d, diagnostic=%d\n", Indent, d.CS, d.Diag)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}
