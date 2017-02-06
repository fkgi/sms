package sms

import (
	"bytes"
)

func CreateSubmit(dest, data string) *TPDU {
	msgRef++

	p := &TPDU{}
	p.Req = true // request message

	p.MTI = 0x01  // message type indicator
	p.RD = false  // reject duplicates
	p.VPF = 0x00  // validity period format
	p.SRR = false // status report request
	p.RP = false  // reply path

	p.MR = msgRef                   // message reference
	p.DA = TPAddr{true, 0, 1, dest} // destination address
	p.PID = 0x00                    // protocol identifer Default store and forward short message
	p.DCS = 0x08                    // data coding scheme UCS2 with no Message Class

	p.UD = data // user data

	return p
}

func (p TPDU) encodeSubmit() []byte {
	var buf bytes.Buffer

	b := byte(p.MTI)
	if p.RD {
		b = b | 0x04
	}
	b = b | ((p.VPF & 0x03) << 3)
	if p.SRR {
		b = b | 0x20
	}
	if p.UDh != nil {
		b = b | 0x40
	}
	if p.RP {
		b = b | 0x80
	}
	buf.WriteByte(b)

	buf.WriteByte(p.MR)
	buf.Write(p.DA.encode())
	buf.WriteByte(p.PID)
	buf.WriteByte(p.DCS)

	if p.VPF != 0x00 {
		buf.Write(p.VP)
	}

	buf.WriteByte(byte(len(p.UD) * 2))
	buf.Write(p.encodeUD())

	return buf.Bytes()
}

func (p TPDU) encodeSubmitReport() []byte {
	var buf bytes.Buffer

	return buf.Bytes()
}
