package sms

import (
	"testing"
	"time"
)

func TestEncodeDeliver(t *testing.T) {
	bytedata := []byte{
		0x40, 0x04, 0x80, 0x21, 0x43, 0x00, 0x08, 0x11,
		0x30, 0x22, 0x41, 0x52, 0x04, 0x63, 0x10, 0x05,
		0x00, 0x03, 0x87, 0x02, 0x01, 0x30, 0x42, 0x30,
		0x44, 0x30, 0x46, 0x30, 0x48, 0x30, 0x4a}
	p, e := DecodeAsMS(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeDeliver(t *testing.T) {
	p := &Deliver{
		MMS: true,
		LP:  false,
		SRI: false,
		RP:  false,
		OA:  Address{TON: 0, NPI: 0},
		PID: 0,
		DCS: &GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   NoMessageClass,
			Charset:    UCS2},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		UDH: []UDH{&ConcatenatedSM{0x84, 0x0a, 0x01}}}
	p.OA.Addr, _ = ParseTBCD("1234")
	_, p.UD, _ = p.DCS.Encode("あいうえお")

	b := p.Encode()
	t.Logf("% x", b)
}

func TestEncodeDeliverReport(t *testing.T) {
	bytedata := []byte{
		0x00, 0x00}
	p, e := DecodeAsSC(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeDeliverReport(t *testing.T) {
	p := &DeliverReport{
		PID: nil,
		DCS: nil,
		UDH: nil,
		UD:  nil}

	b := p.Encode()
	t.Logf("% x", b)
}
