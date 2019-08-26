package sms_test

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/fkgi/sms"
)

func TestGetGSM7bitString(t *testing.T) {
	txt := "Hello0!"
	bin := []byte{0xC8, 0x32, 0x9B, 0xFD, 0x86, 0x85, 0x00}
	cnv := []byte{}

	if s, e := sms.StringToGSM7bit(txt); e != nil {
		t.Fatalf("conversion failure: %s", e)
	} else {
		t.Logf("\ntest text=%s", s)
		cnv = s.Bytes()
	}

	t.Logf("\nconv bin=% x\norig bin=% x", cnv, bin)
	for i := range cnv {
		if bin[i] != cnv[i] {
			t.Fatalf("\nconversion failed in %d byte,"+
				" %x should be %x", i, cnv[i], bin[i])
		}
	}
}

/*
func TestGSM7bitStringLength(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 500; i++ {
		txt := randText(rand.Int() % 1000)
		t.Logf("\ntest text=%s", strconv.QuoteToGraphic(txt))

		s, e := sms.StringToGSM7bit(txt)
		if e != nil {
			t.Fatalf("conversion failure: %s", e)
		}
		t.Logf("\norig len=%d\nconv len=%d",
			utf8.RuneCountInString(txt), s.Length())

		if s.Length() != utf8.RuneCountInString(txt) {
			t.Fatalf("detect length failed orig=%d detect=%d",
				utf8.RuneCountInString(txt), s.Length())
		}
	}
}
*/

func TestGSM7bitStringByteConv(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		org := randText(rand.Int() % 1000)
		l := utf8.RuneCountInString(org)
		o := rand.Int() % 8
		t.Logf("\nid=%d, len=%d, offset=%d\norigin=%s",
			i, l, o, strconv.QuoteToGraphic(org))

		s, e := sms.StringToGSM7bit(org)
		if e != nil {
			t.Fatalf("conversion failure: %s", e)
		}
		b := s.Marshal(o)
		t.Logf("\nlen=%d\nhex=% x\n", len(b), b)

		r := sms.UnmarshalGSM7bitString(o, s.Length(), b).String()
		if r != org {
			t.Fatalf("\ndetect=%s", strconv.QuoteToGraphic(r))
		}
	}
}

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

func randText(len int) string {
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		c := 0x1b
		for code[c] == '\x00' || code[c] == '\x1b' {
			c = rand.Int() % (128 + 16)
		}
		b.WriteRune(code[c])
	}
	return b.String()
}

func randDigit(len int) string {
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		b.WriteString(strconv.Itoa(rand.Int() % 10))
	}
	return b.String()
}
