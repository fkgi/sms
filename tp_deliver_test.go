package sms_test

import (
	"encoding/json"
	"math/rand"
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
	p := sms.Deliver{
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

func randBool() bool {
	return rand.Int31n(2) == 0
}

func randByte() byte {
	return byte(rand.Int31n(256))
}

func TestConvertDeliver(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var e error

	orig := sms.Deliver{
		MMS: randBool(),
		LP:  randBool(),
		SRI: randBool(),
		RP:  randBool(),
		// OA:  ,
		PID: randByte(),
		DCS: getRandomDCS(),
		SCTS: time.Date(
			2000+rand.Int()%100, time.Month(rand.Int()%12+1),
			rand.Int()%32, rand.Int()%24, rand.Int()%60, rand.Int()%60, 0,
			time.FixedZone("", 9*60*60)).Local(),
		UD: sms.UD{}}

	orig.OA, e = genRandomAddress()
	if e != nil {
		t.Fatal(e)
	}

	for i := rand.Int() % 5; i >= 0; i-- {
		h := randByte()
		switch h {
		case 0x00:
			orig.UD.AddUDH(&sms.ConcatenatedSM{
				RefNum: randByte(),
				MaxNum: randByte(),
				SeqNum: randByte()})
		case 0x08:
			orig.UD.AddUDH(&sms.ConcatenatedSM16bit{
				RefNum: uint16(rand.Int() % 65535),
				MaxNum: randByte(),
				SeqNum: randByte()})
		default:
			iei := &sms.GenericIEI{
				K: h,
				V: make([]byte, rand.Int()%5+1)}
			for i := range iei.V {
				iei.V[i] = randByte()
			}
			orig.UD.AddUDH(iei)
		}
	}
	/*
		switch orig.DCS.charset() {
		case sms.CharsetGSM7bit:
			orig.UD.Text = randText(rand.Int() % 100)
		case sms.Charset8bitData:
			tmp := make([]byte, rand.Int()%100)
			for i := range tmp {
				tmp[i] = randByte()
			}
			orig.UD.Set8bitData(tmp)
		case sms.CharsetUCS2:
			orig.UD.Text = "aiueo"
		}
	*/
	t.Logf("%s", orig)
	b := orig.MarshalTP()
	t.Logf("% x", b)
	ocom, e := sms.UnmarshalDeliver(b)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%s", ocom)

	if orig.MMS != ocom.MMS {
		t.Fatal("MMS mismatch")
	}
	if orig.LP != ocom.LP {
		t.Fatal("LP mismatch")
	}
	if orig.SRI != ocom.SRI {
		t.Fatal("SRI mismatch")
	}
	if orig.RP != ocom.RP {
		t.Fatal("RP mismatch")
	}
	if !orig.OA.Equal(ocom.OA) {
		t.Fatal("OA mismatch")
	}
	if orig.PID != ocom.PID {
		t.Fatal("PID mismatch")
	}
	if !orig.DCS.Equal(ocom.DCS) {
		t.Fatal("DCS mismatch")
	}
	if !orig.SCTS.Equal(ocom.SCTS) {
		t.Fatal("SCTS mismatch")
	}
	if !orig.UD.Equal(ocom.UD) {
		t.Fatal("UD text mismatch")
	}
}

func TestMarshalJSON_deliver(t *testing.T) {
	p := sms.Deliver{
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

	bytedata, e := json.Marshal(p)
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
	p := sms.DeliverReport{
		PID: nil,
		DCS: nil}

	b := p.MarshalTP()
	t.Logf("% x", b)
}

func TestMarshalJSON_deliverreport(t *testing.T) {
	p := sms.DeliverReport{
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

	bytedata, e := json.Marshal(p)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))

	if e := json.Unmarshal(bytedata, &p); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(p.String())
}
