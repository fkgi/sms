package sms

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

// GSM7bitString is GSM 7-bit default alphabet of 3GPP TS23.038
type GSM7bitString []byte

var code = [128]rune{
	'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
	'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', '\x1b', 'Æ', 'æ', 'ß', 'É',
	' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
	'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
	'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à'}

func getCode(c rune) byte {
	for i, r := range code {
		if r == c {
			return byte(i)
		}
	}
	return 0x80
}

// GetGSM7bitString generate GSM7bitString from string
func GetGSM7bitString(s string) (GSM7bitString, error) {
	ret := make([]byte, 0, utf8.RuneCountInString(s)*7/8+1)
	for i, r := range []rune(s) {
		c := getCode(r)
		if c == 0x80 {
			return nil, fmt.Errorf("invalid character %c", r)
		}

		switch i % 8 {
		case 0:
			ret = append(ret, c&0x7f)
		case 1:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 7) & 0x80)
			ret = append(ret, (c>>1)&0x3f)
		case 2:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 6) & 0xc0)
			ret = append(ret, (c>>2)&0x1f)
		case 3:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 5) & 0xe0)
			ret = append(ret, (c>>3)&0x0f)
		case 4:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 4) & 0xf0)
			ret = append(ret, (c>>4)&0x07)
		case 5:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 3) & 0xf8)
			ret = append(ret, (c>>5)&0x03)
		case 6:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 2) & 0xfc)
			ret = append(ret, (c>>6)&0x01)
		case 7:
			ret[len(ret)-1] = (ret[len(ret)-1]) | ((c << 1) & 0xfe)
		}
	}
	return ret, nil
}

// Length return length of the GSM 7bit String
func (s GSM7bitString) Length() int {
	bit := len(s) * 8
	bit = bit - bit%7
	return bit / 7
}

// ByteLength return octets of the GSM 7bit String
func (s GSM7bitString) ByteLength() int {
	return len(s)
}

// String return string value of the GSM 7bit String
func (s GSM7bitString) String() string {
	var b bytes.Buffer
	for i, r := range s {
		switch i % 7 {
		case 0:
			b.WriteRune(code[r&0x7f])
		case 6:
			b.WriteRune(code[(r>>1)&0x7f])
		}
		if i+1 != len(s) {
			switch i % 7 {
			case 0:
				b.WriteRune(code[((r>>7)&0x01)|((s[i+1]<<1)&0x7e)])
			case 1:
				b.WriteRune(code[((r>>6)&0x03)|((s[i+1]<<2)&0x7c)])
			case 2:
				b.WriteRune(code[((r>>5)&0x07)|((s[i+1]<<3)&0x78)])
			case 3:
				b.WriteRune(code[((r>>4)&0x0f)|((s[i+1]<<4)&0x70)])
			case 4:
				b.WriteRune(code[((r>>3)&0x1f)|((s[i+1]<<5)&0x60)])
			case 5:
				b.WriteRune(code[((r>>2)&0x3f)|((s[i+1]<<6)&0x40)])
			}
		}
	}
	return b.String()
}

// Bytes return byte data
func (s GSM7bitString) Bytes() []byte {
	return []byte(s)
}
