package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func TestConvertMORPDATA(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.RpData{
			MR: randByte(),
			DA: sms.Address{
				TON: sms.TypeInternational,
				NPI: sms.PlanISDNTelephone},
			UD: randSubmit()}
		tmp := randDigit((rand.Int() % 20) + 1)
		var e error
		orig.DA.Addr, e = teldata.ParseTBCD(tmp)
		if e != nil {
			panic(e)
		}

		t.Logf("%s", orig)
		b := orig.MarshalRPMO()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpData)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if !orig.DA.Equal(ocom.DA) {
			t.Fatal("DA mismatch")
		}
		if orig.UD.String() != ocom.UD.String() {
			t.Fatal("UD mismatch")
		}
	}
}

func TestConvertMTRPDATA(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.RpData{
			MR: randByte(),
			OA: sms.Address{
				TON: sms.TypeInternational,
				NPI: sms.PlanISDNTelephone},
			UD: randDeliver()}
		tmp := randDigit((rand.Int() % 20) + 1)
		var e error
		orig.OA.Addr, e = teldata.ParseTBCD(tmp)
		if e != nil {
			panic(e)
		}

		t.Logf("%s", orig)
		b := orig.MarshalRPMT()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpData)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if !orig.OA.Equal(ocom.OA) {
			t.Fatal("OA mismatch")
		}
		if orig.UD.String() != ocom.UD.String() {
			t.Fatal("UD mismatch")
		}
	}
}
