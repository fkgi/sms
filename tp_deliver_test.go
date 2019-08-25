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

func TestEncodeDeliver(t *testing.T) {
	bytedata := []byte{
		0x40, 0x04, 0x80, 0x21, 0x43, 0x00, 0x08, 0x11,
		0x30, 0x22, 0x41, 0x52, 0x04, 0x63, 0x10, 0x05,
		0x00, 0x03, 0x87, 0x02, 0x01, 0x30, 0x42, 0x30,
		0x44, 0x30, 0x46, 0x30, 0x48, 0x30, 0x4a}
	p, e := sms.UnmarshalTPMT(bytedata)
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
			MsgCharset: sms.CharsetUCS2},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		UD: sms.UserData{Text: "あいうえお"}}
	p.UD.UDH = append(p.UD.UDH, sms.ConcatenatedSM{
		RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
	p.OA.Addr, _ = teldata.ParseTBCD("1234")

	b := p.MarshalTP()
	t.Logf("% x", b)
}

func randDeliver() sms.Deliver {
	orig := sms.Deliver{
		MMS:  randBool(),
		LP:   randBool(),
		SRI:  randBool(),
		RP:   randBool(),
		OA:   randAddress(),
		PID:  randByte(),
		DCS:  randDCS(),
		SCTS: randDate(),
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

func compareTPDeliver(orig, ocom sms.Deliver) error {
	if orig.MMS != ocom.MMS {
		return errors.New("MMS mismatch")
	}
	if orig.LP != ocom.LP {
		return errors.New("LP mismatch")
	}
	if orig.SRI != ocom.SRI {
		return errors.New("SRI mismatch")
	}
	if orig.RP != ocom.RP {
		return errors.New("RP mismatch")
	}
	if !orig.OA.Equal(ocom.OA) {
		return errors.New("OA mismatch")
	}
	if orig.PID != ocom.PID {
		return errors.New("PID mismatch")
	}
	if !orig.DCS.Equal(ocom.DCS) {
		return errors.New("DCS mismatch")
	}
	if !orig.SCTS.Equal(ocom.SCTS) {
		return errors.New("SCTS mismatch")
	}
	if !orig.UD.Equal(ocom.UD) {
		return errors.New("UD text mismatch")
	}
	return nil
}

func compareRPDeliver(orig, ocom sms.Deliver) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if !orig.SCA.Equal(ocom.SCA) {
		return errors.New("SCA mismatch")
	}
	return compareTPDeliver(orig, ocom)
}

func compareCPDeliver(orig, ocom sms.Deliver) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPDeliver(orig, ocom)
}

func TestConvertTPDeliver(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randDeliver()

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalTPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Deliver)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareTPDeliver(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Deliver{}
		e = ocom.UnmarshalTP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareTPDeliver(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertRPDeliver(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randDeliver()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Deliver)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPDeliver(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Deliver{}
		e = ocom.UnmarshalRP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPDeliver(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPDeliver(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randDeliver()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Deliver)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPDeliver(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Deliver{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPDeliver(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
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
			MsgCharset: sms.CharsetUCS2},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		UD: sms.UserData{Text: "あいうえお"}}
	p.UD.UDH = append(p.UD.UDH, sms.ConcatenatedSM{
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
	p, e := sms.UnmarshalTPMO(bytedata)
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
			MsgCharset: sms.CharsetUCS2},
		UD: sms.UserData{Text: "あいうえお"}}
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

func randDeliverreport() sms.DeliverReport {
	orig := sms.DeliverReport{
		FCS: byte(rand.Int31n(129)),
		DCS: sms.UnmarshalDataCoding(randByte()),
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

func compareTPDeliverReport(orig, ocom sms.DeliverReport) error {
	if orig.FCS != ocom.FCS {
		return errors.New("FCS mismatch")
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

func compareRPDeliverReport(orig, ocom sms.DeliverReport) error {
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
	return compareTPDeliverReport(orig, ocom)
}

func compareCPDeliverReport(orig, ocom sms.DeliverReport) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPDeliverReport(orig, ocom)
}

func TestConvertTPDeliverreport(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randDeliverreport()

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalTPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.DeliverReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareTPDeliverReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.DeliverReport{}
		e = ocom.UnmarshalTP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareTPDeliverReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertRPDeliverreport(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randDeliverreport()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.DeliverReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPDeliverReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.DeliverReport{}
		e = ocom.UnmarshalRP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPDeliverReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPDeliverreport(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randDeliverreport()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.DeliverReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPDeliverReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.DeliverReport{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPDeliverReport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
