package sms

import (
	"fmt"
	"io"
	"time"
)

// Deliver is TPDU message from SC to MS
type Deliver struct {
	MMS bool // More Messages to Send (true=more messages)
	LP  bool // Loop Prevention
	SRI bool // Status Report Indication (true=status report shall be returned)
	RP  bool // Reply Path

	OA   Address         // Originating Address
	PID  byte            // Protocol Identifier
	DCS  dcs             // Data Coding Scheme
	SCTS time.Time       // Service Centre Time Stamp
	UDH  map[byte][]byte // User Data Header
	UD   []byte          // User Data
}

// WriteTo output byte data of this TPDU
func (d *Deliver) WriteTo(w io.Writer) (n int64, e error) {
	i := 0
	b := []byte{0x00}
	if !d.MMS {
		b[0] = b[0] | 0x04
	}
	if d.LP {
		b[0] = b[0] | 0x08
	}
	if d.SRI {
		b[0] = b[0] | 0x20
	}
	if d.UDH != nil && len(d.UDH) != 0 {
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

	udh := encodeUDH(d.UDH)
	u := d.DCS.unitSize()
	l := len(udh) + len(d.UD)
	l = ((l * 8) - (l * 8 % u)) / u

	b = make([]byte, 10)
	b[0] = d.PID
	b[1] = d.DCS.encodeDCS()
	for j, k := range encodeTime(d.SCTS) {
		b[j+2] = k
	}
	b[9] = byte(l)
	if i, e = w.Write(b); e != nil {
		return
	}
	n += int64(i)

	if i, e = w.Write(udh); e != nil {
		return
	}
	n += int64(i)
	if i, e = w.Write(d.UD); e != nil {
		return
	}
	n += int64(i)

	return
}

func (d *Deliver) readFrom(h byte, r io.Reader) (n int64, e error) {
	d.MMS = h&0x04 != 0x04
	d.LP = h&0x08 == 0x08
	d.SRI = h&0x20 == 0x20
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
	d.DCS = decodeDCS(b[1])
	if d.DCS == nil {
		e = fmt.Errorf("invalid TP-DCS data: % x", b[1])
		return
	}
	t := [7]byte{}
	for j := 0; j < 7; j++ {
		t[j] = b[j+2]
	}
	d.SCTS = decodeTime(t)

	l := d.DCS.unitSize()
	l *= int(b[9])
	if l%8 != 0 {
		l += 8 - l%8
	}

	d.UD = make([]byte, l/8)
	if i, e = r.Read(d.UD); e != nil {
		return
	} else if i != len(d.UD) {
		e = fmt.Errorf("more data required")
		return
	}
	n += int64(i)

	if h&0x40 == 0x40 {
		d.UDH = decodeUDH(d.UD[0 : d.UD[0]+1])
		d.UD = d.UD[d.UD[0]+1:]
	}

	return
}

// PrintStack show PDU parameter
func (d *Deliver) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "SMS message stack: Deliver\n")
	fmt.Fprintf(w, "TP-MMS:  %s\n", mmsStat(d.MMS))
	fmt.Fprintf(w, "TP-LP:   %s\n", lpStat(d.LP))
	fmt.Fprintf(w, "TP-SRI:  %s\n", sriStat(d.SRI))
	fmt.Fprintf(w, "TP-RP:   %s\n", rpStat(d.RP))

	fmt.Fprintf(w, "TP-OA:   %s\n", d.OA)
	fmt.Fprintf(w, "TP-PID:  %d\n", d.PID)
	fmt.Fprintf(w, "TP-DCS:  %s\n", d.DCS)
	fmt.Fprintf(w, "TP-SCTS: %s\n", d.SCTS)
	fmt.Fprintf(w, "TP-UD:\n")
	for k, v := range d.UDH {
		fmt.Fprintf(w, " IEI:%d = % x\n", k, v)
	}

	if d.UD != nil && len(d.UD) != 0 {
		fmt.Fprintf(w, "%s\n", d.DCS.decodeData(d.UD))
	}
}
