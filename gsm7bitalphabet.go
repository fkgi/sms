package sms

import (
	"unicode/utf8"
)

// GSM7bitString is GSM 7-bit default alphabet of 3GPP TS23.038
type GSM7bitString []rune

var code = [128 + 16]rune{
	'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç',
	'\n', 'Ø', 'ø', '\r', 'Å', 'å', 'Δ', '_', 'Φ', 'Γ',
	'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', '\x1b', 'Æ', 'æ',
	'ß', 'É', ' ', '!', '"', '#', '¤', '%', '&', '\'',
	'(', ')', '*', '+', ',', '-', '.', '/', '0', '1',
	'2', '3', '4', '5', '6', '7', '8', '9', ':', ';',
	'<', '=', '>', '?', '¡', 'A', 'B', 'C', 'D', 'E',
	'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y',
	'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§', '¿', 'a', 'b', 'c',
	'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
	'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
	'|', '\x00', '\x00', '\x00', '^', '€', '\x00', '\x00',
	'{', '}', '\f', '\x00', '[', '~', ']', '\\'}

func getCode(c rune) (bool, byte) {
	if (c > 0x61 && c < 0x7a) || (c > 0x41 && c < 0x5a) ||
		(c > 0x20 && c < 0x23) || (c > 0x25 && c < 0x3f) {
		return false, byte(c)
	}
	if c == '\x00' {
		return false, 0xff
	}

	for i, r := range code {
		if r != c {
			continue
		}
		if i < 0x80 {
			return false, byte(i)
		}
		switch i {
		case 0x80:
			return true, 0x40
		case 0x84:
			return true, 0x14
		case 0x85:
			return true, 0x65
		case 0x88:
			return true, 0x28
		case 0x89:
			return true, 0x29
		case 0x8a:
			return true, 0x0a
		case 0x8c:
			return true, 0x3c
		case 0x8d:
			return true, 0x3d
		case 0x8e:
			return true, 0x3e
		case 0x8f:
			return true, 0x2f
		}
	}
	return false, 0xff
}

// StringToGSM7bit generate GSM7bitString from string
func StringToGSM7bit(s string) (GSM7bitString, error) {
	l := utf8.RuneCountInString(s)
	txt := make([]rune, l)
	i := 0
	for _, r := range s {
		_, c := getCode(r)
		if c == 0xff {
			return nil, UnknownGSM7bitRuneError{R: r}
		}
		txt[i] = r
		i++
	}
	return txt, nil
}

// UnmarshalGSM7bitString generate GSM7bitString from byte slice with offset
func UnmarshalGSM7bitString(o, l int, b []byte) GSM7bitString {
	s := GSM7bitString(make([]rune, 0, l))

	o = 7 - o
	var next byte
	var sh uint
	var esc bool
	var sep int

	for i, r := range b {
		sh = uint((i + o) % 7)
		next |= (r << sh) & 0x7f
		if i != 0 || o == 7 {
			if next == 0x1b {
				esc = true
			} else if esc {
				s = append(s, code[next&0x0f|0x80])
				sep += 2
				esc = false
			} else {
				s = append(s, code[next])
				sep++
			}
		}

		sh = 7 - sh
		next = (r >> sh) & (0x7f >> (sh - 1))
		if sh == 1 && sep < l {
			if next == 0x1b {
				esc = true
			} else if esc {
				s = append(s, code[next&0x0f|0x80])
				sep += 2
				esc = false
			} else {
				s = append(s, code[next])
				sep++
			}
			next = 0x00
		}
	}
	return s
}

// Equal reports a and b are same
func (s GSM7bitString) Equal(b GSM7bitString) bool {
	return s.String() == b.String()
}

// Length return length of the GSM 7bit String
func (s GSM7bitString) Length() int {
	i := 0
	for _, c := range s {
		i++
		if esc, _ := getCode(c); esc {
			i++
		}
	}
	return i
}

func (s GSM7bitString) String() string {
	if s == nil {
		return "<nil>"
	}
	return string(s)
}

// Bytes return byte data
func (s GSM7bitString) Bytes() []byte {
	return s.Marshal(0)
}

// Marshal return byte data with offset shift
func (s GSM7bitString) Marshal(o int) []byte {
	l := s.Length()*7 + o
	b := make([]byte, l/8+1)
	if l%8 == 0 {
		b = b[:l/8]
	} else if l%8 == 1 {
		b[l/8] = 0x1a
	}

	var sh uint
	var c byte
	var esc bool
	i := 0
	o += l * 8
	l = 0
	f := func(code byte) {
		sh = uint((o - i) % 8)
		b[l] |= code << sh
		if sh > 1 {
			sh = 8 - sh
			b[l+1] = code >> sh
			l++
		} else if sh == 1 {
			l++
		}
		i++
	}
	for _, r := range s {
		esc, c = getCode(r)
		if c == 0xff {
			c = 0x20
		}

		if esc {
			f(0x1b)
		}
		f(c)
	}
	return b
}

func (s GSM7bitString) trim(l int) GSM7bitString {
	r := make([]rune, 0, len(s))
	i := 0
	for _, c := range s {
		i++
		if esc, _ := getCode(c); esc {
			i++
		}
		if i > l {
			break
		}
		r = append(r, c)
	}
	return GSM7bitString(r)
}
