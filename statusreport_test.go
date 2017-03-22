package sms

import (
	"bytes"
	"testing"
	"time"
)

func TestEncodeStatusReport(t *testing.T) {

	bytedata := []byte{
		0x06, 0x00, 0x04, 0x80, 0x21, 0x43, 0x11, 0x30,
		0x22, 0x41, 0x52, 0x04, 0x63, 0x11, 0x30, 0x22,
		0x41, 0x52, 0x04, 0x63, 0x00}
	buf := bytes.NewBuffer(bytedata)
	p, _, e := Read(buf, false)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	b := new(bytes.Buffer)
	p.PrintStack(b)
	t.Log(b.String())
}

func TestDecodeStatusReport(t *testing.T) {
	p := &StatusReport{}
	p.MMS = false
	p.LP = false
	p.SRQ = false
	p.MR = 0x00
	p.RA = Address{TON: 0, NPI: 0}
	p.RA.Addr, _ = ParseTBCD("1234")
	p.SCTS = time.Date(
		2011, time.March, 22, 14, 25, 40, 0,
		time.FixedZone("unknown", 9*60*60))
	p.DT = time.Date(
		2011, time.March, 22, 14, 25, 40, 0,
		time.FixedZone("unknown", 9*60*60))
	p.ST = 0x00

	b := new(bytes.Buffer)
	_, e := p.WriteTo(b)
	if e != nil {
		t.Fatalf("deecode failed: %s", e)
	}

	t.Logf("% x", b)
}
