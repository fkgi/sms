package sms_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func randRPErrorMO() sms.ErrorMO {
	orig := sms.ErrorMO{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	orig.CS = randByte()
	if tmp := rand.Int31n(257); tmp != 256 {
		bt := byte(tmp)
		orig.DIAG = &bt
	}
	return orig
}

func compareRPRPErrorMO(orig, ocom sms.ErrorMO) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if orig.CS != ocom.CS {
		return errors.New("CS mismatch")
	}
	if (orig.DIAG == nil) != (ocom.DIAG == nil) {
		return errors.New("DIAG mismatch")
	}
	if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
		return errors.New("DIAG mismatch")
	}
	return nil
}

func compareCPRPErrorMO(orig, ocom sms.ErrorMO) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPErrorMO(orig, ocom)
}

func TestConvertRPRPErrorMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPErrorMO()

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

		e = compareRPRPErrorMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalErrorMO(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPRPErrorMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPRPErrorMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPErrorMO()

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

		e = compareCPRPErrorMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.ErrorMO{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPRPErrorMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func randRPErrorMT() sms.ErrorMT {
	orig := sms.ErrorMT{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	orig.CS = randByte()
	if tmp := rand.Int31n(257); tmp != 256 {
		bt := byte(tmp)
		orig.DIAG = &bt
	}
	return orig
}

func compareRPRPErrorMT(orig, ocom sms.ErrorMT) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if orig.CS != ocom.CS {
		return errors.New("CS mismatch")
	}
	if (orig.DIAG == nil) != (ocom.DIAG == nil) {
		return errors.New("DIAG mismatch")
	}
	if orig.DIAG != nil && ocom.DIAG != nil && *orig.DIAG != *ocom.DIAG {
		return errors.New("DIAG mismatch")
	}
	return nil
}

func compareCPRPErrorMT(orig, ocom sms.ErrorMT) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPErrorMT(orig, ocom)
}

func TestConvertRPRPErrorMT(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPErrorMT()

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

		e = compareRPRPErrorMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalErrorMT(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPRPErrorMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPRPErrorMR(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPErrorMT()

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

		e = compareCPRPErrorMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.ErrorMT{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPRPErrorMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
