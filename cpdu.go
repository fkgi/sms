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
		var c cpData
		rp, e := c.unmarshal(b)
		if e != nil {
			return nil, e
		}
		rpdu, e := UnmarshalRPMO(rp)
		if e != nil {
			return nil, e
		}
		return rpdu, nil
	case 0x04:
		return UnmarshalAck(b)
	case 0x10:
		return UnmarshalError(b)
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
		// return UnmarshalDataMT(b)
	case 0x04:
		return UnmarshalAck(b)
	case 0x10:
		return UnmarshalError(b)
	}
	return nil, UnexpectedMessageTypeError{Actual: b[1]}
}

func unmarshalCpHeader(mti byte, b []byte) (byte, error) {
	if len(b) < 2 {
		return 0, InvalidLengthError{}
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
	d.TI, e = unmarshalCpHeader(0x01, b)
	if e != nil {
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
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
		return
	}
	return
}