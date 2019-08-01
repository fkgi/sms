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
		orig := sms.Ack{
			MR: randByte(),
			UD: randDeliverreport()}

		t.Logf("%s", orig)
		b := orig.MarshalRPMO()
		t.Logf("% x", b)
		res, e := sms.UnmarshalMORP(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Ack)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
	}
}

func TestConvertMTRPACK(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.Ack{
			MR: randByte(),
			UD: randSubmitreport()}

		t.Logf("%s", orig)
		b := orig.MarshalRPMT()
		t.Logf("% x", b)
		res, e := sms.UnmarshalMTRP(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Ack)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
	}
}
