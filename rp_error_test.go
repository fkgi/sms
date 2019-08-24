package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertRPRPErrorMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.ErrorMO{}
		orig.RMR = randByte()
		orig.CS = randByte()
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.DIAG = &bt
		}

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.ErrorMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.RMR != ocom.RMR {
			t.Fatal("MR mismatch")
		}
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
		if (orig.DIAG == nil) != (ocom.DIAG == nil) {
			t.Fatal("DIAG mismatch")
		}
		if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
			t.Fatal("DIAG mismatch")
		}
	}
}

func TestConvertCPRPErrorMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.ErrorMO{}
		orig.TI = randTransactionID()
		orig.RMR = randByte()
		orig.CS = randByte()
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.DIAG = &bt
		}

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.ErrorMO)
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
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
		if (orig.DIAG == nil) != (ocom.DIAG == nil) {
			t.Fatal("DIAG mismatch")
		}
		if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
			t.Fatal("DIAG mismatch")
		}
	}
}

func TestConvertRPRPErrorMT(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.ErrorMT{}
		orig.RMR = randByte()
		orig.CS = randByte()
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.DIAG = &bt
		}

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.ErrorMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.RMR != ocom.RMR {
			t.Fatal("MR mismatch")
		}
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
		if (orig.DIAG == nil) != (ocom.DIAG == nil) {
			t.Fatal("DIAG mismatch")
		}
		if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
			t.Fatal("DIAG mismatch")
		}
	}
}

func TestConvertCPRPErrorMR(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.ErrorMT{}
		orig.TI = randTransactionID()
		orig.RMR = randByte()
		orig.CS = randByte()
		if tmp := rand.Int31n(257); tmp != 256 {
			bt := byte(tmp)
			orig.DIAG = &bt
		}

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.ErrorMT)
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
		if orig.CS != ocom.CS {
			t.Fatal("CS mismatch")
		}
		if (orig.DIAG == nil) != (ocom.DIAG == nil) {
			t.Fatal("DIAG mismatch")
		}
		if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
			t.Fatal("DIAG mismatch")
		}
	}
}
