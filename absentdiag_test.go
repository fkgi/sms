package sms_test

import (
	"testing"

	"github.com/fkgi/sms"
)

func TestAbsentDiag(t *testing.T) {
	t.Log(sms.NoAbsentDiag.String())
	var i byte

	for i = 0; i < 255; i++ {
		a := sms.B2AbsDiag(i)
		t.Log(a.String())
		b := a.Byte()
		if i != b {
			t.Fatal(i, "!=", b)
		}
	}
}
