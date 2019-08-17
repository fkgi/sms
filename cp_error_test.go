package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertMOCPERROR(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.CpError{
			TI: byte(rand.Int31n(16)),
			CS: randByte()}

		t.Logf("%s", orig)
		b := orig.MarshalCPMO()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.CpError)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.TI != ocom.TI {
			t.Fatal("TI mismatch")
		}
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
	}
}

func TestConvertMTCPERROR(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.CpError{
			TI: byte(rand.Int31n(16)),
			CS: randByte()}

		t.Logf("%s", orig)
		b := orig.MarshalCPMT()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.CpError)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.TI != ocom.TI {
			t.Fatal("TI mismatch")
		}
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
	}
}
