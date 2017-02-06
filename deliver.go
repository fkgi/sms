package sms

import (
	"bytes"
	"time"
)

func CreateDeliver(src, data string) *TPDU {
	msgRef++

	p := &TPDU{}
	p.Req = true // request message

	p.MTI = 0x00  // message type indicator
	p.MMS = false // more message to send
	p.LP = false  // loop prevention
	p.SRI = true  // status report indication
	p.RP = false  // reply path

	p.OA = TPAddr{true, 0, 1, src} // originating address
	p.PID = 0x00                   // protocol identifer Default store and forward short message
	p.DCS = 0x08                   // data coding scheme UCS2 with no Message Class
	p.SCTS = getTime(time.Now())   // service center time stamp

	p.UD = data // user data

	return p
}

func (p TPDU) encodeDeliver() []byte {
	var buf bytes.Buffer

	b := byte(p.MTI)
	if !p.MMS {
		b = b | 0x04
	}
	if p.LP {
		b = b | 0x08
	}
	if p.SRI {
		b = b | 0x20
	}
	if p.UDh != nil {
		b = b | 0x40
	}
	if p.RP {
		b = b | 0x80
	}
	buf.WriteByte(b)

	buf.Write(p.OA.encode())
	buf.WriteByte(p.PID)
	buf.WriteByte(p.DCS)
	buf.Write(p.SCTS)

	buf.WriteByte(byte(len(p.UD) * 2))
	buf.Write(p.encodeUD())

	return buf.Bytes()
}

func (p TPDU) encodeDeliverReport() []byte {
	var buf bytes.Buffer

	return buf.Bytes()
}
