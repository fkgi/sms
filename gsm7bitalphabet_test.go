package sms

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unicode/utf8"
)

func TestGetGSM7bitString(t *testing.T) {
	txt := "Hello0!"
	bin := []byte{0xC8, 0x32, 0x9B, 0xFD, 0x86, 0x85, 0x00}
	cnv := []byte{}

	if s, e := StringToGSM7bit(txt); e != nil {
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

func TestGSM7bitStringLength(t *testing.T) {
	rand.Seed(time.Now().Unix())
	txt := randText(rand.Int() % 1000)
	t.Logf("\ntest text=%s", strconv.QuoteToGraphic(txt))

	s, e := StringToGSM7bit(txt)
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

func TestGSM7bitStringByteConv(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 50; i++ {
		org := randText(rand.Int() % 1000)
		l := utf8.RuneCountInString(org)
		o := rand.Int() % 8
		t.Logf("\nid=%d, len=%d, offset=%d\norigin=%s",
			i, l, o, strconv.QuoteToGraphic(org))

		s, e := StringToGSM7bit(org)
		if e != nil {
			t.Fatalf("conversion failure: %s", e)
		}
		b := s.marshal(o)
		t.Logf("\nlen=%d\nhex=% x\n", len(b), b)

		s = GSM7bitString(make([]rune, l))
		s.unmarshal(o, b)

		r := s.String()
		if r != org {
			t.Fatalf("\ndetect=%s", strconv.QuoteToGraphic(r))
		}
	}
}

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
