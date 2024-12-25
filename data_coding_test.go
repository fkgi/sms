package sms_test

import (
	"testing"

	"github.com/fkgi/sms"
)

func TestConvertDCS(t *testing.T) {
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
