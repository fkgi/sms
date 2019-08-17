package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertMORPACK(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.RpAck{
			MR: randByte()}
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
		ocom, ok := res.(sms.RpAck)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if (orig.UD == nil) != (ocom.UD == nil) {
			t.Fatal("UD mismatch")
		}
		if orig.UD != nil && ocom.UD != nil && orig.UD.String() != ocom.UD.String() {
			t.Fatal("UD mismatch")
		}
	}
}

func TestConvertMTRPACK(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.RpAck{
			MR: randByte()}
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
		ocom, ok := res.(sms.RpAck)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
		if (orig.UD == nil) != (ocom.UD == nil) {
			t.Fatal("UD mismatch")
		}
		if orig.UD != nil && ocom.UD != nil && orig.UD.String() != ocom.UD.String() {
			t.Fatal("UD mismatch")
		}
	}
}
