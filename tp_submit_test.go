package sms_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func TestEncodeSubmit(t *testing.T) {

	bytedata := []byte{
		0x41, 0x40, 0x0b, 0x81, 0x90, 0x10, 0x32, 0x54,
		0x76, 0xf8, 0x00, 0x08, 0x10, 0x05, 0x00, 0x03,
		0x84, 0x0a, 0x01, 0x30, 0x42, 0x30, 0x44, 0x30,
		0x46, 0x30, 0x48, 0x30, 0x4a}
	p, e := sms.DecodeAsSC(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}

	t.Log(p.String())
}

func TestDecodeSubmit(t *testing.T) {
	p := &sms.Submit{
		RD:  false,
		SRR: false,
		RP:  false,
		MR:  64,
		DA:  sms.Address{TON: 0, NPI: 1},
		PID: 0,
		DCS: &sms.GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   sms.NoMessageClass,
			Charset:    sms.CharsetUCS2},
		VP: nil,
		UD: sms.UD{Text: "あいうえお"}}
	p.UD.AddUDH(&sms.ConcatenatedSM{
		RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
	p.DA.Addr, _ = teldata.ParseTBCD("09012345678")

	b := p.Encode()
	t.Logf("% x", b)
}

func TestMarshalJSON_submit(t *testing.T) {
	p := &sms.Submit{
		RD:  false,
		SRR: false,
		RP:  false,
		MR:  64,
		DA:  sms.Address{TON: 0, NPI: 1},
		PID: 0,
		DCS: &sms.GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   sms.NoMessageClass,
			Charset:    sms.CharsetUCS2},
		VP: nil,
		UD: sms.UD{Text: "あいうえお"}}
	p.DA.Addr, _ = teldata.ParseTBCD("09012345678")

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

func TestEncodeSubmitReport(t *testing.T) {
	bytedata := []byte{
		0x01, 0x00, 0x11, 0x30, 0x22, 0x41, 0x52, 0x04, 0x63}
	p, e := sms.DecodeAsMS(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeSubmitReport(t *testing.T) {
	p := &sms.SubmitReport{
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("", 9*60*60)).Local(),
		PID: nil,
		DCS: nil}

	b := p.Encode()
	t.Logf("% x", b)
}

func TestMarshalJSON_submitreport(t *testing.T) {
	p := &sms.SubmitReport{
		FCS: 0xC0,
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("", 9*60*60)).Local(),
		PID: nil,
		DCS: nil}
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
