package sms

import (
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
	WriteTo(w io.Writer) (n int64, e error)
}

func (a Address) String() string {
	return fmt.Sprintf("TON/NPI=%d/%d addr=%s", a.TON, a.NPI, a.Addr)
}

// RegexpMatch check matching text of address
func (a Address) RegexpMatch(re *regexp.Regexp) bool {
	return re.MatchString(a.Addr.String())
}

// WriteTo wite binary data to io.Writer
func (a Address) WriteTo(w io.Writer) (n int64, e error) {
	i := 0
	switch a.Addr.(type) {
	case TBCD:
		i = a.Addr.Length()
		if a.TON == 0x05 {
			e = fmt.Errorf("invalid TON for digit address")
			return
		}
	case GSM7bitString:
		i = a.Addr.ByteLength() * 2
		if a.TON != 0x05 || a.NPI != 0x00 {
			e = fmt.Errorf("invalid TON/NPI for alphanumeric address")
			return
		}
	}

	b := []byte{byte(i), 0x80}
	b[1] = b[1] | (a.TON&0x07)<<4
	b[1] = b[1] | (a.NPI & 0x0f)
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	nn, e := a.Addr.WriteTo(w)
	n += nn

	if e == nil && n > 12 {
		e = fmt.Errorf("too much long address data %d", n)
	}
	return
}

// ReadFrom read byte data and set parameter of the Address
func (a *Address) ReadFrom(r io.Reader) (n int64, e error) {
	b := make([]byte, 2)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	l := int(b[0])
	a.TON = (b[1] >> 4) & 0x07
	a.NPI = b[1] & 0x0f

	if l%2 == 1 {
		l++
	}
	b = make([]byte, l/2)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	if a.TON == 0x05 {
		a.Addr = GSM7bitString(b)
	} else {
		a.Addr = TBCD(b)
	}
	return
}
