package sms_test

import (
	"testing"

	"github.com/fkgi/sms"
)

func TestConvertCPAckMO(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := sms.CpAck{
			TI: randTransactionID()}

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.CpAck)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.TI != ocom.TI {
			t.Fatal("TI mismatch")
		}
	}
}

func TestConvertCPAckMT(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := sms.CpAck{
			TI: randTransactionID()}

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.CpAck)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.TI != ocom.TI {
			t.Fatal("TI mismatch")
		}
	}
}
