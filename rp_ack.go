package sms

import (
	"bytes"
	"fmt"
	"io"
)

// RpAck is RP-ACK RPDU
type RpAck struct {
	MR byte `json:"mr"`           // M / Message Reference
	UD TPDU `json:"ud,omitempty"` // O / User Data
}

// MarshalRPMO returns binary data
func (d RpAck) MarshalRPMO() []byte {
	return d.marshal(2)
}

// MarshalRPMT returns binary data
func (d RpAck) MarshalRPMT() []byte {
	return d.marshal(3)
}

func (d RpAck) marshal(mti byte) []byte {
	w := new(bytes.Buffer)

	w.WriteByte(mti)
	w.WriteByte(d.MR)
	if d.UD != nil {
		b := d.UD.MarshalTP()
		w.WriteByte(0x41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}

	return w.Bytes()
}

// UnmarshalAckMO decode Ack MO from bytes
func UnmarshalAckMO(b []byte) (a RpAck, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *RpAck) UnmarshalRPMO(b []byte) error {
	ud, e := d.unmarshal(b, 2)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalTPMO(ud)
	}
	return e
}

// UnmarshalAckMT decode Ack MT from bytes
func UnmarshalAckMT(b []byte) (a RpAck, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *RpAck) UnmarshalRPMT(b []byte) error {
	ud, e := d.unmarshal(b, 3)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalTPMT(ud)
	}
	return e
}

func (d *RpAck) unmarshal(b []byte, mti byte) ([]byte, error) {
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
	} else if tmp != 0x41 {
		return nil, UnexpectedInformationElementError{Expected: 0x41, Actual: tmp}
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

func (d RpAck) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: RP-Ack\n")
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.MR)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD: %s\n", Indent, d.UD)
	}

	return w.String()
}
