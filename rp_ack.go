package sms

import (
	"bytes"
	"fmt"
	"io"
)

// AckMO is MO RP-ACK RPDU
type AckMO struct {
	MR byte
	UD TPDU
}

// Encode returns binary data
func (d *AckMO) Encode() []byte {
	if d == nil {
		return []byte{}
	}
	return encodeAck(2, d.MR, d.UD)
}

// Decode reads binary data
func (d *AckMO) Decode(b []byte) (e error) {
	if d == nil {
		e = fmt.Errorf("nil data")
	} else {
		d.MR, d.UD, e = decodeAck(b, 2)
	}
	return
}

func (d *AckMO) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MO Ack\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}

// AckMT is MT RP-ACK RPDU
type AckMT struct {
	MR byte
	UD TPDU
}

// Encode returns binary data
func (d *AckMT) Encode() []byte {
	if d == nil {
		return []byte{}
	}
	return encodeAck(3, d.MR, d.UD)
}

// Decode reads binary data
func (d *AckMT) Decode(b []byte) (e error) {
	if d == nil {
		e = fmt.Errorf("nil data")
	} else {
		d.MR, d.UD, e = decodeAck(b, 3)
	}
	return
}

func (d *AckMT) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MT Ack\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}

func encodeAck(mti, mr byte, ud TPDU) []byte {
	w := new(bytes.Buffer)
	w.WriteByte(mti)
	w.WriteByte(mr)
	if ud != nil {
		b := ud.Encode()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

func decodeAck(b []byte, mti byte) (mr byte, ud TPDU, e error) {
	if len(b) < 2 {
		e = io.EOF
	} else if b[0] != mti {
		e = fmt.Errorf("invalid data")
	} else {
		mr = b[1]
		ud, e = DecodeAsSC(b[2:])
	}
	return
}
