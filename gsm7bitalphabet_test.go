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
	s, e := GetGSM7bitString("Design@Home")
	if e != nil {
		t.Fatalf("conversion failure: %s", e)
	}
	b := []byte{0xC4, 0xF2, 0x3C, 0x7D, 0x76, 0x03, 0x90, 0xEF, 0x76, 0x19}

	for i := range b {
		if s[i] != b[i] {
			t.Fatalf("conversion failed in %d byte", i)
		}
	}
}

func TestGSM7bitStringLength(t *testing.T) {
	org := randText(200)
	s, e := GetGSM7bitString(org)
	if e != nil {
		t.Fatalf("conversion failure: %s", e)
	}

	if s.Length() != utf8.RuneCountInString(org) {
		t.Fatalf("detect length failed orig=%d detect=%d",
			utf8.RuneCountInString(org), s.Length())
	}
}

func TestGSM7bitStringString(t *testing.T) {
	org := randText(200)
	s, e := GetGSM7bitString(org)
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

func randText(len int) string {
	rand.Seed(time.Now().Unix())
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		b.WriteRune(code[rand.Int()%128])
	}
	return b.String()
}
