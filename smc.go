package sms

type state byte

const (
	idle       state = iota
	connecting state = iota
	waitAck    state = iota
	waitData   state = iota
)

// SMC is SM-CP handler
type SMC struct {
	rxStack   [7]chan CPDU
	txStack   [7]chan CPDU
	SCAddress Address

	CtrlReq     func(CPDU)
	TranspInd   func(TPDU) (TPDU, error)
	MemAvailInd func() error
}

// CtrlInd handle CP-DATA/CP-ACK/CP-ERROR
func (smc SMC) CtrlInd(pdu CPDU) {
	var ti byte
	var isTx bool

	switch v := pdu.(type) {
	case Submit:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
		if !isTx && smc.rxStack[ti] == nil {
			c := make(chan CPDU)
			smc.rxStack[ti] = c
			go smc.rxHandler(c, ti)
		}
	case Command:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
		if !isTx && smc.rxStack[ti] == nil {
			c := make(chan CPDU)
			smc.rxStack[ti] = c
			go smc.rxHandler(c, ti)
		}
	case Deliver:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
		if !isTx && smc.rxStack[ti] == nil {
			c := make(chan CPDU)
			smc.rxStack[ti] = c
			go smc.rxHandler(c, ti)
		}
	case StatusReport:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
		if !isTx && smc.rxStack[ti] == nil {
			c := make(chan CPDU)
			smc.rxStack[ti] = c
			go smc.rxHandler(c, ti)
		}
	case SubmitReport:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case DeliverReport:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case MemoryAvailable:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case RpAckMT:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case RpAckMO:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case RpErrorMT:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case RpErrorMO:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case CpAck:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	case CpError:
		ti = v.TI & 0x07
		isTx = v.TI&0x08 == 0x08
	default:
		smc.CtrlReq(CpError{TI: ti, CS: 98})
		return
	}

	var c chan CPDU
	if isTx {
		c = smc.txStack[ti]
	} else {
		c = smc.rxStack[ti]
	}
	if c != nil {
		c <- pdu
	} else if isTx {
		smc.CtrlReq(CpError{TI: ti, CS: 81})
	} else {
		smc.CtrlReq(CpError{TI: ti | 0x08, CS: 81})
	}
}

func (smc SMC) rxHandler(c chan CPDU, ti byte) {
	ti = ti | 0x08
	defer func() {
		smc.rxStack[ti&0x07] = nil
		close(c)
	}()

	pdu := <-c

	var a TPDU
	var e error
	var mr byte
	var isNW bool

	switch v := pdu.(type) {
	case Submit:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = true
		if smc.TranspInd == nil {
			e = RpError{CS: 97}
		} else {
			a, e = smc.TranspInd(v)
		}
	case Command:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = true
		if smc.TranspInd == nil {
			e = RpError{CS: 97}
		} else {
			a, e = smc.TranspInd(v)
		}
	case Deliver:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = false
		if smc.TranspInd == nil {
			e = RpError{CS: 97}
		} else {
			a, e = smc.TranspInd(v)
		}
	case StatusReport:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = false
		if smc.TranspInd == nil {
			e = RpError{CS: 97}
		} else {
			a, e = smc.TranspInd(v)
		}
	case MemoryAvailable:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = true
		if smc.MemAvailInd == nil {
			e = RpError{CS: 97}
		} else {
			e = smc.MemAvailInd()
		}
	case DeliverReport:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = true
		e = RpError{CS: 95}
	case SubmitReport:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = false
		e = RpError{CS: 95}
	case RpAckMO:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = true
		e = RpError{CS: 95}
	case RpAckMT:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = false
		e = RpError{CS: 95}
	case RpErrorMO:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = true
		e = RpError{CS: 95}
	case RpErrorMT:
		smc.CtrlReq(CpAck{TI: ti})
		mr = v.RMR
		isNW = false
		e = RpError{CS: 95}
	case CpAck:
		return
	case CpError:
		return
	}

	if e == nil {
		switch v := a.(type) {
		case SubmitReport:
			if isNW {
				v.RMR = mr
				v.TI = ti
				smc.CtrlReq(v)
			} else {
				e = RpError{CS: 111}
			}
		case DeliverReport:
			if !isNW {
				v.RMR = mr
				v.TI = ti
				smc.CtrlReq(v)
			} else {
				e = RpError{CS: 111}
			}
		case nil:
			a := RpAck{RMR: mr}
			a.TI = ti
			if isNW {
				smc.CtrlReq(RpAckMT(a))
			} else {
				smc.CtrlReq(RpAckMO(a))
			}
		default:
			e = RpError{CS: 111}
		}
	}

	switch v := e.(type) {
	case RpError:
		v.RMR = mr
		v.TI = ti
		if isNW {
			smc.CtrlReq(RpErrorMT(v))
		} else {
			smc.CtrlReq(RpErrorMO(v))
		}
	case nil:
	default:
		a := RpError{RMR: mr, CS: 111}
		a.TI = ti
		if isNW {
			smc.CtrlReq(RpErrorMT(a))
		} else {
			smc.CtrlReq(RpErrorMO(a))
		}
	}

	pdu = <-c
	return
}

