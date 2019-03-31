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

// Address is SMS originator/destination address
type Address struct {
	TON  byte
	NPI  byte
	Addr addrValue
}

type jaddr struct {
	TON  byte   `json:"ton"`
	NPI  byte   `json:"npi"`
	Addr string `json:"addr"`
}

func (a Address) String() string {
	if a.Addr == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"TON/NPI=%d/%d addr=%s", a.TON, a.NPI, a.Addr)
}

// MarshalJSON provide custom marshaller
func (a Address) MarshalJSON() ([]byte, error) {
	addr := jaddr{
		TON: a.TON,
		NPI: a.NPI}
	if a.Addr != nil {
		addr.Addr = a.Addr.String()
	}
	return json.Marshal(addr)
}

// UnmarshalJSON provide custom marshaller
func (a *Address) UnmarshalJSON(b []byte) (e error) {
	addr := jaddr{}
	if e = json.Unmarshal(b, &addr); e != nil {
		return e
	}
	a.TON = addr.TON
	a.NPI = addr.NPI

	if len(addr.Addr) == 0 {
		a.Addr = nil
	} else if a.TON == 0x05 {
		a.Addr, e = StringToGSM7bit(addr.Addr)
	} else {
		a.Addr, e = teldata.ParseTBCD(addr.Addr)
	}

	return e
}

// RegexpMatch check matching text of address
func (a Address) RegexpMatch(re *regexp.Regexp) bool {
	return re.MatchString(a.Addr.String())
}

// Encode generate binary data and semi-octet length of this Address
func (a Address) Encode() (l byte, b []byte) {
	switch a.Addr.(type) {
	case teldata.TBCD:
		l = byte(a.Addr.Length())
		if a.TON == 0x05 {
			a.TON = 0x00
		}
	case GSM7bitString:
		l = byte(a.Addr.Length() * 7 / 4)
		if a.Addr.Length()*7%4 != 0 {
			l++
		}
		a.TON = 0x05
		a.NPI = 0x00
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

// Decode make Address from binary data and semi-octet length
func (a *Address) Decode(l byte, b []byte) {
	a.TON = (b[0] >> 4) & 0x07
	a.NPI = b[0] & 0x0f

	b = b[1:]
	if a.TON == 0x05 {
		l = l * 4 / 7
		a.Addr = DecodeGSM7bit(int(l), b)
	} else {
		if l%2 == 1 {
			b[len(b)-1] |= 0xf0
		}
		a.Addr = teldata.TBCD(b)
	}
	return
}

func readAddr(r *bytes.Reader) (Address, error) {
	a := Address{}
	l, e := r.ReadByte()
	if e != nil {
		return a, e
	}

	b := make([]byte, l/2+l%2+1)
	var i int
	if i, e = r.Read(b); e == nil {
		if i != len(b) {
			e = io.EOF
		} else {
			a.Decode(l, b)
		}
	}
	return a, e
}
