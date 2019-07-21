package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/fkgi/teldata"
)

type addrValue interface {
	Length() int
	String() string
	Bytes() []byte
}

const (
	// TypeUnknown of TON
	TypeUnknown byte = 0
	// TypeInternational of TON
	TypeInternational byte = 1
	// TypeNational of TON
	TypeNational byte = 2
	// TypeNetworkSpecific of TON
	TypeNetworkSpecific byte = 3
	// TypeSubscriber of TON
	TypeSubscriber byte = 4
	// TypeAlphanumeric of TON
	TypeAlphanumeric byte = 5
	// TypeAbbreviated of TON
	TypeAbbreviated byte = 6

	// PlanUnknown of NPI
	PlanUnknown byte = 0
	// PlanISDNTelephone of NPI
	PlanISDNTelephone byte = 1
	// PlanData of NPI
	PlanData byte = 3
	// PlanTelex of NPI
	PlanTelex byte = 4
	// PlanSCSpecific1 of NPI
	PlanSCSpecific1 byte = 5
	// PlanSCSpecific2 of NPI
	PlanSCSpecific2 byte = 6
	// PlanNational of NPI
	PlanNational byte = 8
	// PlanPrivate of NPI
	PlanPrivate byte = 9
	// PlanERMES of NPI
	PlanERMES byte = 10
)

// Address is SMS originator/destination address
type Address struct {
	TON  byte      `json:"ton"`
	NPI  byte      `json:"npi"`
	Addr addrValue `json:"addr"`
}

func (a Address) String() string {
	if a.Addr == nil {
		return fmt.Sprintf(
			"TON/NPI=%d/%d addr=<empty>", a.TON, a.NPI)
	}
	return fmt.Sprintf(
		"TON/NPI=%d/%d addr=%s", a.TON, a.NPI, a.Addr)
}

// UnmarshalJSON provide custom marshaller
func (a *Address) UnmarshalJSON(b []byte) (e error) {
	type alias Address
	al := struct {
		Addr string `json:"addr"`
		*alias
	}{
		alias: (*alias)(a),
	}
	if e = json.Unmarshal(b, &al); e != nil {
		return
	}

	if len(al.Addr) == 0 {
		a.Addr = nil
	} else if a.TON == 0x05 {
		a.Addr, e = StringToGSM7bit(al.Addr)
	} else {
		a.Addr, e = teldata.ParseTBCD(al.Addr)
	}
	return
}

// MarshalJSON provide custom marshaller
func (a Address) MarshalJSON() ([]byte, error) {
	type alias Address
	al := struct {
		*alias
		Addr string `json:"addr"`
	}{
		alias: (*alias)(&a)}
	if a.Addr == nil {
		al.Addr = ""
	} else {
		al.Addr = a.Addr.String()
	}
	return json.Marshal(al)
}

// RegexpMatch check matching text of address
func (a Address) RegexpMatch(re *regexp.Regexp) bool {
	return re.MatchString(a.Addr.String())
}

// Marshal generate binary data and semi-octet length of this Address
func (a Address) Marshal() (l byte, b []byte) {
	switch a.Addr.(type) {
	case teldata.TBCD:
		l = byte(a.Addr.Length())
		if a.TON == TypeAlphanumeric {
			a.TON = TypeUnknown
		}
	case GSM7bitString:
		l = byte(a.Addr.Length() * 7 / 4)
		if a.Addr.Length()*7%4 != 0 {
			l++
		}
		a.TON = TypeAlphanumeric
		a.NPI = PlanUnknown
	default:
		// null addr
	}

	b = []byte{0x80}
	b[0] |= (a.TON & 0x07) << 4
	b[0] |= a.NPI & 0x0f

	if a.Addr != nil {
		b = append(b, a.Addr.Bytes()...)
	}
	return
}

// Unmarshal make Address from binary data and semi-octet length
func (a *Address) Unmarshal(l byte, b []byte) {
	if a == nil {
		return
	}

	a.TON = (b[0] >> 4) & 0x07
	a.NPI = b[0] & 0x0f

	b = b[1:]
	if a.TON == TypeAlphanumeric {
		l = l * 4 / 7
		a.Addr = UnmarshalGSM7bit(int(l), b)
	} else {
		if l%2 == 1 {
			b[len(b)-1] |= 0xf0
		}
		a.Addr = teldata.TBCD(b)
	}
	return
}

func readAddr(r *bytes.Reader) (a Address, e error) {
	var l byte
	if l, e = r.ReadByte(); e != nil {
		return
	}

	b := make([]byte, l/2+l%2+1)
	var i int
	if i, e = r.Read(b); e == nil {
		if i != len(b) {
			e = io.EOF
		} else {
			a.Unmarshal(l, b)
		}
	}
	return
}
