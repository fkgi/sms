package sms_test

import (
	"encoding/json"
	"testing"

	"github.com/fkgi/sms"
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
