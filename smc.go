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

	CtrlReq func(CPDU)
	// relayInd [2][7]func(RPDU) (RPDU, error)

	TranspInd   func(TPDU) (TPDU, error)
	MemAvailInd func() error
}

// CtrlInd handle CP-DATA/CP-ACK/CP-ERROR
func (smc SMC) CtrlInd(pdu CPDU) {
	var ti byte
	switch v := pdu.(type) {
	case Submit:
		ti = v.TI
		if ti&0x80 != 0x80 && smc.rxStack[ti&0x07] == nil {
			smc.rxStack[ti&0x07] = make(chan CPDU)
			go smc.rxHandler(smc.rxStack[ti&0x07], ti&0x07|0x08)
		}
	case Command:
		ti = v.TI
		if ti&0x80 != 0x80 && smc.rxStack[ti&0x07] == nil {
			smc.rxStack[ti&0x07] = make(chan CPDU)
			go smc.rxHandler(smc.rxStack[ti&0x07], ti&0x07|0x08)
		}
	case Deliver:
		ti = v.TI
		if ti&0x80 != 0x80 && smc.rxStack[ti&0x07] == nil {
			smc.rxStack[ti&0x07] = make(chan CPDU)
			go smc.rxHandler(smc.rxStack[ti&0x07], ti&0x07|0x08)
		}
	case StatusReport:
		ti = v.TI
		if ti&0x80 != 0x80 && smc.rxStack[ti&0x07] == nil {
			smc.rxStack[ti&0x07] = make(chan CPDU)
			go smc.rxHandler(smc.rxStack[ti&0x07], ti&0x07|0x08)
		}
	case SubmitReport:
		ti = v.TI
	case DeliverReport:
		ti = v.TI
	case MemoryAvailable:
		ti = v.TI
	case RpAckMT:
		ti = v.TI
	case RpAckMO:
		ti = v.TI
	case RpErrorMT:
		ti = v.TI
	case RpErrorMO:
		ti = v.TI
	case CpAck:
		ti = v.TI
	case CpError:
		ti = v.TI
	default:
		smc.CtrlReq(CpError{TI: ti, CS: 98})
		return
	}

	var c chan CPDU
	if ti&0x80 == 0x80 {
		c = smc.txStack[ti&0x07]
	} else {
		c = smc.rxStack[ti&0x07]
	}
	if c == nil {
		smc.CtrlReq(CpError{TI: ti, CS: 81})
		return
	}
	c <- pdu
}

func (smc SMC) rxHandler(c chan CPDU, ti byte) {
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
		a, e = smc.TranspInd(v)
	case Command:
	case Deliver:
	case StatusReport:
	case MemoryAvailable:
	case SubmitReport, DeliverReport, RpAckMT, RpAckMO, RpErrorMT, RpErrorMO:
		smc.CtrlReq(CpAck{TI: ti})
		smc.CtrlReq(CpError{TI: ti, CS: 111})
		smc.rxStack[ti&0x07] = nil
		return
	case CpAck:
		smc.rxStack[ti&0x07] = nil
		return
	case CpError:
		smc.rxStack[ti&0x07] = nil
		return
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
		switch v := a.(type) {
		case SubmitReport:
			if isNW {
				v.RMR = mr
				v.TI = ti
				smc.CtrlReq(v)
			}
		case DeliverReport:
			if !isNW {
				v.RMR = mr
				v.TI = ti
				smc.CtrlReq(v)
			}
		case nil:
			a := RpAck{RMR: mr}
			a.TI = ti
			if isNW {
				smc.CtrlReq(RpAckMT(a))
			} else {
				smc.CtrlReq(RpAckMO(a))
			}
		}
	default:
		if isNW {
			v := RpErrorMT{CS: 111}
			v.TI = ti
			v.RMR = mr
			smc.CtrlReq(v)
		} else {
			v := RpErrorMO{CS: 111}
			v.TI = ti
			v.RMR = mr
			smc.CtrlReq(v)
		}
	}

	pdu = <-c
	smc.rxStack[ti&0x07] = nil
	close(c)
	return
}

