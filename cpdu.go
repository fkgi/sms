package sms

import (
	"bytes"
	"fmt"
	"io"
)

// CPDU represents a SMS CP PDU
type CPDU interface {
	MarshalCP() []byte
	fmt.Stringer
}

// UnmarshalerCP is the interface implemented by types
// that can unmarshal a CPDU
type UnmarshalerCP interface {
	UnmarshalCP([]byte) error
}

// UnmarshalCPMO parse byte data to CPDU.
func UnmarshalCPMO(b []byte) (CPDU, error) {
	if len(b) < 2 {
		return nil, io.EOF
	}
	switch b[1] {
	case 0x01:
		var c cpData
		var e error
		if b, e = c.unmarshal(b); e != nil {
			return nil, e
		}
		return unmarshalRPMO(b, c)
	case 0x04:
		return UnmarshalAck(b)
	case 0x10:
		return UnmarshalError(b)
	}
	return nil, UnknownMessageTypeError{Actual: b[1]}
}

// UnmarshalCPMT parse byte data to CPDU.
func UnmarshalCPMT(b []byte) (CPDU, error) {
	if len(b) < 2 {
		return nil, io.EOF
	}
	switch b[1] {
	case 0x01:
		var c cpData
		var e error
		if b, e = c.unmarshal(b); e != nil {
			return nil, e
		}
		return unmarshalRPMT(b, c)
	case 0x04:
		return UnmarshalAck(b)
	case 0x10:
		return UnmarshalError(b)
	}
	return nil, UnknownMessageTypeError{Actual: b[1]}
}

type cpData struct {
	TI byte `json:"ti"` // M / Transaction identifier
}

func (d cpData) marshal(rp []byte) []byte {
	w := new(bytes.Buffer)

	b := (d.TI & 0x0f) << 4
	b |= 0x09
	w.WriteByte(b)
	w.WriteByte(0x01)
	w.WriteByte(byte(len(rp)))
	w.Write(rp)

	return w.Bytes()
}

func (d *cpData) unmarshal(b []byte) (rp []byte, e error) {
	if d.TI, e = unmarshalCpHeader(0x01, b); e != nil {
		return
	}
	r := bytes.NewReader(b[2:])

	var tmp byte
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	rp = make([]byte, int(tmp))

	var l int
	if l, e = r.Read(rp); e != nil {
		return
	}
	if l != len(rp) {
		e = io.EOF
	} else if r.Len() != 0 {
		e = ErrExtraData
	}
	return
}

func unmarshalCpHeader(mti byte, b []byte) (byte, error) {
	if len(b) < 2 {
		return 0, io.EOF
	}

	if b[0]&0x0f != 0x09 {
		return 0, UnexpectedMessageTypeError{
			Expected: 0x09, Actual: b[0] & 0x0f}
	}
	ti := b[0] >> 4
	ti &= 0x0f

	if b[1] != mti {
		return 0, UnexpectedMessageTypeError{
			Expected: mti, Actual: b[1]}
	}
	return ti, nil
}
