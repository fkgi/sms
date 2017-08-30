package sms

import (
	"testing"
	"time"
)

func TestEncodeSubmit(t *testing.T) {

	bytedata := []byte{
		0x41, 0x40, 0x0b, 0x81, 0x90, 0x10, 0x32, 0x54,
		0x76, 0xf8, 0x00, 0x08, 0x10, 0x05, 0x00, 0x03,
		0x84, 0x0a, 0x01, 0x30, 0x42, 0x30, 0x44, 0x30,
		0x46, 0x30, 0x48, 0x30, 0x4a}
	p, e := DecodeAsSC(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}

	t.Log(p.String())
}

func TestDecodeSubmit(t *testing.T) {
	p := &Submit{
		RD:  false,
		SRR: false,
		RP:  false,
		MR:  64,
		DA:  Address{TON: 0, NPI: 1},
		PID: 0,
		DCS: &GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   NoMessageClass,
			Charset:    UCS2},
		VP:  nil,
		UDH: []UDH{&ConcatenatedSM{0x84, 0x0a, 0x01}}}
	p.DA.Addr, _ = ParseTBCD("09012345678")
	_, p.UD, _ = p.DCS.Encode("あいうえお")

	b := p.Encode()
	t.Logf("% x", b)
}

func TestEncodeSubmitReport(t *testing.T) {
	bytedata := []byte{
		0x01, 0x00, 0x11, 0x30, 0x22, 0x41, 0x52, 0x04, 0x63}
	p, e := DecodeAsMS(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeSubmitReport(t *testing.T) {
	p := &SubmitReport{
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		PID: nil,
		DCS: nil,
		UDH: nil,
		UD:  nil}

	b := p.Encode()
	t.Logf("% x", b)
}
