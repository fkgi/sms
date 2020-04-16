package sms

import (
	"bytes"
	"fmt"
	"io"
)

func rpCauseStat(c byte) string {
	switch c {
	case 1:
		return "Unassigned (unallocated) number"
	case 8:
		return "Operator determined barring"
	case 10:
		return "Call barred"
	case 11:
		return "Reserved"
	case 21:
		return "Short message transfer rejected"
	case 22:
		return "Memory capacity exceeded"
	case 27:
		return "Destination out of order"
	case 28:
		return "Unidentified subscriber"
	case 29:
		return "Facility rejected"
	case 30:
		return "Unknown subscriber"
	case 38:
		return "Network out of order"
	case 41:
		return "Temporary failure"
	case 42:
		return "Congestion"
	case 47:
		return "Resources unavailable, unspecified"
	case 50:
		return "Requested facility not subscribed"
	case 69:
		return "Requested facility not implemented"
	case 81:
		return "Invalid short message transfer reference value"
	case 95:
		return "Semantically incorrect message"
	case 96:
		return "Invalid mandatory information"
	case 97:
		return "Message type non existent or not implemented"
	case 98:
		return "Message not compatible with short message protocol state"
	case 99:
		return "Information element non existent or not implemented"
	case 111:
		return "Protocol error, unspecified"
	case 127:
		return "Interworking, unspecified"
	}
	return fmt.Sprintf("Reserved(%d)", c)
}

// ErrorMO is MO RP-ERROR RPDU
type ErrorMO struct {
	rpAnswer
}

// ErrorMT is MT RP-ERROR RPDU
type ErrorMT struct {
	rpAnswer
}

// MarshalRP output byte data of this RPDU
func (d ErrorMO) MarshalRP() []byte {
	return d.marshalErr(true, nil)
}

// MarshalRP output byte data of this RPDU
func (d ErrorMT) MarshalRP() []byte {
	return d.marshalErr(false, nil)
}

func (d rpAnswer) marshalErr(mo bool, tp []byte) []byte {
	w := new(bytes.Buffer)

	if mo {
		w.WriteByte(4)
	} else {
		w.WriteByte(5)
	}
	w.WriteByte(d.RMR)

	if d.DIAG != nil {
		w.WriteByte(2)
		w.WriteByte(d.CS)
		w.WriteByte(*d.DIAG)
	} else {
		w.WriteByte(1)
		w.WriteByte(d.CS)
	}

	if tp != nil {
		w.WriteByte(0x41)
		w.WriteByte(byte(len(tp)))
		w.Write(tp)
	}

	return w.Bytes()
}

// MarshalCP output byte data of this CPDU
func (d ErrorMO) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// MarshalCP output byte data of this CPDU
func (d ErrorMT) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalErrorMO decode Error MO from bytes
func UnmarshalErrorMO(b []byte) (a ErrorMO, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *ErrorMO) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshalErr(true, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

// UnmarshalErrorMT decode Error MO from bytes
func UnmarshalErrorMT(b []byte) (a ErrorMT, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *ErrorMT) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshalErr(false, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

func (d *rpAnswer) unmarshalErr(mo bool, b []byte) (tp []byte, e error) {
	if mo {
		d.RMR, e = unmarshalRpHeader(4, b)
	} else {
		d.RMR, e = unmarshalRpHeader(5, b)
	}
	if e != nil {
		return
	}

	r := bytes.NewReader(b[2:])
	var tmp byte
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp == 0 {
		e = io.EOF
		return
	}
	if tmp > 2 {
		e = ErrExtraData
		return
	}
	if d.CS, e = r.ReadByte(); e != nil {
		return
	}
	if tmp == 2 {
		var diag byte
		if diag, e = r.ReadByte(); e != nil {
			return
		}
		d.DIAG = &diag
	}

	if r.Len() == 0 {
		return
	}

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
func (d *ErrorMO) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *ErrorMT) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

func (d ErrorMO) String() string {
	return d.stringErr()
}

func (d ErrorMT) String() string {
	return d.stringErr()
}

func (d rpAnswer) stringErr() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "RP-Error\n")
	fmt.Fprintf(w, "%sCP-TI: %s\n", Indent, cpTIStat(d.TI))
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.RMR)
	fmt.Fprintf(w, "%sRP-CS: cause=%s", Indent, rpCauseStat(d.CS))
	if d.DIAG != nil {
		fmt.Fprintf(w, ", diagnostic=%d", *d.DIAG)
	}
	fmt.Fprintf(w, "\n")

	return w.String()
}
