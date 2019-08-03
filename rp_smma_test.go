package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertMORPSMMA(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := sms.MemoryAvailable{
			MR: randByte()}

		t.Logf("%s", orig)
		b := orig.MarshalRPMO()
		t.Logf("% x", b)
		res, e := sms.UnmarshalMORP(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.MemoryAvailable)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		if orig.MR != ocom.MR {
			t.Fatal("MR mismatch")
		}
	}
}
