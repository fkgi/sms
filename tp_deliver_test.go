package sms_test

import (
	"encoding/json"
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
	p, e := sms.UnmarshalMTTP(bytedata)
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

	b := p.MarshalTP()
	t.Logf("% x", b)
}

func TestMarshalJSON_deliver(t *testing.T) {
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

	t.Log(p.String())

	var e error
	bytedata, e = json.Marshal(p)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))

	if e := json.Unmarshal(bytedata, &p); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(p.String())
}

func TestEncodeDeliverReport(t *testing.T) {
	bytedata := []byte{
		0x00, 0x00}
	p, e := sms.UnmarshalMOTP(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeDeliverReport(t *testing.T) {
	p := &sms.DeliverReport{
		PID: nil,
		DCS: nil}

	b := p.MarshalTP()
	t.Logf("% x", b)
}

func TestMarshalJSON_deliverreport(t *testing.T) {
	p := &sms.DeliverReport{
		FCS: 0xC0,
		DCS: &sms.GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   sms.NoMessageClass,
			Charset:    sms.CharsetUCS2},
		UD: sms.UD{Text: "あいうえお"}}
	tmp := byte(0x01)
	p.PID = &tmp
	t.Log(p.String())

	var e error
	bytedata, e = json.Marshal(p)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))

	if e := json.Unmarshal(bytedata, &p); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(p.String())
}
