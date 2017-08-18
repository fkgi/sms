package sms

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

// Address is SMS originator/destination address
type Address struct {
	TON  byte
	NPI  byte
	Addr addrValue
}

type addrValue interface {
	Length() int
	ByteLength() int
	String() string
	Bytes() []byte
	//	WriteTo(w io.Writer) (n int64, e error)
}

func (a Address) String() string {
	return fmt.Sprintf("TON/NPI=%d/%d addr=%s", a.TON, a.NPI, a.Addr)
}

// RegexpMatch check matching text of address
func (a Address) RegexpMatch(re *regexp.Regexp) bool {
	return re.MatchString(a.Addr.String())
}

// Encode generate binary data of this Address
func (a Address) Encode() (l byte, b []byte) {
	switch a.Addr.(type) {
	case TBCD:
		l = byte(a.Addr.Length())
		if a.TON == 0x05 {
			a.TON = 0x00
		}
	case GSM7bitString:
		l = byte(a.Addr.ByteLength() * 2)
		a.TON = 0x05
		a.NPI = 0x00
	}

	b = []byte{0x80}
	b[0] = b[0] | (a.TON&0x07)<<4
	b[0] = b[0] | (a.NPI & 0x0f)

	if a.Addr != nil {
		b = append(b, a.Addr.Bytes()...)
	}
	return
}

// Decode make Address from binary data
func (a *Address) Decode(b []byte) {
	a.TON = (b[0] >> 4) & 0x07
	a.NPI = b[0] & 0x0f
	if a.TON == 0x05 {
		a.Addr = GSM7bitString(b[1:])
	} else {
		a.Addr = TBCD(b[1:])
	}
	return
}

func readAddr(r *bytes.Reader) (a Address, e error) {
	var l byte
	if l, e = r.ReadByte(); e != nil {
		return
	}
	if l%2 == 1 {
		l++
	}
	b := make([]byte, l/2+1)
	var i int
	if i, e = r.Read(b); e == nil {
		if i != len(b) {
			e = io.EOF
		} else {
			a.Decode(b)
		}
	}
	return
}
