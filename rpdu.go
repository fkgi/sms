package sms

import (
	"bytes"
	"fmt"
	"io"
)

// RPDU represents a SMS RP PDU
type RPDU interface {
	CPDU
	MarshalRP() []byte
}

// UnmarshalerRP is the interface implemented by types
// that can unmarshal a RPDU
type UnmarshalerRP interface {
	UnmarshalRP([]byte) error
}

// UnmarshalRPMO parse byte data to TPDU as SC.
func UnmarshalRPMO(b []byte) (RPDU, error) {
	return unmarshalRPMO(b, cpData{})
}

func unmarshalRPMO(b []byte, c cpData) (RPDU, error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	switch b[0] & 0x07 {
	case 0x00:
		var rp rpRequest
		var e error
		b, e = rp.unmarshal(true, b)
		rp.cpData = c
		if e != nil {
			return nil, e
		}
		switch b[0] & 0x03 {
		case 0x01:
			var tp Submit
			e = tp.UnmarshalTP(b)
			tp.rpRequest = rp
			return tp, e
		case 0x02:
			// var tp Command
			// e = tp.UnmarshalTP(b)
			// tp.rpRequest = rp
			// return tp, e
		}
	case 0x02:
		var rp rpAnswer
		var e error
		b, e = rp.unmarshalAck(true, b)
		rp.cpData = c
		if e != nil {
			return nil, e
		}
		if b == nil {
			return AckMO{rp}, nil
		}
		var tp DeliverReport
		e = tp.UnmarshalTP(b)
		tp.rpAnswer = rp
		return tp, e
	case 0x04:
		var rp rpAnswer
		var e error
		b, e = rp.unmarshalErr(true, b)
		rp.cpData = c
		if e != nil {
			return nil, e
		}
		if b == nil {
			return ErrorMO{rp}, nil
		}
		var tp DeliverReport
		e = tp.UnmarshalTP(b)
		tp.rpAnswer = rp
		return tp, e
	case 0x06:
		rp, e := UnmarshalMemoryAvailable(b)
		rp.cpData = c
		return rp, e
	}
	return nil, UnexpectedMessageTypeError{Actual: b[0]}
}

// UnmarshalRPMT parse byte data to TPDU as MS.
func UnmarshalRPMT(b []byte) (RPDU, error) {
	return unmarshalRPMT(b, cpData{})
}

func unmarshalRPMT(b []byte, c cpData) (RPDU, error) {
	if len(b) == 0 {
		return nil, io.EOF
	}
	switch b[0] & 0x07 {
	case 0x01:
		var rp rpRequest
		var e error
		b, e = rp.unmarshal(false, b)
		rp.cpData = c
		if e != nil {
			return nil, e
		}
		switch b[0] & 0x03 {
		case 0x00:
			var tp Deliver
			e = tp.UnmarshalTP(b)
			tp.rpRequest = rp
			return tp, e
		case 0x02:
			var tp StatusReport
			e = tp.UnmarshalTP(b)
			tp.rpRequest = rp
			return tp, e
		}
	case 0x03:
		var rp rpAnswer
		var e error
		b, e = rp.unmarshalAck(false, b)
		rp.cpData = c
		if e != nil {
			return nil, e
		}
		if b == nil {
			return AckMT{rp}, nil
		}
		var tp SubmitReport
		e = tp.UnmarshalTP(b)
		tp.rpAnswer = rp
		return tp, e
	case 0x05:
		var rp rpAnswer
		var e error
		b, e = rp.unmarshalErr(false, b)
		rp.cpData = c
		if e != nil {
			return nil, e
		}
		if b == nil {
			return ErrorMT{rp}, nil
		}
		var tp DeliverReport
		e = tp.UnmarshalTP(b)
		tp.rpAnswer = rp
		return tp, e
	}
	return nil, UnexpectedMessageTypeError{Actual: b[0]}
}

type rpRequest struct {
	cpData

	RMR byte    `json:"rmr"` // M / Message Reference for RP
	SCA Address `json:"sca"` // M / Destination SC Address
}

func (d rpRequest) marshal(mo bool, tp []byte) []byte {
	w := new(bytes.Buffer)

	if mo {
		w.WriteByte(0) // MTI
		w.WriteByte(d.RMR)
		w.WriteByte(0) // OA
		_, a := d.SCA.Marshal()
		w.WriteByte(byte(len(a)))
		w.Write(a)
	} else {
		w.WriteByte(1) // MTI
		w.WriteByte(d.RMR)
		_, a := d.SCA.Marshal()
		w.WriteByte(byte(len(a)))
		w.Write(a)
		w.WriteByte(0) // DA
	}
	w.WriteByte(byte(len(tp)))
	w.Write(tp)

	return w.Bytes()
}

