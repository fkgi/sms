package sms_test

import (
	"errors"
	"testing"

	"github.com/fkgi/sms"
)

func randRPAckMO() sms.RpAckMO {
	orig := sms.RpAckMO{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	return orig
}

func compareRPRPAckMO(orig, ocom sms.RpAckMO) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	return nil
}

func compareCPRPAckMO(orig, ocom sms.RpAckMO) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPAckMO(orig, ocom)
}

func TestConvertRPRPAckMO(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randRPAckMO()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpAckMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPRPAckMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalRpAckMO(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPRPAckMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPRPAckMO(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randRPAckMO()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpAckMO)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPRPAckMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.RpAckMO{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPRPAckMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func randRPAckMT() sms.RpAckMT {
	orig := sms.RpAckMT{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	return orig
}

func compareRPRPAckMT(orig, ocom sms.RpAckMT) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	return nil
}

func compareCPRPAckMT(orig, ocom sms.RpAckMT) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPAckMT(orig, ocom)
}

func TestConvertRPRPAckMT(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randRPAckMT()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpAckMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPRPAckMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalRpAckMT(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPRPAckMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPRPAckMT(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randRPAckMT()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMT(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.RpAckMT)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPRPAckMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.RpAckMT{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPRPAckMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
