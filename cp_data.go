package sms

import (
	"bytes"
	"fmt"
	"io"
)

// CpData is CP-DATA CPDU
type CpData struct {
	TI byte `json:"ti"`           // M / Transaction identifier
	UD RPDU `json:"ud,omitempty"` // M / User Data
}

// MarshalCPMO returns binary data
func (d CpData) MarshalCPMO() []byte {
	w := new(bytes.Buffer)

	b := (d.TI & 0x0f) << 4
	b |= 0x09
	w.WriteByte(b)

	w.WriteByte(0x01)

	if d.UD != nil {
		data := d.UD.MarshalRPMO()
		w.WriteByte(byte(len(data)))
		w.Write(data)
	} else {
		w.WriteByte(0)
	}

	return w.Bytes()
}

// MarshalCPMT returns binary data
func (d CpData) MarshalCPMT() []byte {
	w := new(bytes.Buffer)

	b := (d.TI & 0x0f) << 4
	b |= 0x09
	w.WriteByte(b)

	w.WriteByte(0x01)

	if d.UD != nil {
		data := d.UD.MarshalRPMT()
		w.WriteByte(byte(len(data)))
		w.Write(data)
	} else {
		w.WriteByte(0)
	}

	return w.Bytes()
}

// UnmarshalCpDataMO decode Ack MT from bytes
func UnmarshalCpDataMO(b []byte) (a CpData, e error) {
	e = a.UnmarshalCPMO(b)
	return
}

// UnmarshalCPMO reads binary data
func (d *CpData) UnmarshalCPMO(b []byte) error {
	ud, e := d.unmarshal(b)
	if e != nil {
		return e
	}
	if len(ud) != 0 {
		d.UD, e = UnmarshalRPMO(ud)
	}
	return e
}

// UnmarshalCpDataMT decode Ack MT from bytes
func UnmarshalCpDataMT(b []byte) (a CpData, e error) {
	e = a.UnmarshalCPMT(b)
	return
}

// UnmarshalCPMT reads binary data
func (d *CpData) UnmarshalCPMT(b []byte) error {
	ud, e := d.unmarshal(b)
	if e != nil {
		return e
	}
	if len(ud) != 0 {
		d.UD, e = UnmarshalRPMT(ud)
	}
	return e
}

func (d *CpData) unmarshal(b []byte) ([]byte, error) {
	r := bytes.NewReader(b)

	if tmp, e := r.ReadByte(); e != nil {
		return nil, e
	} else if tmp&0x0f != 0x09 {
		return nil, UnexpectedMessageTypeError{
			Expected: 0x09, Actual: tmp & 0x0f}
	} else {
		d.TI = tmp >> 4
		d.TI &= 0x0f

	}
	if tmp, e := r.ReadByte(); e != nil {
		return nil, e
	} else if tmp != 0x01 {
		return nil, UnexpectedMessageTypeError{
			Expected: 0x01, Actual: tmp}
	}

	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return nil, e
	}
	if len(b) != 0 {
		if n, e := r.Read(b); e != nil {
			return nil, e
		} else if n != len(b) {
			return nil, io.EOF
		}
	}
	if r.Len() != 0 {
		return nil, InvalidLengthError{}
	}
	return b, nil
}

func (d CpData) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: CP-Data\n")
	fmt.Fprintf(w, "%sCP-TI: %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sCP-UD: %s\n", Indent, d.UD)

	return w.String()
}
