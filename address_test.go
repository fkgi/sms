package sms_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

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
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig, e := genRandomAddress()
		if e != nil {
			t.Fatal(e)
		}
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

func genRandomAddress() (a sms.Address, e error) {
	a.TON = byte(rand.Int() % 7)
	if a.TON == sms.TypeAlphanumeric {
		a.NPI = sms.PlanUnknown
		tmp := randText(rand.Int() % 10)
		a.Addr, e = sms.StringToGSM7bit(tmp)
	} else {
		a.NPI = byte(rand.Int() % 11)
		tmp := randDigit(rand.Int() % 10)
		a.Addr, e = teldata.ParseTBCD(tmp)
	}
	return
}
