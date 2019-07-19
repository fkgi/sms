package sms_test

import (
	"testing"

	"github.com/fkgi/sms"
)

func MakeRPDATA(t *testing.T) {
	d := sms.Data{
		MR: 0x01,
	}
	//	if e != nil {
	//		t.Fatalf("encode failed: %s", e)
	//	}
	t.Log(d.String())
}
