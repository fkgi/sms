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

// UnmarshalCPMO parse byte data to CPDU.
func UnmarshalCPMO(b []byte) (CPDU, error) {
	if len(b) < 2 {
		return nil, io.EOF
	}
	switch b[1] {
	case 0x01:
		ti, rp, e := unmarshalCpDataWith(b)
		if e != nil {
			return nil, e
		}
		rpdu, e := UnmarshalRPMO(rp)
		if e != nil {
			return nil, e
		}
		return rpdu, nil
	case 0x04:
		return UnmarshalCpAck(b)
	case 0x10:
		return UnmarshalCpError(b)
	}
	return nil, UnexpectedMessageTypeError{Actual: b[1]}
}

// UnmarshalCPMT parse byte data to CPDU.
func UnmarshalCPMT(b []byte) (CPDU, error) {
	if len(b) < 2 {
		return nil, io.EOF
	}
	switch b[1] {
	case 0x01:
		return UnmarshalCpDataMT(b)
	case 0x04:
		return UnmarshalCpAck(b)
	case 0x10:
		return UnmarshalCpError(b)
	}
	return nil, UnexpectedMessageTypeError{Actual: b[1]}
}

func marshalCpDataWith(ti byte, rp []byte) []byte {
	w := new(bytes.Buffer)

	b := (ti & 0x0f) << 4
	b |= 0x09
	w.WriteByte(b)
	w.WriteByte(0x01)
	w.WriteByte(byte(len(rp)))
	w.Write(rp)

	return w.Bytes()
}

func unmarshalCpDataWith(b []byte) (ti byte, rp []byte, e error) {
	r := bytes.NewReader(b)

	var tmp byte
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp&0x0f != 0x09 {
		e = UnexpectedMessageTypeError{
			Expected: 0x09, Actual: tmp & 0x0f}
		return
	}
	ti = tmp >> 4
	ti &= 0x0f
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 0x01 {
		e = UnexpectedMessageTypeError{
			Expected: 0x01, Actual: tmp}
		return
	}

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
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
		return
	}
	return
}
