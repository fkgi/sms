package sms

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// TBCD is Telephony Binary Coded Decimal value
type TBCD []byte

// ParseTBCD create TBCD value from string
func ParseTBCD(s string) (TBCD, error) {
	if strings.ContainsRune(s, '\x00') {
		return nil, fmt.Errorf("invalid charactor")
	} else if len(s)%2 != 0 {
		s = s + "\x00"
	}

	r := make([]byte, len(s)/2)
	for i, c := range s {
		var v byte
		switch c {
		case '0':
			v = 0x00
		case '1':
			v = 0x01
		case '2':
			v = 0x02
		case '3':
			v = 0x03
		case '4':
			v = 0x04
		case '5':
			v = 0x05
		case '6':
			v = 0x06
		case '7':
			v = 0x07
		case '8':
			v = 0x08
		case '9':
			v = 0x09
		case '*':
			v = 0x0a
		case '#':
			v = 0x0b
		case 'a', 'A':
			v = 0x0c
		case 'b', 'B':
			v = 0x0d
		case 'c', 'C':
			v = 0x0e
		case '\x00':
			v = 0x0f
		default:
			return r, fmt.Errorf("invalid charactor %c", c)
		}
		if i%2 == 1 {
			v = v << 4
		}
		r[i/2] = r[i/2] | v
	}
	return r, nil
}

// Length return length of the TBCD digit
func (t TBCD) Length() int {
	ret := len(t) * 2
	if (t[len(t)-1]&0xf0)>>4 == 0x0f {
		ret--
	}
	return ret
}

// ByteLength return octets of the TBCD digit
func (t TBCD) ByteLength() int {
	return len(t)
}

// String return string value of the TBCD digit
func (t TBCD) String() string {
	var b bytes.Buffer
	so := [2]byte{}
	for _, c := range t {
		so[0] = c & 0x0f
		so[1] = (c & 0xf0) >> 4
		for _, s := range so {
			switch s {
			case 0x00:
				b.WriteRune('0')
			case 0x01:
				b.WriteRune('1')
			case 0x02:
				b.WriteRune('2')
			case 0x03:
				b.WriteRune('3')
			case 0x04:
				b.WriteRune('4')
			case 0x05:
				b.WriteRune('5')
			case 0x06:
				b.WriteRune('6')
			case 0x07:
				b.WriteRune('7')
			case 0x08:
				b.WriteRune('8')
			case 0x09:
				b.WriteRune('9')
			case 0x0a:
				b.WriteRune('*')
			case 0x0b:
				b.WriteRune('#')
			case 0x0c:
				b.WriteRune('a')
			case 0x0d:
				b.WriteRune('b')
			case 0x0e:
				b.WriteRune('c')
			case 0x0f:
			}
		}
	}
	return b.String()
}

// WriteTo wite binary data to io.Writer
func (t TBCD) WriteTo(w io.Writer) (n int64, e error) {
	i := 0
	i, e = w.Write(t)
	return int64(i), e
}
