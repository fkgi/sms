package sms

import (
	"bytes"
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
	l := utf8.RuneCountInString(s)
	o := l%8 + l*8
	ret := make([]byte, (l*7+l%8)/8)

	l = 0
	for i, r := range []rune(s) {
		c := getCode(r)
		if c == 0x80 {
			return nil, &InvalidDataError{
				Name:  "GSM7bit string",
				Bytes: []byte(string(r))}
		}

		shift := uint((o - i) % 8)
		ret[l] = ret[l] | (c << shift)
		if shift > 1 {
			shift = 8 - shift
			ret[l+1] = (c >> shift)
			l++
		} else if shift == 1 {
			l++
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

var mask = []byte{0x01, 0x03, 0x07, 0x0f, 0x1f, 0x3f, 0x7f}

// String return string value of the GSM 7bit String
func (s GSM7bitString) String() string {
	o := (len(s) * 8) % 7
	var b bytes.Buffer
	var next byte

	o = 7 - o
	for i, r := range s {
		shift := uint((i + o) % 7)

		next = next | ((r << shift) & 0x7f)
		b.WriteRune(code[next])
		next = (r >> (7 - shift)) & mask[shift]
		if shift == 6 {
			b.WriteRune(code[next])
			next = 0x00
		}
	}

	return b.String()[1:]
}

// Bytes return byte data
func (s GSM7bitString) Bytes() []byte {
	return []byte(s)
}
