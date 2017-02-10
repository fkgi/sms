package sms

import (
	"fmt"
	"io"
)

// Deliver is TPDU message from SC to MS
type Deliver struct {
	MTI  byte // Message Type Indicator
	MMS  bool // More Messages to Send (true=more messages)
	LP   bool // Loop Prevention
	SRI  bool // Status Report Indication (true=status report shall be returned)
	UDHI bool // User Data Header Indicator
	RP   bool // Reply Path

	OA   Address   // Originating Address
	PID  byte      // Protocol Identifier
	DCS  byte      // Data Coding Scheme
	SCTS TimeStamp // Service Centre Time Stamp
	UDL  byte      // User Data Length
	UD   []byte    // User Data
}

// WriteTo output byte data of this TPDU
func (d *Deliver) WriteTo(w io.Writer) (n int64, e error) {
	i := 0
	b := []byte{d.MTI}
	if !d.MMS {
		b[0] = b[0] | 0x04
	}
	if d.LP {
		b[0] = b[0] | 0x08
	}
	if d.SRI {
		b[0] = b[0] | 0x20
	}
	if d.UDHI {
		b[0] = b[0] | 0x40
	}
	if d.RP {
		b[0] = b[0] | 0x80
	}
	if i, e = w.Write(b); e != nil {
		return
	}

	if n, e = d.OA.WriteTo(w); e != nil {
		return
	}
	n += int64(i)

	b = make([]byte, 10)
	b[0] = d.PID
	b[1] = d.DCS
	for j := 0; j < 7; j++ {
		b[j+2] = d.SCTS[j]
	}
	b[9] = d.UDL
	if i, e = w.Write(b); e != nil {
		return
	}
	n += int64(i)

	if i, e = w.Write(d.UD); e != nil {
		return
	}
	n += int64(i)

	return
}

// ReadFrom read byte data and set parameter of the TPDU
func (d *Deliver) ReadFrom(h byte, r io.Reader) (n int64, e error) {
	d.MTI = h & 0x03
	d.MMS = h&0x04 != 0x04
	d.LP = h&0x08 == 0x08
	d.SRI = h&0x20 == 0x20
	d.UDHI = h&0x40 == 0x40
	d.RP = h&0x80 == 0x80

	d.OA = Address{}
	if n, e = d.OA.ReadFrom(r); e != nil {
		return
	}

	i := 0
	b := make([]byte, 10)
	if i, e = r.Read(b); e != nil {
		return
	} else if i != len(b) {
		e = fmt.Errorf("more data required")
		return
	}
	n += int64(i)

	d.PID = b[0]
	d.DCS = b[1]
	for j := 0; j < 7; j++ {
		d.SCTS[j] = b[j+2]
	}
	d.UDL = b[9]

	d.UD = make([]byte, int(d.UDL))
	if i, e = r.Read(d.UD); e != nil {
		return
	} else if i != len(d.UD) {
		e = fmt.Errorf("more data required")
		return
	}
	n += int64(i)

	return
}

/*
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
*/
