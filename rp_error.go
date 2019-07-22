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
	e = a.UnmarshalMORP(b)
	return
}

// UnmarshalMORP reads binary data
func (d *Error) UnmarshalMORP(b []byte) error {
	ud, e := d.unmarshal(b, 4)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalMOTP(ud)
	}
	return e
}

// UnmarshalErrorMT decode Error MO from bytes
func UnmarshalErrorMT(b []byte) (a Error, e error) {
	e = a.UnmarshalMTRP(b)
	return
}

// UnmarshalMTRP reads binary data
func (d *Error) UnmarshalMTRP(b []byte) error {
	ud, e := d.unmarshal(b, 5)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalMTTP(ud)
	}
	return e
}

func (d *Error) unmarshal(b []byte, mti byte) ([]byte, error) {
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
	if l, e := r.ReadByte(); e != nil {
		return nil, e
	} else if l == 0 || l > 2 {
		return nil, InvalidLengthError{}
	} else if d.CS, e = r.ReadByte(); e != nil {
		return nil, e
	} else if l == 2 {
		if l, e = r.ReadByte(); e != nil {
			return nil, e
		}
		d.Diag = &l
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
