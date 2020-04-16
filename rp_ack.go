package sms

import (
	"bytes"
	"fmt"
	"io"
)

// AckMO is MO RP-ACK RPDU
type AckMO struct {
	rpAnswer
}

// AckMT is MT RP-ACK RPDU
type AckMT struct {
	rpAnswer
}

// MarshalRP output byte data of this RPDU
func (d AckMO) MarshalRP() []byte {
	return d.marshalAck(true, nil)
}

// MarshalRP output byte data of this RPDU
func (d AckMT) MarshalRP() []byte {
	return d.marshalAck(false, nil)
}

func (d rpAnswer) marshalAck(mo bool, tp []byte) []byte {
	w := new(bytes.Buffer)

	if mo {
		w.WriteByte(2)
	} else {
		w.WriteByte(3)
	}
	w.WriteByte(d.RMR)

	if tp != nil {
		w.WriteByte(0x41)
		w.WriteByte(byte(len(tp)))
		w.Write(tp)
	}

	return w.Bytes()
}

// MarshalCP output byte data of this CPDU
func (d AckMO) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// MarshalCP output byte data of this CPDU
func (d AckMT) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalAckMO decode Ack MO from bytes
func UnmarshalAckMO(b []byte) (a AckMO, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *AckMO) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshalAck(true, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

// UnmarshalAckMT decode Ack MT from bytes
func UnmarshalAckMT(b []byte) (a AckMT, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *AckMT) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshalAck(false, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

func (d *rpAnswer) unmarshalAck(mo bool, b []byte) (tp []byte, e error) {
	if mo {
		d.RMR, e = unmarshalRpHeader(2, b)
	} else {
		d.RMR, e = unmarshalRpHeader(3, b)
	}
	if e != nil || len(b) == 2 {
		return
	}

	r := bytes.NewReader(b[2:])
	var tmp byte
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 0x41 {
		e = UnexpectedInformationElementError{
			Expected: 0x41, Actual: tmp}
		return
	}
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	tp = make([]byte, int(tmp))
	var n int
	if n, e = r.Read(tp); e != nil {
		return
	}
	if n != len(tp) {
		e = io.EOF
	} else if r.Len() != 0 {
		e = ErrExtraData
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *AckMO) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *AckMT) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

func (d AckMO) String() string {
	return d.stringAck()
}

func (d AckMT) String() string {
	return d.stringAck()
}

func (d rpAnswer) stringAck() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "RP-Ack\n")
	fmt.Fprintf(w, "%sCP-TI: %s\n", Indent, cpTIStat(d.TI))
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.RMR)

	return w.String()
}
