package sms

import (
	"bytes"
	"testing"
	"time"
)

func TestEncodeDeliver(t *testing.T) {

	bytedata := []byte{
		0x40, 0x04, 0x80, 0x21, 0x43, 0x00, 0x08, 0x11,
		0x30, 0x22, 0x41, 0x52, 0x04, 0x63, 0x10, 0x05,
		0x00, 0x03, 0x87, 0x02, 0x01, 0x30, 0x42, 0x30,
		0x44, 0x30, 0x46, 0x30, 0x48, 0x30, 0x4a, 0xff}
	buf := bytes.NewBuffer(bytedata)
	p, _, e := ReadAsSM(buf, false)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}
	b := new(bytes.Buffer)
	p.PrintStack(b)
	t.Log(b.String())
}

func TestDecodeDeliver(t *testing.T) {
	p := &Deliver{}
	p.MMS = true
	p.LP = false
	p.SRI = false
	p.RP = false
	p.OA = Address{TON: 0, NPI: 0}
	p.OA.Addr, _ = ParseTBCD("1234")
	p.PID = 0
	p.DCS = &GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   NoMessageClass,
		Charset:    UCS2}
	p.SCTS = time.Date(
		2011, time.March, 22, 14, 25, 40, 0,
		time.FixedZone("unknown", 9*60*60))
	p.UDH = []udh{&ConcatenatedSM{0x84, 0x0a, 0x01}}
	p.UD, _ = p.DCS.encodeData("あいうえお")

	b := new(bytes.Buffer)
	_, e := p.WriteTo(b)
	if e != nil {
		t.Fatalf("deecode failed: %s", e)
	}

	t.Logf("% x", b)
}
