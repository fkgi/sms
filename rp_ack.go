package sms

import (
	"bytes"
	"fmt"
	"io"
)

// Ack is RP-ACK RPDU
type Ack struct {
	MR byte `json:"mr"`           // M / Message Reference
	UD TPDU `json:"ud,omitempty"` // O / User Data
}

// MarshalRPMO returns binary data
func (d Ack) MarshalRPMO() []byte {
	return d.marshal(2)
}

// MarshalRPMT returns binary data
func (d Ack) MarshalRPMT() []byte {
	return d.marshal(3)
}

func (d Ack) marshal(mti byte) []byte {
	w := new(bytes.Buffer)
	w.WriteByte(mti)
	w.WriteByte(d.MR)
	if d.UD != nil {
		b := d.UD.MarshalTP()
		w.WriteByte(41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}
	return w.Bytes()
}

// UnmarshalAckMO decode Ack MO from bytes
func UnmarshalAckMO(b []byte) (a Ack, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *Ack) UnmarshalRPMO(b []byte) error {
	ud, e := d.unmarshal(b, 2)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = UnmarshalMOTP(ud)
	}
	return e
}

// UnmarshalAckMT decode Ack MT from bytes
func UnmarshalAckMT(b []byte) (a Ack, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *Ack) UnmarshalRPMT(b []byte) error {
	ud, e := d.unmarshal(b, 3)
	if e != nil {
		return e
	}
	if ud != nil {
		d.UD, e = UnmarshalMTTP(ud)
	}
	return e
}

func (d *Ack) unmarshal(b []byte, mti byte) ([]byte, error) {
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

func (d Ack) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Ack\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}
	return w.String()
}
