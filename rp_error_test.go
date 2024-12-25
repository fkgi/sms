package sms_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/fkgi/sms"
)

func randRPErrorMO() sms.RpErrorMO {
	orig := sms.RpErrorMO{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	orig.CS = randByte()
	if tmp := rand.Int31n(257); tmp != 256 {
		bt := byte(tmp)
		orig.DIAG = &bt
	}
	return orig
}

func compareRPRPErrorMO(orig, ocom sms.RpErrorMO) error {
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

func compareCPRPErrorMO(orig, ocom sms.RpErrorMO) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPErrorMO(orig, ocom)
}

func TestConvertRPRPErrorMO(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randRPErrorMO()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpErrorMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPRPErrorMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalRpErrorMO(b)
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
	for i := 0; i < 1000; i++ {
		orig := randRPErrorMO()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpErrorMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPRPErrorMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.RpErrorMO{}
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

func randRPErrorMT() sms.RpErrorMT {
	orig := sms.RpErrorMT{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	orig.CS = randByte()
	if tmp := rand.Int31n(257); tmp != 256 {
		bt := byte(tmp)
		orig.DIAG = &bt
	}
	return orig
}

func compareRPRPErrorMT(orig, ocom sms.RpErrorMT) error {
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

func compareCPRPErrorMT(orig, ocom sms.RpErrorMT) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPErrorMT(orig, ocom)
}

func TestConvertRPRPErrorMT(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randRPErrorMT()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpErrorMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPRPErrorMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalRpErrorMT(b)
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
	for i := 0; i < 1000; i++ {
		orig := randRPErrorMT()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpErrorMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPRPErrorMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.RpErrorMT{}
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
