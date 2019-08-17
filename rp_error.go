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
	MR   byte  `json:"mr"`             // M / Message Reference
	CS   byte  `json:"cs"`             // M / Cause
	Diag *byte `json:"diag,omitempty"` // O / Diagnostics
	UD   TPDU  `json:"ud,omitempty"`   // O / User Data
}

// MarshalRPMO returns binary data
func (d RpError) MarshalRPMO() []byte {
	return d.marshal(4)
}

// MarshalRPMT returns binary data
func (d RpError) MarshalRPMT() []byte {
	return d.marshal(5)
}

func (d RpError) marshal(mti byte) []byte {
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
		w.WriteByte(0x41)
		w.WriteByte(byte(len(b)))
		w.Write(b)
	}

	return w.Bytes()
}

// UnmarshalErrorMO decode Error MO from bytes
func UnmarshalErrorMO(b []byte) (a RpError, e error) {
	e = a.UnmarshalRPMO(b)
	return
}

// UnmarshalRPMO reads binary data
func (d *RpError) UnmarshalRPMO(b []byte) error {
	ud, e := d.unmarshal(b, 4)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalTPMO(ud)
	}
	return e
}

// UnmarshalErrorMT decode Error MO from bytes
func UnmarshalErrorMT(b []byte) (a RpError, e error) {
	e = a.UnmarshalRPMT(b)
	return
}

// UnmarshalRPMT reads binary data
func (d *RpError) UnmarshalRPMT(b []byte) error {
	ud, e := d.unmarshal(b, 5)
	if e == nil && ud != nil {
		d.UD, e = UnmarshalTPMT(ud)
	}
	return e
}

func (d *RpError) unmarshal(b []byte, mti byte) ([]byte, error) {
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

func (d RpError) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: RP-Error\n")
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-CS: cause=%s, diagnostic=%d\n",
		Indent, rpCauseStat(d.CS), *d.Diag)
	if d.UD != nil {
		fmt.Fprintf(w, "%sRP-UD: %s\n", Indent, d.UD)
	}

	return w.String()
}
