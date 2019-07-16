package sms

import (
	"bytes"
	"fmt"
	"io"
)

// Ack is RP-ACK RPDU
type Ack struct {
	MR byte // M / Message Reference
	UD TPDU // O / User Data
}

// EncodeMO returns binary data
func (d *Ack) EncodeMO() []byte {
	return d.encode(2)
}

// EncodeMT returns binary data
func (d *Ack) EncodeMT() []byte {
	return d.encode(3)
}

func (d *Ack) encode(mti byte) []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)
	w.WriteByte(mti)
	w.WriteByte(d.MR)
	if d.UD != nil {
		b := d.UD.Encode()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

// DecodeMO reads binary data
func (d *Ack) DecodeMO(b []byte) error {
	ud, e := d.decode(b, 2)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = DecodeAsSC(ud)
	}
	return e
}

// DecodeMT reads binary data
func (d *Ack) DecodeMT(b []byte) error {
	ud, e := d.decode(b, 3)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = DecodeAsMS(ud)
	}
	return e
}

func (d *Ack) decode(b []byte, mti byte) ([]byte, error) {
	if d == nil {
		return nil, fmt.Errorf("nil data")
	}
	if len(b) < 2 {
		return nil, io.EOF
	}
	if b[0] != mti {
		return nil, fmt.Errorf("invalid data")
	}
	d.MR = b[1]
	return readOptionalUD(b[2:])
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
