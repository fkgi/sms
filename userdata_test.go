package sms_test

import (
	"encoding/json"
	"testing"

	"github.com/fkgi/sms"
)

func TestMarshalJSON(t *testing.T) {
	u := sms.UD{
		Text: "hoge",
	}
	u.AddUDH(&sms.ConcatenatedSM{RefNum: 3, MaxNum: 2, SeqNum: 1})
	u.AddUDH(&sms.GenericIEI{K: 99, V: []byte{0x01}})
	t.Log(u.String())

	var e error
	bytedata, e = json.Marshal(u)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))

	if e := json.Unmarshal(bytedata, &u); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(u.String())
}
