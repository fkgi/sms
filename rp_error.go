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

// RpError is RP-ERROR RPDU
type RpError struct {
	cpData

	RMR  byte  `json:"rp-mr"`          // M / Message Reference
	CS   byte  `json:"rp-cs"`          // M / Cause
	DIAG *byte `json:"diag,omitempty"` // O / Diagnostics
}

// RpErrorMO is MO RP-ERROR RPDU
type RpErrorMO RpError

// RpErrorMT is MT RP-ERROR RPDU
type RpErrorMT RpError

// MarshalRP output byte data of this RPDU
func (d RpErrorMO) MarshalRP() []byte {
	return RpError(d).marshalRP(true, nil)
}

// MarshalRP output byte data of this RPDU
func (d RpErrorMT) MarshalRP() []byte {
	return RpError(d).marshalRP(false, nil)
}

func (d RpError) marshalRP(mo bool, tp []byte) []byte {
	// func (d rpAnswer) marshalErr(mo bool, tp []byte) []byte {
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
func (d RpErrorMO) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// MarshalCP output byte data of this CPDU
func (d RpErrorMT) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalRpErrorMO decode Error MO from bytes
func UnmarshalRpErrorMO(b []byte) (a RpErrorMO, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *RpErrorMO) UnmarshalRP(b []byte) (e error) {
	if b, e = (*RpError)(d).unmarshalRP(true, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

// UnmarshalRpErrorMT decode Error MO from bytes
func UnmarshalRpErrorMT(b []byte) (a RpErrorMT, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP get data of this RPDU
func (d *RpErrorMT) UnmarshalRP(b []byte) (e error) {
	if b, e = (*RpError)(d).unmarshalRP(false, b); e != nil && b != nil {
		e = ErrExtraData
	}
	return
}

func (d *RpError) unmarshalRP(mo bool, b []byte) (tp []byte, e error) {
	// func (d *rpAnswer) unmarshalErr(mo bool, b []byte) (tp []byte, e error) {
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
func (d *RpErrorMO) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *RpErrorMT) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

func (d RpErrorMO) String() string {
	return RpError(d).String()
}

func (d RpErrorMT) String() string {
	return RpError(d).String()
}

func (d RpError) String() string {
	// func (d rpAnswer) stringErr() string {
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

func (d RpError) Error() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "RP-Error, cause=%s", rpCauseStat(d.CS))
	if d.DIAG != nil {
		fmt.Fprintf(w, ", diagnostic=%d", *d.DIAG)
	}
	return w.String()
}
