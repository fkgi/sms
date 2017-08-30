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
	bin := []byte{0x00, 0x64, 0x99, 0xCD, 0x7E, 0xC3, 0x42}
	cnv := []byte{}

	t.Logf("\ntest text=%s", txt)
	if s, e := GetGSM7bitString(txt); e != nil {
		t.Fatalf("conversion failure: %s", e)
	} else {
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

	s, e := GetGSM7bitString(txt)
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
		t.Logf("\nid=%d, len=%d, bin=%d, offset=%d\norigin=%s",
			i, l, (l*7+l%8)/8, l%8, strconv.QuoteToGraphic(org))

		s, e := GetGSM7bitString(org)
		if e != nil {
			t.Fatalf("conversion failure: %s", e)
		}
		b := s.Bytes()
		t.Logf("\nhex=% x\n", b)

		s = GetGSM7bitByte(l, b)
		r := s.String()
		if r != org {
			t.Fatalf("\ndetect=%s", strconv.QuoteToGraphic(r))
		}
	}
}

func randText(len int) string {
	//rand.Seed(time.Now().Unix())
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		b.WriteRune(code[rand.Int()%128])
	}
	return b.String()
}
