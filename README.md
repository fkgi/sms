# sms
TPDU (3GPP TS23.040), RPDU and CPDU (3GPP TS24.011) message encoding implementation by Golang.

[![GoDoc](https://godoc.org/github.com/fkgi/sms?status.svg)](https://godoc.org/github.com/fkgi/sms)

## About SMS

[SMS](https://en.wikipedia.org/wiki/GSM_03.40) is a text messaging service component of mobile device systems.

## Project Status

This project is still WIP.

# Getting Started

## Installation

```shell-session
go get -u github.com/fkgi/teldata
go get -u github.com/fkgi/sms
```

## Running examples

Encode TP-DELIVER as below.

```go
p := sms.Deliver{
	MMS: true,
	LP:  false,
	SRI: false,
	RP:  false,
	OA:  sms.Address{TON: 0, NPI: 0},
	PID: 0,
	DCS: &sms.GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   sms.NoMessageClass,
		MsgCharset: sms.CharsetUCS2},
	SCTS: time.Date(
		2011, time.March, 22, 14, 25, 40, 0,
		time.FixedZone("unknown", 9*60*60)),
	UD: sms.UserData{Text: "あいうえお"}}
p.UD.UDH = append(p.UD.UDH, sms.ConcatenatedSM{
	RefNum: 0x84, MaxNum: 0x0a, SeqNum: 0x01})
p.OA.Addr, _ = teldata.ParseTBCD("1234")

b := p.MarshalTP()
fmt.Printf("% x", b)
```

Decode TP-DELIVER as below.

```go
bytedata := []byte{
	0x40, 0x04, 0x80, 0x21, 0x43, 0x00, 0x08, 0x11,
	0x30, 0x22, 0x41, 0x52, 0x04, 0x63, 0x10, 0x05,
	0x00, 0x03, 0x87, 0x02, 0x01, 0x30, 0x42, 0x30,
	0x44, 0x30, 0x46, 0x30, 0x48, 0x30, 0x4a}
p, e := sms.UnmarshalTPMT(bytedata)
if e != nil {
	fmt.Printf("encode failed: %s", e)
}
fmt.Println(p.String())
```

Refer each _test.go files to see each message decoding/encoding.

# LICENSE

[MIT](https://github.com/fkgi/sms/blob/master/LICENSE)
