package sms

import (
	"bytes"
	"fmt"
)

// Ack is RPDU ack answer message
type Ack struct {
	MR byte
	UD TPDU
}

func (d *Ack) encode(mo bool) []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)
	if mo {
		w.WriteByte(2)
	} else {
		w.WriteByte(3)
	}
	w.WriteByte(d.MR)

	if d.UD != nil {
		b := d.UD.Encode()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

func (d *Ack) EncodeMO() []byte {
	return d.encode(true)
}

func (d *Ack) EncodeMT() []byte {
	return d.encode(false)
}

func (d *Ack) DecodeMO(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	if b[0] != 2 {
		return fmt.Errorf("not MO Ack data")
	}
	d.MR = b[1]
	tpdu, e := DecodeAsSC(b[2:])
	d.UD = tpdu
	return e
}

func (d *Ack) DecodeMT(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	if b[0] != 3 {
		return fmt.Errorf("not MT Ack data")
	}
	d.MR = b[1]
	tpdu, e := DecodeAsMS(b[2:])
	d.UD = tpdu
	return e
}

func (d *Ack) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Ack\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}
