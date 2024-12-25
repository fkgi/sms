package sms_test

import (
	"testing"

	"github.com/fkgi/sms"
)

func TestConvertCPErrorMO(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := sms.CpError{
			TI: randTransactionID(),
			CS: randByte()}

		t.Logf("%s", orig)
		b := orig.MarshalCP()
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

func TestConvertCPErrorMT(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := sms.CpError{
			TI: randTransactionID(),
			CS: randByte()}

		t.Logf("%s", orig)
		b := orig.MarshalCP()
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
