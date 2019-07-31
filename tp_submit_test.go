package sms_test

import (
	"encoding/json"
	"math/rand"
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
	p, e := sms.UnmarshalMOTP(bytedata)
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
			MsgCharset: sms.CharsetUCS2},
		VP: nil,
		UD: sms.UD{Text: "あいうえお"}}
	p.UD.UDH = append(p.UD.UDH, sms.ConcatenatedSM{
		RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
	p.DA.Addr, _ = teldata.ParseTBCD("09012345678")

	b := p.MarshalTP()
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
			MsgCharset: sms.CharsetUCS2},
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

func randVP() sms.VP {
	switch rand.Int31n(4) {
	case 1:
		return sms.VPRelative(randByte())
	case 2:
		t := randDate()
		var r [7]byte
		r[0] = int2SemiOctet(t.Year())
		r[1] = int2SemiOctet(int(t.Month()))
		r[2] = int2SemiOctet(t.Day())
		r[3] = int2SemiOctet(t.Hour())
		r[4] = int2SemiOctet(t.Minute())
		r[5] = int2SemiOctet(t.Second())

		_, z := t.Zone()
		z /= 900
		if z < 0 {
			z = -z
			r[6] = byte(z % 10)
			r[6] = (r[6] << 4) | byte(((z/10)%10)&0x07)
			r[6] = r[6] | 0x08
		} else {
			r[6] = byte(z % 10)
			r[6] = (r[6] << 4) | byte(((z/10)%10)&0x07)
		}
		return sms.VPAbsolute(r)
	case 3:
		return sms.VPEnhanced{
			randByte(), randByte(), randByte(), randByte(), randByte(), randByte(), randByte()}
	}
	return nil
}
func int2SemiOctet(i int) (b byte) {
	b = byte(i % 10)
	b = (b << 4) | byte((i/10)%10)
	return
}

func TestConvertSubmit(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.Submit{
			RD:  randBool(),
			SRR: randBool(),
			RP:  randBool(),
			MR:  randByte(),
			DA:  genRandomAddress(),
			PID: randByte(),
			DCS: getRandomDCS(),
			VP:  randVP(),
		}

		orig.UD = getRandomUD(orig.DCS)

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		ocom, e := sms.UnmarshalSubmit(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		if orig.RD != ocom.RD {
			t.Fatal("RD mismatch")
		}
		if orig.SRR != ocom.SRR {
			t.Fatal("SRR mismatch")
		}
		if orig.RP != ocom.RP {
			t.Fatal("RP mismatch")
		}
		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if !orig.DA.Equal(ocom.DA) {
			t.Fatal("DA mismatch")
		}
		if orig.PID != ocom.PID {
			t.Fatal("PID mismatch")
		}
		if !orig.DCS.Equal(ocom.DCS) {
			t.Fatal("DCS mismatch")
		}
		if orig.VP == nil && ocom.VP != nil {
			t.Fatal("VP mismatch")
		}
		if orig.VP != nil && ocom.VP == nil {
			t.Fatal("VP mismatch")
		}
		if orig.VP != nil && ocom.VP != nil && !orig.VP.Equal(ocom.VP) {
			t.Fatal("VP mismatch")
		}
		if !orig.UD.Equal(ocom.UD) {
			t.Fatal("UD text mismatch")
		}
	}
}

func TestEncodeSubmitReport(t *testing.T) {
	bytedata := []byte{
		0x01, 0x00, 0x11, 0x30, 0x22, 0x41, 0x52, 0x04, 0x63}
	p, e := sms.UnmarshalMTTP(bytedata)
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

	b := p.MarshalTP()
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

func TestConvertSubmitreport(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.SubmitReport{
			FCS:  byte(rand.Int31n(129)),
			SCTS: randDate(),
			DCS:  sms.UnmarshalDCS(randByte()),
		}

		if orig.FCS == 128 {
			orig.FCS = 0
		} else {
			orig.FCS += 128
		}
		if tmp := rand.Int31n(257); tmp != 256 {
			b := byte(tmp)
			orig.PID = &b
		}
		if orig.DCS != nil {
			orig.UD = getRandomUD(orig.DCS)
		}

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		ocom, e := sms.UnmarshalSubmitReport(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		if orig.FCS != ocom.FCS {
			t.Fatal("FCS mismatch")
		}
		if !orig.SCTS.Equal(ocom.SCTS) {
			t.Fatal("SCTS text mismatch")
		}
		if orig.PID == nil && ocom.PID != nil {
			t.Fatal("PID mismatch")
		}
		if orig.PID != nil && ocom.PID == nil {
			t.Fatal("PID mismatch")
		}
		if orig.PID != nil && ocom.PID != nil && *orig.PID != *ocom.PID {
			t.Fatal("PID mismatch")
		}
		if orig.DCS == nil && ocom.DCS != nil {
			t.Fatal("DCS mismatch")
		}
		if orig.DCS != nil && ocom.DCS == nil {
			t.Fatal("DCS mismatch")
		}
		if orig.DCS != nil && ocom.DCS != nil && !orig.DCS.Equal(ocom.DCS) {
			t.Fatal("DCS mismatch")
		}
		if !orig.UD.Equal(ocom.UD) {
			t.Fatal("UD text mismatch")
		}
	}
}
