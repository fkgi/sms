package sms_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/fkgi/sms"
)

func TestMarshalJSON_multiUDH(t *testing.T) {
	u := sms.UserData{
		Text: "hoge",
	}
	u.UDH = append(u.UDH, sms.ConcatenatedSM{RefNum: 3, MaxNum: 2, SeqNum: 1})
	u.UDH = append(u.UDH, sms.GenericIEI{K: 99, V: []byte{0x01}})
	t.Log(u.String())
	subfuncMarshalJSON(u, t)
}

func TestMarshalJSON_emptytxt(t *testing.T) {
	u := sms.UserData{}
	u.UDH = append(u.UDH, sms.ConcatenatedSM{RefNum: 3, MaxNum: 2, SeqNum: 1})
	t.Log(u.String())
	subfuncMarshalJSON(u, t)
}

func TestMarshalJSON_emptyhdr(t *testing.T) {
	u := sms.UserData{
		Text: "hoge",
	}
	t.Log(u.String())
	subfuncMarshalJSON(u, t)
}

func TestMarshalJSON_empty(t *testing.T) {
	u := sms.UserData{}
	t.Log(u.String())
	subfuncMarshalJSON(u, t)
}

func subfuncMarshalJSON(u sms.UserData, t *testing.T) {
	var e error
	bytedata, e := json.Marshal(u)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))

	if e := json.Unmarshal(bytedata, &u); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(u.String())
}

func TestConvertUDH(t *testing.T) {
	origs := make([]sms.UserDataHdr, rand.Int()%10)
	for i := range origs {
		origs[i] = randUDH()
	}
	b := sms.MarshalUDHs(origs)
	t.Logf("\ndata=% x", b)

	ocoms := sms.UnmarshalUDHs(b)
	for i := range ocoms {
		t.Logf("%s", origs[i])
		t.Logf("%s", ocoms[i])
		if !origs[i].Equal(ocoms[i]) {
			t.Fatalf("mismatch orig=%s ocom=%s", origs[i], ocoms[i])
		}
	}
}

func randUDH() sms.UserDataHdr {
	h := randByte()
	switch h {
	case 0x00:
		return sms.ConcatenatedSM{
			RefNum: randByte(),
			MaxNum: randByte(),
			SeqNum: randByte()}
	case 0x08:
		return sms.ConcatenatedSM16bit{
			RefNum: uint16(rand.Int31n(65536)),
			MaxNum: randByte(),
			SeqNum: randByte()}
	default:
		iei := sms.GenericIEI{
			K: h,
			V: make([]byte, rand.Int31n(5))}
		for i := range iei.V {
			iei.V[i] = randByte()
		}
		return iei
	}
}

/*
func TestConvertUD(t *testing.T) {
	for i := 0; i < 500; i++ {
		d := getRandomDCS()
		orig := getRandomUDText(d)
		t.Logf("%s", orig)
		b := orig.Marshal()
		t.Logf("\ndata=% x", b)
		ocom := sms.UnmarshalUD(b)
		t.Logf("%s", ocom)
		if !orig.Equal(ocom) {
			t.Fatalf("mismatch orig=%s ocom=%s", orig, ocom)
		}
	}
}
*/

func randUD(d sms.DataCoding) sms.UserData {
	u := sms.UserData{}
	u.UDH = make([]sms.UserDataHdr, rand.Int31n(5))
	for i := range u.UDH {
		u.UDH[i] = randUDH()
	}
	l := len(sms.MarshalUDHs(u.UDH))
	l = 140 - l

	switch d.Charset() {
	case sms.CharsetGSM7bit:
		l = l / 7
		l++
		u.Text = randText(rand.Int() % l)
	case sms.Charset8bitData:
		l = l / 8
		l++
		tmp := make([]byte, rand.Int()%l)
		for i := range tmp {
			tmp[i] = randByte()
		}
		u.Set8bitData(tmp)
	case sms.CharsetUCS2:
		l = l / 8
		l++
		tmp := make([]rune, rand.Int()%l)
		for i := range tmp {
			for !unicode.IsPrint(tmp[i]) {
				tmp[i] = int32(rand.Int() % 2147483648)
			}
		}
		u.Text = string(tmp)
	}
	return u
}

func TestMakeSeparatedText(t *testing.T) {
	for j := 0; j < 100; j++ {
		var txt string
		if randBool() {
			txt = randText(int(rand.Int31n(500)))
		} else {
			tmp := make([]rune, int(rand.Int31n(500)))
			for i := range tmp {
				for !unicode.IsPrint(tmp[i]) {
					tmp[i] = int32(rand.Int() % 2147483648)
				}
			}
			txt = string(tmp)
		}
		ud, cs := sms.MakeSeparatedText(txt, randByte())

		t.Log("orig=", txt)
		switch cs {
		case sms.CharsetGSM7bit:
			t.Log("char=GSM-7bit")
		case sms.CharsetUCS2:
			t.Log("char=UCS2")
		}
		ocom := ""
		for i, u := range ud {
			t.Log("txt[", i, "]=", u)
			ocom += u.Text

			if cs == sms.CharsetGSM7bit {
				g7s, e := sms.StringToGSM7bit(u.Text)
				if e != nil {
					t.Fatal(e)
				}
				if len(ud) == 1 && g7s.Length() > 160 {
					t.Fatal("too long text", g7s.Length())
				}
				if len(ud) > 1 && g7s.Length() > 153 {
					t.Fatal("too long text", g7s.Length())
				}
			} else {
				if len(ud) == 1 && utf8.RuneCountInString(u.Text) > 70 {
					t.Fatal("too long text", len(u.Text))
				}
				if len(ud) > 1 && utf8.RuneCountInString(u.Text) > 67 {
					t.Fatal("too long text", len(u.Text))
				}

			}
		}

		if ocom != txt {
			t.Fatal("text missmatch")
		}
	}
}
