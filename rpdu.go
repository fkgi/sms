package sms

import (
	"bytes"
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
		var rp rpData
		var e error
		if b, e = rp.unmarshal(true, b); e != nil {
			return nil, e
		}
		rp.cpData = c

		switch b[0] & 0x03 {
		case 0x01:
			var tp Submit
			e = tp.UnmarshalTP(b)
			tp.rpData = rp
			return tp, e
		case 0x02:
			var tp Command
			e = tp.UnmarshalTP(b)
			tp.rpData = rp
			return tp, e
		}
	case 0x02:
		var rp RpAck
		var e error
		if b, e = rp.unmarshalRP(true, b); e != nil {
			return nil, e
		}
		rp.cpData = c
		if b == nil {
			return RpAckMO(rp), nil
		}

		var tp DeliverReport
		e = tp.UnmarshalTP(b)
		tp.cpData = rp.cpData
		tp.RMR = rp.RMR
		return tp, e
	case 0x04:
		var rp RpError
		var e error
		if b, e = rp.unmarshalRP(true, b); e != nil {
			return nil, e
		}
		rp.cpData = c
		if b == nil {
			return RpErrorMO(rp), nil
		}

		var tp DeliverReport
		e = tp.UnmarshalTP(b)
		tp.cpData = rp.cpData
		tp.RMR = rp.RMR
		tp.CS = rp.CS
		tp.DIAG = rp.DIAG
		return tp, e
	case 0x06:
		rp, e := UnmarshalMemoryAvailable(b)
		rp.cpData = c
		return rp, e
	}
	return nil, UnknownMessageTypeError{Actual: b[0]}
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
		var rp rpData
		var e error
		if b, e = rp.unmarshal(false, b); e != nil {
			return nil, e
		}
		rp.cpData = c

		switch b[0] & 0x03 {
		case 0x00:
			var tp Deliver
			e = tp.UnmarshalTP(b)
			tp.rpData = rp
			return tp, e
		case 0x02:
			var tp StatusReport
			e = tp.UnmarshalTP(b)
			tp.rpData = rp
			return tp, e
		}
	case 0x03:
		var rp RpAck
		var e error
		if b, e = rp.unmarshalRP(false, b); e != nil {
			return nil, e
		}
		rp.cpData = c
		if b == nil {
			return RpAckMT(rp), nil
		}

		var tp SubmitReport
		e = tp.UnmarshalTP(b)
		tp.cpData = rp.cpData
		tp.RMR = rp.RMR
		return tp, e
	case 0x05:
		var rp RpError
		var e error
		if b, e = rp.unmarshalRP(false, b); e != nil {
			return nil, e
		}
		rp.cpData = c
		if b == nil {
			return RpErrorMT(rp), nil
		}

		var tp SubmitReport
		e = tp.UnmarshalTP(b)
		tp.cpData = rp.cpData
		tp.RMR = rp.RMR
		tp.CS = rp.CS
		tp.DIAG = rp.DIAG
		return tp, e
	}
	return nil, UnknownMessageTypeError{Actual: b[0]}
}

func unmarshalRpHeader(mti byte, b []byte) (byte, error) {
	if len(b) < 2 {
		return 0, io.EOF
	}

	if b[0] != mti {
		return 0, UnexpectedMessageTypeError{
			Expected: mti, Actual: b[0]}
	}
	return b[1], nil
}

type rpData struct {
	cpData

	RMR byte    `json:"rp-mr"`  // M / Message Reference for RP
	SCA Address `json:"rp-sca"` // M / SC Address
}

func (d rpData) marshal(mo bool, tp []byte) []byte {
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

func (d *rpData) unmarshal(mo bool, b []byte) (tp []byte, e error) {
	if mo {
		d.RMR, e = unmarshalRpHeader(0, b)
	} else {
		d.RMR, e = unmarshalRpHeader(1, b)
	}
	if e != nil {
		return
	}

	r := bytes.NewReader(b[2:])
	var tmp byte
	if mo {
		if tmp, e = r.ReadByte(); e != nil {
			return
		} else if tmp != 0 {
			e = ErrInvalidLength
			return
		}
	}
	if d.SCA, e = readRPAddr(r); e != nil {
		return
	}
	if !mo {
		if tmp, e = r.ReadByte(); e != nil {
			return
		} else if tmp != 0 {
			e = ErrInvalidLength
			return
		}
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

/*
type rpAResult interface {
	marshalRP(bool, []byte) []byte
}

/*
type rpAnswer struct {
	cpData

	RMR  byte  `json:"rmr"`            // M / Message Reference
	CS   byte  `json:"cs"`             // M / Cause
	DIAG *byte `json:"diag,omitempty"` // O / Diagnostics
}
*/
