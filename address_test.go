package sms_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func TestUnmarshalJSON(t *testing.T) {
	bytedata := []byte(
		`{"ton":0,"npi":0,"addr":"1"}`,
	)

	a := sms.Address{}
	if e := json.Unmarshal(bytedata, &a); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(a.String())

	var e error
	bytedata, e = json.Marshal(a)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))
}

func TestConvertAddress(t *testing.T) {
	for i := 0; i < 1000; i++ {
		orig := randAddress()
		t.Logf("%s", orig)
		l, b := orig.Marshal()
		t.Logf("\nlen=%d\ndata=% x", l, b)
		ocom := sms.UnmarshalAddress(l, b)
		t.Logf("%s", ocom)
		if !orig.Equal(ocom) {
			t.Fatalf("mismatch orig=%s ocom=%s", orig, ocom)
		}
	}
}

func randAddress() (a sms.Address) {
	a.TON = byte(rand.Int() % 7)
	if a.TON == sms.TypeAlphanumeric {
		a.NPI = sms.PlanUnknown
		tmp := randText(rand.Int() % 11)
		var e error
		a.Addr, e = sms.StringToGSM7bit(tmp)
		if e != nil {
			panic(e)
		}
	} else {
		a.NPI = byte(rand.Int() % 11)
		tmp := randDigit(rand.Int() % 20)
		var e error
		a.Addr, e = teldata.ParseTBCD(tmp)
		if e != nil {
			panic(e)
		}
	}
	return
}

func TestNilAddr(t *testing.T) {
	for i := 0; i < 100; i++ {
		orig := randAddress()
		orig.Addr = nil
		t.Logf("%s", orig)
	}
}
