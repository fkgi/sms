package sms

import (
	"bytes"
	"fmt"
	"io"
)

// RpAck is RP-ACK RPDU
type RpAck struct {
	cpData

	RMR byte `json:"rp-mr"` // M / Message Reference
}

// RpAckMO is MO RP-ACK RPDU
type RpAckMO RpAck

// RpAckMT is MT RP-ACK RPDU
type RpAckMT RpAck

// MarshalRP output byte data of this RPDU
func (d RpAckMO) MarshalRP() []byte {
	return RpAck(d).marshalRP(true, nil)
}

// MarshalRP output byte data of this RPDU
func (d RpAckMT) MarshalRP() []byte {
	return RpAck(d).marshalRP(false, nil)
}

func (d RpAck) marshalRP(mo bool, tp []byte) []byte {
	// func (d rpAnswer) marshalAck(mo bool, tp []byte) []byte {
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
func (d RpAckMO) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// MarshalCP output byte data of this CPDU
func (d RpAckMT) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalRpAckMO decode Ack MO from bytes
func UnmarshalRpAckMO(b []byte) (a RpAckMO, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *RpAckMO) UnmarshalRP(b []byte) (e error) {
	if b, e = (*RpAck)(d).unmarshalRP(true, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

// UnmarshalRpAckMT decode Ack MT from bytes
func UnmarshalRpAckMT(b []byte) (a RpAckMT, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *RpAckMT) UnmarshalRP(b []byte) (e error) {
	if b, e = (*RpAck)(d).unmarshalRP(false, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

func (d *RpAck) unmarshalRP(mo bool, b []byte) (tp []byte, e error) {
	// func (d *rpAnswer) unmarshalAck(mo bool, b []byte) (tp []byte, e error) {
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
func (d *RpAckMO) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *RpAckMT) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

func (d RpAckMO) String() string {
	return RpAck(d).String()
}

func (d RpAckMT) String() string {
	return RpAck(d).String()
}

func (d RpAck) String() string {
	// func (d rpAnswer) stringAck() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "RP-Ack\n")
	fmt.Fprintf(w, "%sCP-TI: %s\n", Indent, cpTIStat(d.TI))
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.RMR)

	return w.String()
}