func (d *rpRequest) unmarshal(mo bool, b []byte) (tp []byte, e error) {
	r := bytes.NewReader(b)
	var tmp byte

	if mo {
		if tmp, e = r.ReadByte(); e != nil {
			return
		} else if tmp != 0 {
			e = UnexpectedMessageTypeError{
				Expected: 0, Actual: tmp}
			return
		}
		if d.RMR, e = r.ReadByte(); e != nil {
			return
		}
		if tmp, e = r.ReadByte(); e != nil {
			return
		} else if tmp != 0 {
			e = InvalidLengthError{}
			return
		}
		if d.SCA, e = readRPAddr(r); e != nil {
			return
		}
	} else {
		if tmp, e = r.ReadByte(); e != nil {
			return
		} else if tmp != 1 {
			e = UnexpectedMessageTypeError{
				Expected: 1, Actual: tmp}
			return
		}
		if d.RMR, e = r.ReadByte(); e != nil {
			return
		}
		if d.SCA, e = readRPAddr(r); e != nil {
			return
		}
		if tmp, e = r.ReadByte(); e != nil {
			return
		} else if tmp != 0 {
			e = InvalidLengthError{}
			return
		}
	}

	if tmp, e = r.ReadByte(); e == nil {
		tp = make([]byte, int(tmp))
	} else {
		return
	}
	var n int
	if n, e = r.Read(tp); e != nil {
		return
	} else if n != len(tp) {
		e = io.EOF
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
	}
	return
}

type rpAnswer struct {
	cpData

	RMR  byte  `json:"rmr"`            // M / Message Reference
	CS   byte  `json:"cs"`             // M / Cause
	DIAG *byte `json:"diag,omitempty"` // O / Diagnostics
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

func (d *rpAnswer) unmarshalAck(mo bool, b []byte) (tp []byte, e error) {
	r := bytes.NewReader(b)

	var tmp byte
	var mti byte = 3
	if mo {
		mti = 2
	}
	if tmp, e = r.ReadByte(); e != nil {
		return
	} else if tmp != mti {
		e = UnexpectedMessageTypeError{
			Expected: mti, Actual: b[0]}
		return
	}
	if d.RMR, e = r.ReadByte(); e != nil {
		return
	}

	if r.Len() == 0 {
		return
	}

	if tmp, e = r.ReadByte(); e != nil {
		return
	} else if tmp != 0x41 {
		e = UnexpectedInformationElementError{
			Expected: 0x41, Actual: tmp}
		return
	}
	if tmp, e = r.ReadByte(); e == nil {
		tp = make([]byte, int(tmp))
	} else {
		return
	}
	var n int
	if n, e = r.Read(tp); e != nil {
		return
	} else if n != len(tp) {
		e = io.EOF
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
	}
	return
}

func (d *rpAnswer) unmarshalErr(mo bool, b []byte) (tp []byte, e error) {
	r := bytes.NewReader(b)

	var tmp byte
	var mti byte = 5
	if mo {
		mti = 4
	}
	if tmp, e = r.ReadByte(); e != nil {
		return
	} else if tmp != mti {
		e = UnexpectedMessageTypeError{
			Expected: mti, Actual: b[0]}
		return
	}
	if d.RMR, e = r.ReadByte(); e != nil {
		return
	}
	if tmp, e = r.ReadByte(); e != nil {
		return
	} else if tmp == 0 || tmp > 2 {
		e = InvalidLengthError{}
		return
	} else if d.CS, e = r.ReadByte(); e != nil {
		return
	} else if tmp == 2 {
		if tmp, e = r.ReadByte(); e != nil {
			return
		}
		d.DIAG = &tmp
	}

	if r.Len() == 0 {
		return
	}

	if tmp, e = r.ReadByte(); e != nil {
		return
	} else if tmp != 0x41 {
		e = UnexpectedInformationElementError{
			Expected: 0x41, Actual: tmp}
		return
	}
	if tmp, e = r.ReadByte(); e == nil {
		tp = make([]byte, int(tmp))
	} else {
		return
	}
	var n int
	if n, e = r.Read(tp); e != nil {
		return
	} else if n != len(tp) {
		e = io.EOF
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
	}
	return
}

func (d rpAnswer) stringAck() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "RP-Ack\n")
	fmt.Fprintf(w, "%sCP-TI: %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.RMR)

	return w.String()
}

func (d rpAnswer) stringErr() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "RP-Error\n")
	fmt.Fprintf(w, "%sCP-TI: %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sRP-MR: %d\n", Indent, d.RMR)
	fmt.Fprintf(w, "%sRP-CS: cause=%s", Indent, rpCauseStat(d.CS))
	if d.DIAG != nil {
		fmt.Fprintf(w, ", diagnostic=%d", *d.DIAG)
	}
	fmt.Fprintf(w, "\n")

	return w.String()
}
