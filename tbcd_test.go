package sms

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unicode/utf8"
)

var digit = [15]rune{
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', '*', '#', 'a', 'b', 'c'}

func TestParseTBCD(t *testing.T) {
	s, e := ParseTBCD("0123456789*#abc")
	if e != nil {
		t.Fatalf("conversion failure: %s", e)
	}
	b := []byte{0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe}

	for i := range b {
		if s[i] != b[i] {
			t.Fatalf("conversion failed in %d byte", i)
		}
	}
}

func TestTBCDLength(t *testing.T) {
	rand.Seed(time.Now().Unix())
	org := randDigit(rand.Int() % 100)
	s, e := ParseTBCD(org)
	if e != nil {
		t.Fatalf("conversion failure: %s", e)
	}

	if s.Length() != utf8.RuneCountInString(org) {
		t.Fatalf("detect length failed orig=%d detect=%d",
			utf8.RuneCountInString(org), s.Length())
	}
}

func TestTBCDString(t *testing.T) {
	rand.Seed(time.Now().Unix())
	org := randDigit(rand.Int() % 100)
	s, e := ParseTBCD(org)
	if e != nil {
		t.Fatalf("conversion failure: %s", e)
	}

	r := s.String()
	if r != org {
		t.Fatalf("different text orig=%s detect=%s",
			strconv.QuoteToGraphic(org),
			strconv.QuoteToGraphic(r))
	}
}

func randDigit(len int) string {
	rand.Seed(time.Now().Unix())
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		b.WriteRune(digit[rand.Int()%15])
	}
	return b.String()
}
