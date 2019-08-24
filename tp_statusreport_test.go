package sms_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

var bytedata = []byte{
	0x06, 0x00, 0x04, 0x80, 0x21, 0x43, 0x11, 0x30,
	0x22, 0x41, 0x52, 0x04, 0x63, 0x11, 0x30, 0x22,
	0x41, 0x52, 0x04, 0x63, 0x00}

func TestEncodeStatusReport(t *testing.T) {
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

func randStatusreport() sms.StatusReport {
	orig := sms.StatusReport{
		MMS:  randBool(),
		LP:   randBool(),
		SRQ:  randBool(),
		TMR:  randByte(),
		RA:   randAddress(),
		SCTS: randDate(),
		DT:   randDate(),
		DCS:  sms.UnmarshalDataCoding(randByte()),
	}
	if tmp := rand.Int31n(257); tmp != 256 {
		b := byte(tmp)
		orig.PID = &b
	}
	if orig.DCS != nil {
		orig.UD = randUD(orig.DCS)
	}
	return orig
}

func TestConvertStatusreport(t *testing.T) {
	rand.Seed(time.Now().Unix())

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

		if orig.MMS != ocom.MMS {
			t.Fatal("MMS mismatch")
		}
		if orig.LP != ocom.LP {
			t.Fatal("LP mismatch")
		}
		if orig.SRQ != ocom.SRQ {
			t.Fatal("SRQ mismatch")
		}
		if orig.TMR != ocom.TMR {
			t.Fatal("MR mismatch")
		}
		if !orig.RA.Equal(ocom.RA) {
			t.Fatal("RA mismatch")
		}
		if !orig.SCTS.Equal(ocom.SCTS) {
			t.Fatal("SCTS mismatch")
		}
		if !orig.DT.Equal(ocom.DT) {
			t.Fatal("DT mismatch")
		}
		if (orig.PID == nil) != (ocom.PID == nil) {
			t.Fatal("PID mismatch")
		}
		if orig.PID != nil && ocom.PID != nil && *orig.PID != *ocom.PID {
			t.Fatal("PID mismatch")
		}
		if (orig.DCS == nil) != (ocom.DCS == nil) {
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
