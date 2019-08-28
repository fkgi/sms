package sms_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func TestVP(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		orig := time.Duration(rand.Int63n(100000)) * time.Second
		vp := sms.ValidityPeriodOf(orig, randBool())
		t.Log(vp)
		ocom := vp.Duration()
		if orig != ocom {
			t.Fatalf("duration missmatch, orig=%ds, ocom=%ds",
				orig/time.Second, ocom/time.Second)
		}

		now := time.Now()
		exp1 := now.Add(orig)
		exp2 := vp.ExpireTime(now)
		if !exp1.Equal(exp2) {
			t.Fatalf("expire missmatch, exp1=%s, exp2=%s",
				exp1, exp2)
		}
	}
}