// TranspReq handle TP-DATA
func (smc SMC) TranspReq(pdu TPDU) (TPDU, error) {
	ti := -1
	c := make(chan CPDU)
	for i, c := range smc.txStack {
		if c == nil {
			ti = i
			smc.txStack[i] = c
			break
		}
	}
	if ti == -1 {
		return nil, RpError{CS: 42}
	}

	var isNW bool
	switch v := pdu.(type) {
	case Submit:
		isNW = false
		v.TI = byte(ti)
		v.RMR = byte(ti)
		v.SCA = smc.SCAddress
		smc.CtrlReq(v)
	case Command:
		isNW = false
		v.TI = byte(ti)
		v.RMR = byte(ti)
		v.SCA = smc.SCAddress
		smc.CtrlReq(v)
	case Deliver:
		isNW = true
		v.TI = byte(ti)
		v.RMR = byte(ti)
		v.SCA = smc.SCAddress
		smc.CtrlReq(v)
	case StatusReport:
		isNW = true
		v.TI = byte(ti)
		v.RMR = byte(ti)
		v.SCA = smc.SCAddress
		smc.CtrlReq(v)
	default:
		smc.txStack[ti] = nil
		close(c)
		return nil, RpError{CS: 97}
	}

	a := <-c
	switch v := a.(type) {
	case DeliverReport:
		if v.RMR != byte(ti) {
			smc.CtrlReq(CpError{TI: byte(ti), CS: 81})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError{CS: 81}
		}
		if isNW {
			smc.CtrlReq(CpAck{TI: byte(ti)})
			smc.txStack[ti] = nil
			close(c)
			return v, nil
		}
	case SubmitReport:
		if v.RMR != byte(ti) {
			smc.CtrlReq(CpError{TI: byte(ti), CS: 81})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError{CS: 81}
		}
		if !isNW {
			smc.CtrlReq(CpAck{TI: byte(ti)})
			smc.txStack[ti] = nil
			close(c)
			return v, nil
		}
	case RpAckMT:
		if v.RMR != byte(ti) {
			smc.CtrlReq(CpError{TI: byte(ti), CS: 81})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError{CS: 81}
		}
		if isNW {
			smc.CtrlReq(CpAck{TI: byte(ti)})
			smc.txStack[ti] = nil
			close(c)
			return nil, nil
		}
	case RpAckMO:
		if v.RMR != byte(ti) {
			smc.CtrlReq(CpError{TI: byte(ti), CS: 81})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError{CS: 81}
		}
		if !isNW {
			smc.CtrlReq(CpAck{TI: byte(ti)})
			smc.txStack[ti] = nil
			close(c)
			return nil, nil
		}
	case RpErrorMT:
		if v.RMR != byte(ti) {
			smc.CtrlReq(CpError{TI: byte(ti), CS: 81})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError{CS: 81}
		}
		if isNW {
			smc.CtrlReq(CpAck{TI: byte(ti)})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError(v)
		}
	case RpErrorMO:
		if v.RMR != byte(ti) {
			smc.CtrlReq(CpError{TI: byte(ti), CS: 81})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError{CS: 81}
		}
		if !isNW {
			smc.CtrlReq(CpAck{TI: byte(ti)})
			smc.txStack[ti] = nil
			close(c)
			return nil, RpError(v)
		}
	}
	smc.CtrlReq(CpAck{TI: byte(ti)})
	smc.txStack[ti] = nil
	close(c)
	return nil, RpError{CS: 95}
}

/*
	ti := -1
	for i, s := range c.state[0] {
		if s == idle {
			ti = i
			break
		}
	}
	if ti == -1 {
		return CpError{TI: byte(ti), CS: 22}
	}
	c.state[0][ti] = connecting

	if c.Connector == nil {
		c.state[0][ti] = waitAck
	} else if c.Connector() {
		c.state[0][ti] = waitAck
	} else {
		c.state[0][ti] = idle
		return CpError{TI: byte(ti), CS: 17}
	}

	switch v := t.(type) {
	case MemoryAvailable:
		v.TI = byte(ti)
		c.TxMM <- v
	case Submit:
		v.TI = byte(ti)
		c.TxMM <- v
	case Command:
		v.TI = byte(ti)
		c.TxMM <- v
	case Deliver:
		v.TI = byte(ti)
		c.TxMM <- v
	case StatusReport:
		v.TI = byte(ti)
		c.TxMM <- v
	default:
		c.state[0][ti] = idle
		return CpError{TI: byte(ti), CS: 95}
	}
	c.state[0][ti] = waitAck
	// start TC1 timer
	a := <-c.RxMM
	// stop TC1 timer

	switch v := a.(type) {
	case CpAck:
	case CpError:
		c.state[0][ti] = idle
		return v
	default:
		c.state[0][ti] = idle
		return CpError{TI: byte(ti), CS: 95}
	}

	a = <-c.RxMM
	switch v := a.(type) {
	case RpAckMO:
	case RpAckMT:
	case SubmitReport:
	case DeliverReport:
	case CpError:
		c.state[0][ti] = idle
		return v
	default:
		c.state[0][ti] = idle
		return CpError{TI: byte(ti), CS: 95}
	}
	return a
}

/*
func (SMC) DataInd() {}
func (SMC) EstablishInd() {}
func (SMC) ErrorInd() {}
*/