// TranspReq handle TP-DATA
func (smc SMC) TranspReq(pdu TPDU) (TPDU, error) {
	var ti byte = 255
	c := make(chan CPDU)
	for i, sc := range smc.txStack {
		if sc == nil {
			ti = byte(i)
			smc.txStack[i] = c
			break
		}
	}
	if ti == 255 {
		return nil, RpError{CS: 42}
	}

	defer func() {
		smc.txStack[ti] = nil
		close(c)
	}()

	var isNW bool
	switch v := pdu.(type) {
	case Submit:
		isNW = false
		v.TI = ti
		v.RMR = ti
		v.SCA = smc.SCAddress
	case Command:
		isNW = false
		v.TI = ti
		v.RMR = ti
		v.SCA = smc.SCAddress
	case Deliver:
		isNW = true
		v.TI = ti
		v.RMR = ti
		v.SCA = smc.SCAddress
	case StatusReport:
		isNW = true
		v.TI = ti
		v.RMR = ti
		v.SCA = smc.SCAddress
	default:
		return nil, RpError{CS: 97}
	}
	smc.CtrlReq(pdu)

	a := <-c
	switch v := a.(type) {
	case CpAck:
		a = <-c
	case CpError:
		return nil, v
	}

	switch v := a.(type) {
	case SubmitReport:
		if v.RMR != ti {
			smc.CtrlReq(CpError{TI: ti, CS: 81})
			return nil, RpError{CS: 81}
		}
		if !isNW {
			smc.CtrlReq(CpAck{TI: ti})
			return v, nil
		}
	case DeliverReport:
		if v.RMR != ti {
			smc.CtrlReq(CpError{TI: ti, CS: 81})
			return nil, RpError{CS: 81}
		}
		if isNW {
			smc.CtrlReq(CpAck{TI: ti})
			return v, nil
		}
	case RpAckMT:
		if v.RMR != ti {
			smc.CtrlReq(CpError{TI: ti, CS: 81})
			return nil, RpError{CS: 81}
		}
		if !isNW {
			smc.CtrlReq(CpAck{TI: ti})
			return nil, nil
		}
	case RpAckMO:
		if v.RMR != ti {
			smc.CtrlReq(CpError{TI: ti, CS: 81})
			return nil, RpError{CS: 81}
		}
		if isNW {
			smc.CtrlReq(CpAck{TI: ti})
			return nil, nil
		}
	case RpErrorMT:
		if v.RMR != ti {
			smc.CtrlReq(CpError{TI: ti, CS: 81})
			return nil, RpError{CS: 81}
		}
		if !isNW {
			smc.CtrlReq(CpAck{TI: ti})
			return nil, RpError(v)
		}
	case RpErrorMO:
		if v.RMR != ti {
			smc.CtrlReq(CpError{TI: ti, CS: 81})
			return nil, RpError{CS: 81}
		}
		if isNW {
			smc.CtrlReq(CpAck{TI: ti})
			return nil, RpError(v)
		}
	}

	smc.CtrlReq(CpError{TI: ti, CS: 95})
	return nil, RpError{CS: 95}
}
