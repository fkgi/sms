package sms

import (
	"fmt"
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

// MarshalRP returns binary data
func (d ErrorMO) MarshalRP() []byte {
	return d.marshalErr(true, nil)
}

// MarshalRP returns binary data
func (d ErrorMT) MarshalRP() []byte {
	return d.marshalErr(false, nil)
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

// UnmarshalRP reads binary data
func (d *ErrorMO) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshalErr(true, b); e != nil && b != nil {
		e = InvalidLengthError{}
	}
	return
}

// UnmarshalErrorMT decode Error MO from bytes
func UnmarshalErrorMT(b []byte) (a ErrorMT, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP reads binary data
func (d *ErrorMT) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshalErr(false, b); e != nil && b != nil {
		e = InvalidLengthError{}
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
