package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertMORPERROR(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.RpError{
			MR: randByte(),
			CS: randByte()}
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.Diag = &bt
		}
		if randByte() != 0 {
			orig.UD = randDeliverreport()
		}

		t.Logf("%s", orig)
		b := orig.MarshalRPMO()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpError)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
		if (orig.UD == nil) != (ocom.UD == nil) {
			t.Fatal("UD mismatch")
		}
		if orig.UD != nil && ocom.UD != nil && orig.UD.String() != ocom.UD.String() {
			t.Fatal("UD mismatch")
		}
	}
}

func TestConvertMTRPERROR(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.RpError{
			MR: randByte(),
			CS: randByte()}
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.Diag = &bt
		}
		if randByte() != 0 {
			orig.UD = randSubmitreport()
		}

		t.Logf("%s", orig)
		b := orig.MarshalRPMT()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpError)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
		if (orig.UD == nil) != (ocom.UD == nil) {
			t.Fatal("UD mismatch")
		}
		if orig.UD != nil && ocom.UD != nil && orig.UD.String() != ocom.UD.String() {
			t.Fatal("UD mismatch")
		}
	}
}
