package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestConvertDCS(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 256; i++ {
		orig := sms.UnmarshalDataCoding(byte(i))
		if orig == nil {
			continue
		}
		t.Logf("%s", orig)
		b := orig.Marshal()
		t.Logf("% x", b)
		ocom := sms.UnmarshalDataCoding(b)
		t.Logf("%s", ocom)
		if !orig.Equal(ocom) {
			t.Fatalf("mismatch orig=%s ocom=%s", orig, ocom)
		}
	}
}

func randDCS() (d sms.DataCoding) {
	for d == nil {
		d = sms.UnmarshalDataCoding(randByte())
	}
	return
}
