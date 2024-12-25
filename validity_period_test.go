package sms_test

import (
	"log"
	"testing"
	"time"

	"github.com/fkgi/sms"
)

func randDuration() (time.Duration, bool) {
	// return time.Duration(rand.Int63n(3122064000)) * time.Second
	vp := randVP()
	for vp == nil {
		vp = randVP()
	}
	log.Println(vp)
	return vp.Duration(), vp.SingleAttempt()
}

func TestMakeVPfmDuration(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig, sa := randDuration()
		vp := sms.ValidityPeriodOf(orig, sa)
		t.Log(vp)

		ocom := vp.Duration()
		if orig+time.Second*2 < ocom || orig-time.Second*2 > ocom {
			t.Fatalf("duration missmatch, orig=%s, ocom=%s", orig, ocom)
		}
		if sa != vp.SingleAttempt() {
			t.Fatalf("single attempt missmatch, orig=%t, conv=%t",
				sa, vp.SingleAttempt())
		}
	}
}

func TestDurationValue(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig, sa := randDuration()
		if orig == 0 {
			continue
		}
		vp := sms.ValidityPeriodOf(orig, sa)
		t.Log(vp)

		now := time.Now()
		exp1 := now.Add(orig)
		exp2 := vp.ExpireTime(now)
		if exp1.Add(-time.Second*2).After(exp2) ||
			exp1.Add(time.Second*2).Before(exp2) {
			t.Fatalf("expire missmatch, exp1=%s, exp2=%s", exp1, exp2)
		}
	}
}
