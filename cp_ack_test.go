package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertMOCPACK(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.CpAck{
			TI: byte(rand.Int31n(16))}

		t.Logf("%s", orig)
		b := orig.MarshalCPMO()
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

func TestConvertMTCPACK(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.CpAck{
			TI: byte(rand.Int31n(16))}

		t.Logf("%s", orig)
		b := orig.MarshalCPMT()
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
