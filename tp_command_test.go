package sms_test

import (
	"encoding/json"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func TestEncodeCommand(t *testing.T) {

	bytedata := []byte{
		0x42, 0x00, 0x00, 0x00, 0x00, 0x04, 0x80, 0x21,
		0x43, 0x0d, 0x05, 0x00, 0x03, 0x84, 0x0a, 0x01,
		0x5e, 0x93, 0xdd, 0x51, 0x7b, 0x86, 0x50}
	p, e := sms.UnmarshalTPMO(bytedata)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}

	t.Log(p.String())
}

func TestDecodeCommand(t *testing.T) {
	p := &sms.Command{
		SRR: false,
		TMR: 0x00,
		PID: 0x00,
		CT:  0x00,
		MN:  0x00,
		DA:  sms.Address{TON: 0, NPI: 0},
		CD:  sms.UserData{Text: "XpPdUXuGUA=="}}
	p.CD.UDH = append(p.CD.UDH, sms.ConcatenatedSM{
		RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
	p.DA.Addr, _ = teldata.ParseTBCD("1234")

	b := p.MarshalTP()
	t.Logf("% x", b)
}

func TestMarshalJSON_command(t *testing.T) {
	p := &sms.Command{
		SRR: false,
		TMR: 0x00,
		PID: 0x00,
		CT:  0x00,
		MN:  0x00,
		DA:  sms.Address{TON: 0, NPI: 0},
		CD:  sms.UserData{Text: "XpPdUXuGUA=="}}
	p.CD.UDH = append(p.CD.UDH, sms.ConcatenatedSM{
		RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
	p.DA.Addr, _ = teldata.ParseTBCD("1234")
	t.Log(p.String())

	var e error
	bytedata, e := json.Marshal(p)
	if e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(string(bytedata))

	if e := json.Unmarshal(bytedata, &p); e != nil {
		t.Fatalf("unmarshal failed: %s", e)
	}
	t.Log(p.String())
}

func randCommand() sms.Command {
	orig := sms.Command{
		SRR: randBool(),
		TMR: randByte(),
		PID: randByte(),
		CT:  randByte(),
		MN:  randByte(),
		DA:  randAddress()}
	orig.CD = randUD(sms.GeneralDataCoding{MsgCharset: sms.Charset8bitData})

	orig.TI = randTransactionID()
	orig.RMR = randByte()
	orig.SCA = sms.Address{
		TON: sms.TypeInternational,
		NPI: sms.PlanISDNTelephone}
	tmp := randDigit((rand.Int() % 20) + 1)
	var e error
	if orig.SCA.Addr, e = teldata.ParseTBCD(tmp); e != nil {
		panic(e)
	}
	return orig
}

func compareTPCommand(orig, ocom sms.Command) error {
	if orig.SRR != ocom.SRR {
		return errors.New("SRR mismatch")
	}
	if orig.TMR != ocom.TMR {
		return errors.New("TMR mismatch")
	}
	if orig.PID != ocom.PID {
		return errors.New("PID mismatch")
	}
	if orig.CT != ocom.CT {
		return errors.New("CT mismatch")
	}
	if orig.MN != ocom.MN {
		return errors.New("MN mismatch")
	}
	if !orig.DA.Equal(ocom.DA) {
		return errors.New("DA mismatch")
	}
	if !orig.CD.Equal(ocom.CD) {
		return errors.New("CD text mismatch")
	}
	return nil
}

func compareRPCommand(orig, ocom sms.Command) error {
	if orig.RMR != ocom.RMR {
		return errors.New("MR mismatch")
	}
	if !orig.SCA.Equal(ocom.SCA) {
		return errors.New("SCA mismatch")
	}
	return compareTPCommand(orig, ocom)
}

func compareCPCommand(orig, ocom sms.Command) error {
	if orig.TI != ocom.TI {
		return errors.New("TI mismatch")
	}
	return compareRPCommand(orig, ocom)
}

func TestConvertTPCommand(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randCommand()

		t.Logf("%s", orig)
		b := orig.MarshalTP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalTPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Command)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareTPCommand(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Command{}
		e = ocom.UnmarshalTP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareTPCommand(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertRPCommand(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randCommand()

		t.Logf("%s", orig)
		b := orig.MarshalRP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalRPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Command)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareRPCommand(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Command{}
		e = ocom.UnmarshalRP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareRPCommand(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestConvertCPCommand(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		orig := randCommand()

		t.Logf("%s", orig)
		b := orig.MarshalCP()
		t.Logf("% x", b)
		res, e := sms.UnmarshalCPMO(b)
		if e != nil {
			t.Fatal(e)
		}
		ocom, ok := res.(sms.Command)
		if !ok {
			t.Fatal("mti mismatch")
		}
		t.Logf("%s", ocom)

		e = compareCPCommand(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}

		ocom = sms.Command{}
		e = ocom.UnmarshalCP(b)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("%s", ocom)

		e = compareCPCommand(orig, ocom)
		if e != nil {
			t.Fatal(e)
		}
	}
}
