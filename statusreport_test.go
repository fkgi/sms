package sms

import (
	"testing"
	"time"
)

var bytedata = []byte{
	0x06, 0x00, 0x04, 0x80, 0x21, 0x43, 0x11, 0x30,
	0x22, 0x41, 0x52, 0x04, 0x63, 0x11, 0x30, 0x22,
	0x41, 0x52, 0x04, 0x63, 0x00}

func TestEncodeStatusReport(t *testing.T) {
	p, e := DecodeAsMS(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeStatusReport(t *testing.T) {
	p := &StatusReport{
		MMS: false,
		LP:  false,
		SRQ: false,
		MR:  0x00,
		RA:  Address{TON: 0, NPI: 0},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		DT: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		ST: 0x00}
	p.RA.Addr, _ = ParseTBCD("1234")

	b := p.Encode()
	t.Logf("% x", b)
}
