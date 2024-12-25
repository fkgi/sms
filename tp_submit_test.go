package sms_test

import (
	"encoding/json"
	"errors"
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
	p, e := sms.UnmarshalTPMO(bytedata)
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
		TMR: 64,
		DA:  sms.Address{TON: 0, NPI: 1},
		PID: 0,
		DCS: &sms.GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   sms.NoMessageClass,
			MsgCharset: sms.CharsetUCS2},
		VP: nil,
		UD: sms.UserData{Text: "あいうえお"}}
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
		TMR: 64,
		DA:  sms.Address{TON: 0, NPI: 1},
		PID: 0,
		DCS: &sms.GeneralDataCoding{
			AutoDelete: false,
			Compressed: false,
			MsgClass:   sms.NoMessageClass,
			MsgCharset: sms.CharsetUCS2},
		VP: nil,
		UD: sms.UserData{Text: "あいうえお"}}
	p.DA.Addr, _ = teldata.ParseTBCD("09012345678")

	t.Log(p.String())

	var e error
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

func randSubmit() sms.Submit {
	orig := sms.Submit{
		RD:  randBool(),
		SRR: randBool(),
		RP:  randBool(),
		TMR: randByte(),
		DA:  randAddress(),
		PID: randByte(),
		DCS: randDCS(),
		VP:  randVP(),
	}

	orig.UD = randUD(orig.DCS)

	orig.TI = randTransactionID()
	orig.RMR = randByte()
	orig.SCA = sms.Address{
		TON: sms.TypeInternational,
		NPI: sms.PlanISDNTelephone}
	tmp := randDigit((rand.Int() % 20) + 1)
	var e error
	orig.SCA.Addr, e = teldata.ParseTBCD(tmp)
	if e != nil {
		panic(e)
	}
	return orig
}

func compareTPSubmit(orig, ocom sms.Submit) error {
	if orig.RD != ocom.RD {
		return errors.New("RD mismatch")
	}
	if orig.SRR != ocom.SRR {
		return errors.New("SRR mismatch")
	}
	if orig.RP != ocom.RP {
		return errors.New("RP mismatch")
	}
	if orig.TMR != ocom.TMR {
		return errors.New("MR mismatch")
	}
	if !orig.DA.Equal(ocom.DA) {
		return errors.New("DA mismatch")
	}
	if orig.PID != ocom.PID {
		return errors.New("PID mismatch")
	}
	if !orig.DCS.Equal(ocom.DCS) {
		return errors.New("DCS mismatch")
	}
	if (orig.VP == nil) != (ocom.VP == nil) {
		return errors.New("VP mismatch")
	}
	if orig.VP != nil && ocom.VP != nil && !orig.VP.Equal(ocom.VP) {
		return errors.New("VP mismatch")
	}
	if !orig.UD.Equal(ocom.UD) {
		return errors.New("UD text mismatch")
	}
	return nil
}

func compareRPSubmit(orig, ocom sms.Submit) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if !orig.SCA.Equal(ocom.SCA) {
		return errors.New("SCA mismatch")
	}
	return compareTPSubmit(orig, ocom)
}

func compareCPSubmit(orig, ocom sms.Submit) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPSubmit(orig, ocom)
}

func TestConvertTPSubmit(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randSubmit()

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalTPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Submit)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareTPSubmit(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Submit{}
		e = ocom.UnmarshalTP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareTPSubmit(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertRPSubmit(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randSubmit()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Submit)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPSubmit(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Submit{}
		e = ocom.UnmarshalRP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPSubmit(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPSubmit(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randSubmit()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Submit)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPSubmit(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Submit{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPSubmit(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestEncodeSubmitReport(t *testing.T) {
	bytedata := []byte{
		0x01, 0x00, 0x11, 0x30, 0x22, 0x41, 0x52, 0x04, 0x63}
	p, e := sms.UnmarshalTPMT(bytedata)
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

func randSubmitreport() sms.SubmitReport {
	orig := sms.SubmitReport{
		FCS:  byte(rand.Int31n(129)),
		SCTS: randDate(),
		DCS:  sms.UnmarshalDataCoding(randByte()),
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
		orig.UD = randUD(orig.DCS)
	}

	orig.TI = randTransactionID()
	orig.RMR = randByte()
	if orig.FCS != 0 {
		orig.CS = randByte()
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.DIAG = &bt
		}
	}
	return orig
}

func compareTPSubmitReport(orig, ocom sms.SubmitReport) error {
	if orig.FCS != ocom.FCS {
		return errors.New("FCS mismatch")
	}
	if !orig.SCTS.Equal(ocom.SCTS) {
		return errors.New("SCTS text mismatch")
	}
	if (orig.PID == nil) != (ocom.PID == nil) {
		return errors.New("PID mismatch")
	}
	if orig.PID != nil && ocom.PID != nil && *orig.PID != *ocom.PID {
		return errors.New("PID mismatch")
	}
	if (orig.DCS == nil) != (ocom.DCS == nil) {
		return errors.New("DCS mismatch")
	}
	if orig.DCS != nil && ocom.DCS != nil && !orig.DCS.Equal(ocom.DCS) {
		return errors.New("DCS mismatch")
	}
	if !orig.UD.Equal(ocom.UD) {
		return errors.New("UD text mismatch")
	}
	return nil
}

func compareRPSubmitReport(orig, ocom sms.SubmitReport) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if orig.CS != ocom.CS {
		return errors.New("CS mismatch")
	}
	if (orig.DIAG == nil) != (ocom.DIAG == nil) {
		return errors.New("DIAG mismatch")
	}
	if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
		return errors.New("DIAG mismatch")
	}
	return compareTPSubmitReport(orig, ocom)
}

func compareCPSubmitReport(orig, ocom sms.SubmitReport) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPSubmitReport(orig, ocom)
}

func TestConvertTPSubmitreport(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randSubmitreport()

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalTPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.SubmitReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareTPSubmitReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.SubmitReport{}
		e = ocom.UnmarshalTP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareTPSubmitReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertRPSubmitreport(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randSubmitreport()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.SubmitReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPSubmitReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.SubmitReport{}
		e = ocom.UnmarshalRP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPSubmitReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPSubmitreport(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randSubmitreport()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.SubmitReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPSubmitReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.SubmitReport{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPSubmitReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
