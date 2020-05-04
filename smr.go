package sms

import "time"

var (
	// TR1M timer waiting time
	TR1M = time.Duration(35 * time.Second)
)

// SMR is SM-RP handler
type SMR struct {
	stack     [256]chan RPDU
	SCAddress Address

	RelayReq    func(RPDU) error
	TranspInd   func(TPDU) (TPDU, error)
	MemAvailInd func() error
}

// TranspReq handle TP-DATA
func (smr SMR) TranspReq(r TPDU) (TPDU, error) {
	if smr.RelayReq == nil {
		return nil, RpError{CS: 97}
	}

	mr := -1
	ch := make(chan RPDU)
	for i, c := range smr.stack {
		if c == nil {
			mr = i
			smr.stack[i] = ch
			break
		}
	}
	if mr == -1 {
		return nil, RpError{CS: 42}
	}

	var e error
	var isNW bool
	switch v := r.(type) {
	case Submit:
		isNW = false
		v.RMR = byte(mr)
		v.SCA = smr.SCAddress
		e = smr.RelayReq(v)
	case Command:
		isNW = false
		v.RMR = byte(mr)
		v.SCA = smr.SCAddress
		e = smr.RelayReq(v)
	case Deliver:
		isNW = true
		v.RMR = byte(mr)
		v.SCA = smr.SCAddress
		e = smr.RelayReq(v)
	case StatusReport:
		isNW = true
		v.RMR = byte(mr)
		v.SCA = smr.SCAddress
		e = smr.RelayReq(v)
	default:
		e = RpError{CS: 97}
	}

	if e != nil {
		smr.stack[mr] = nil
		close(ch)

		switch e.(type) {
		case RpError:
			return nil, e
		case CpError:
			return nil, RpError{CS: 47}
		default:
			return nil, RpError{CS: 111}
		}
	}

	t := time.AfterFunc(TR1M, func() {
		if isNW {
			ch <- RpErrorMO{RMR: byte(mr), CS: 47}
		}
		ch <- RpErrorMT{RMR: byte(mr), CS: 47}
	})
	a := <-ch
	t.Stop()
	smr.stack[mr] = nil
	close(ch)

	switch v := a.(type) {
	case DeliverReport:
		if v.RMR != byte(mr) {
			return nil, RpError{CS: 81}
		}
		if isNW {
			return v, nil
		}
	case SubmitReport:
		if v.RMR != byte(mr) {
			return nil, RpError{CS: 81}
		}
		if !isNW {
			return v, nil
		}
	case RpAckMT:
		if v.RMR != byte(mr) {
			return nil, RpError{CS: 81}
		}
		if isNW {
			return nil, nil
		}
	case RpAckMO:
		if v.RMR != byte(mr) {
			return nil, RpError{CS: 81}
		}
		if !isNW {
			return nil, nil
		}
	case RpErrorMT:
		if v.RMR != byte(mr) {
			return nil, RpError{CS: 81}
		}
		if isNW {
			return nil, RpError(v)
		}
	case RpErrorMO:
		if v.RMR != byte(mr) {
			return nil, RpError{CS: 81}
		}
		if !isNW {
			return nil, RpError(v)
		}
	}

	return nil, RpError{CS: 95}
}

// RelayInd handle RP-DATA/SMMA
func (smr SMR) RelayInd(r RPDU) (RPDU, error) {
	var a TPDU
	var e error
	var mr byte
	var isNW bool

	switch v := r.(type) {
	case MemoryAvailable:
		if smr.MemAvailInd == nil {
			return nil, CpError{CS: 97}
		}
		mr = v.RMR
		isNW = true
		e = smr.MemAvailInd()
	case Submit:
		if smr.TranspInd == nil {
			return nil, CpError{CS: 97}
		}
		mr = v.RMR
		isNW = true
		a, e = smr.TranspInd(v)
	case Command:
		if smr.TranspInd == nil {
			return nil, CpError{CS: 97}
		}
		mr = v.RMR
		isNW = true
		a, e = smr.TranspInd(v)
	case Deliver:
		if smr.TranspInd == nil {
			return nil, CpError{CS: 97}
		}
		mr = v.RMR
		isNW = false
		a, e = smr.TranspInd(v)
	case StatusReport:
		if smr.TranspInd == nil {
			return nil, CpError{CS: 97}
		}
		mr = v.RMR
		isNW = false
		a, e = smr.TranspInd(v)
	default:
		return nil, CpError{CS: 97}
	}

	switch v := e.(type) {
	case RpError:
		v.RMR = mr
		if isNW {
			return RpErrorMT(v), nil
		}
		return RpErrorMO(v), nil
	case nil:
	default:
		return nil, CpError{CS: 111}
	}

	switch v := a.(type) {
	case SubmitReport:
		if isNW {
			v.RMR = mr
			return v, nil
		}
	case DeliverReport:
		if !isNW {
			v.RMR = mr
			return v, nil
		}
	case nil:
		a := RpAck{RMR: mr}
		if isNW {
			return RpAckMT(a), nil
		}
		return RpAckMO(a), nil
	}
	return nil, CpError{CS: 111}
}
