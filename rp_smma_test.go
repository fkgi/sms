package sms_test

import (
	"errors"
	"testing"

	"github.com/fkgi/sms"
)

func randMemoryAvailable() sms.MemoryAvailable {
	orig := sms.MemoryAvailable{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	return orig
}

func compareRPMemoryAvailable(orig, ocom sms.MemoryAvailable) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	return nil
}

func compareCPMemoryAvailable(orig, ocom sms.MemoryAvailable) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPMemoryAvailable(orig, ocom)
}

func TestConvertRPMemoryAvailable(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randMemoryAvailable()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.MemoryAvailable)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPMemoryAvailable(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom, e = sms.UnmarshalMemoryAvailable(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPMemoryAvailable(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPMemoryAvailable(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randMemoryAvailable()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.MemoryAvailable)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPMemoryAvailable(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.MemoryAvailable{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPMemoryAvailable(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
