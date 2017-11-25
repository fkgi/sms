package sms_test

import (
	"testing"
	"time"

	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func TestEncodeDeliver(t *testing.T) {
	bytedata := []byte{
		0x40, 0x04, 0x80, 0x21, 0x43, 0x00, 0x08, 0x11,
		0x30, 0x22, 0x41, 0x52, 0x04, 0x63, 0x10, 0x05,
		0x00, 0x03, 0x87, 0x02, 0x01, 0x30, 0x42, 0x30,
		0x44, 0x30, 0x46, 0x30, 0x48, 0x30, 0x4a}
	p, e := sms.DecodeAsMS(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeDeliver(t *testing.T) {
	p := &sms.Deliver{
		MMS: true,
		LP:  false,
		SRI: false,
		RP:  false,
		OA:  sms.Address{TON: 0, NPI: 0},
		PID: 0,
		DCS: &sms.GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   sms.NoMessageClass,
			Charset:    sms.CharsetUCS2},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		UD: sms.UD{Text: "あいうえお"}}
	p.UD.AddUDH(&sms.ConcatenatedSM{
		RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
	p.OA.Addr, _ = teldata.ParseTBCD("1234")

	b := p.Encode()
	t.Logf("% x", b)
}

func TestEncodeDeliverReport(t *testing.T) {
	bytedata := []byte{
		0x00, 0x00}
	p, e := sms.DecodeAsSC(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeDeliverReport(t *testing.T) {
	p := &sms.DeliverReport{
		PID: nil,
		DCS: nil}

	b := p.Encode()
	t.Logf("% x", b)
}
