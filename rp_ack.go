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
	e = a.UnmarshalMORP(b)
	return
}

// UnmarshalMORP reads binary data
func (d *Ack) UnmarshalMORP(b []byte) error {
	ud, e := d.unmarshal(b, 2)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalMOTP(ud)
	}
	return e
}

// UnmarshalAckMT decode Ack MT from bytes
func UnmarshalAckMT(b []byte) (a Ack, e error) {
	e = a.UnmarshalMTRP(b)
	return
}

// UnmarshalMTRP reads binary data
func (d *Ack) UnmarshalMTRP(b []byte) error {
	ud, e := d.unmarshal(b, 3)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalMTTP(ud)
	}
	return e
}

func (d *Ack) unmarshal(b []byte, mti byte) ([]byte, error) {
	r := bytes.NewReader(b)
	var e error

	if tmp, e := r.ReadByte(); e != nil {
		return nil, e
	} else if tmp != mti {
		return nil, UnexpectedMessageTypeError{Expected: mti, Actual: b[0]}
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return nil, e
	}
	if tmp, e := r.ReadByte(); e == io.EOF {
		return nil, nil
	} else if tmp != 41 {
		return nil, UnexpectedInformationElementError{Expected: 41, Actual: tmp}
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

func (d Ack) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: Ack\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	}

	return w.String()
}
