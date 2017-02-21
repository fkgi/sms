package sms

import (
	"bytes"
	"testing"
)

func TestEncodeSubmit(t *testing.T) {

	bytedata := []byte{
		0x41, 0x40, 0x0b, 0x81, 0x90, 0x10, 0x32, 0x54,
		0x76, 0xf8, 0x00, 0x08, 0x10, 0x05, 0x00, 0x03,
		0x84, 0x0a, 0x01, 0x30, 0x42, 0x30, 0x44, 0x30,
		0x46, 0x30, 0x48, 0x30, 0x4a, 0xff, 0xff, 0xff}
	buf := bytes.NewBuffer(bytedata)
	p, _, e := ReadAsSC(buf)
	if e != nil {
		t.Fatalf("encode failed: %s", e)
	}

	b := new(bytes.Buffer)
	p.PrintStack(b)
	t.Log(b.String())
}

func TestDecodeSubmit(t *testing.T) {
	p := &Submit{}
	p.RD = false
	p.SRR = false
	p.RP = false
	p.MR = 64
	p.DA = Address{TON: 0, NPI: 1}
	p.DA.Addr, _ = ParseTBCD("09012345678")
	p.PID = 0
	p.DCS = &GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   NoMessageClass,
		Charset:    UCS2}
	p.VP = nil
	p.UDH = make(map[byte][]byte)
	p.UDH[0] = []byte{0x84, 0x0a, 0x01}
	p.UD, _ = p.DCS.encodeData("あいうえお")

	b := new(bytes.Buffer)
	_, e := p.WriteTo(b)
	if e != nil {
		t.Fatalf("deecode failed: %s", e)
	}

	t.Logf("% x", b)
}
