package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertRPRPAckMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.AckMO{}
		orig.RMR = randByte()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.AckMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.RMR != ocom.RMR {
			t.Fatal("MR mismatch")
		}
	}
}

func TestConvertCPRPAckMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.AckMO{}
		orig.TI = randTransactionID()
		orig.RMR = randByte()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.AckMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.TI != ocom.TI {
			t.Fatal("TI mismatch")
		}
		if orig.RMR != ocom.RMR {
			t.Fatal("MR mismatch")
		}
	}
}

func TestConvertRPRPAckMT(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.AckMT{}
		orig.RMR = randByte()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.AckMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.RMR != ocom.RMR {
			t.Fatal("MR mismatch")
		}
	}
}

func TestConvertCPRPAckMT(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.AckMT{}
		orig.TI = randTransactionID()
		orig.RMR = randByte()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.AckMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.TI != ocom.TI {
			t.Fatal("TI mismatch")
		}
		if orig.RMR != ocom.RMR {
			t.Fatal("MR mismatch")
		}
	}
}
