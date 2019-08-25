package sms_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func randRPAckMO() sms.AckMO {
	orig := sms.AckMO{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	return orig
}

func compareRPRPAckMO(orig, ocom sms.AckMO) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	return nil
}

func compareCPRPAckMO(orig, ocom sms.AckMO) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPAckMO(orig, ocom)
}

func TestConvertRPRPAckMO(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPAckMO()

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

		e = compareRPRPAckMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.AckMO{}
		e = ocom.UnmarshalRP(b)
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
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPAckMO()

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

		e = compareCPRPAckMO(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.AckMO{}
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

func randRPAckMT() sms.AckMT {
	orig := sms.AckMT{}
	orig.TI = randTransactionID()
	orig.RMR = randByte()
	return orig
}

func compareRPRPAckMT(orig, ocom sms.AckMT) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	return nil
}

func compareCPRPAckMT(orig, ocom sms.AckMT) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPRPAckMT(orig, ocom)
}

func TestConvertRPRPAckMT(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPAckMT()

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

		e = compareRPRPAckMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.AckMT{}
		e = ocom.UnmarshalRP(b)
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
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randRPAckMT()

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

		e = compareCPRPAckMT(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.AckMT{}
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
