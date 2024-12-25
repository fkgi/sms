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

func TestEncodeStatusReport(t *testing.T) {
	bytedata := []byte{
		0x06, 0x00, 0x04, 0x80, 0x21, 0x43, 0x11, 0x30,
		0x22, 0x41, 0x52, 0x04, 0x63, 0x11, 0x30, 0x22,
		0x41, 0x52, 0x04, 0x63, 0x00}
	p, e := sms.UnmarshalTPMT(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	t.Log(p.String())
}

func TestDecodeStatusReport(t *testing.T) {
	p := &sms.StatusReport{
		MMS: false,
		LP:  false,
		SRQ: false,
		TMR: 0x00,
		RA:  sms.Address{TON: 0, NPI: 0},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		DT: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		ST: 0x00}
	p.RA.Addr, _ = teldata.ParseTBCD("1234")

	b := p.MarshalTP()
	t.Logf("% x", b)
}

func TestMarshalJSON_statusreport(t *testing.T) {
	p := &sms.StatusReport{
		MMS: false,
		LP:  false,
		SRQ: false,
		TMR: 0x00,
		RA:  sms.Address{TON: 0, NPI: 0},
		SCTS: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		DT: time.Date(
			2011, time.March, 22, 14, 25, 40, 0,
			time.FixedZone("unknown", 9*60*60)),
		ST: 0x00}
	p.RA.Addr, _ = teldata.ParseTBCD("1234")
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

func randStatusreport() sms.StatusReport {
	orig := sms.StatusReport{
		MMS:  randBool(),
		LP:   randBool(),
		SRQ:  randBool(),
		TMR:  randByte(),
		RA:   randAddress(),
		SCTS: randDate(),
		DT:   randDate(),
		ST:   randByte(),
		DCS:  sms.UnmarshalDataCoding(randByte()),
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
	orig.SCA = sms.Address{
		TON: sms.TypeInternational,
		NPI: sms.PlanISDNTelephone}
	tmp := randDigit((rand.Int() % 20) + 1)
	var e error
	if orig.SCA.Addr, e = teldata.ParseTBCD(tmp); e != nil {
		panic(e)
	}
	return orig
}

func compareTPStatusreport(orig, ocom sms.StatusReport) error {
	if orig.MMS != ocom.MMS {
		return errors.New("MMS mismatch")
	}
	if orig.LP != ocom.LP {
		return errors.New("LP mismatch")
	}
	if orig.SRQ != ocom.SRQ {
		return errors.New("SRQ mismatch")
	}
	if orig.TMR != ocom.TMR {
		return errors.New("MR mismatch")
	}
	if !orig.RA.Equal(ocom.RA) {
		return errors.New("RA mismatch")
	}
	if !orig.SCTS.Equal(ocom.SCTS) {
		return errors.New("SCTS mismatch")
	}
	if !orig.DT.Equal(ocom.DT) {
		return errors.New("DT mismatch")
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

func compareRPStatusreport(orig, ocom sms.StatusReport) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if !orig.SCA.Equal(ocom.SCA) {
		return errors.New("SCA mismatch")
	}
	return compareTPStatusreport(orig, ocom)
}

func compareCPStatusreport(orig, ocom sms.StatusReport) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPStatusreport(orig, ocom)
}

func TestConvertTPStatusreport(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randStatusreport()

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalTPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.StatusReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareTPStatusreport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.StatusReport{}
		e = ocom.UnmarshalTP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareTPStatusreport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertRPStatusreport(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randStatusreport()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.StatusReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPStatusreport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.StatusReport{}
		e = ocom.UnmarshalRP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPStatusreport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPStatusreport(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randStatusreport()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.StatusReport)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPStatusreport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.StatusReport{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPStatusreport(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
