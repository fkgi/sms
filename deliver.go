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

	OA   Address   // Originating Address
	PID  byte      // Protocol Identifier
	DCS  dcs       // Data Coding Scheme
	SCTS time.Time // Service Centre Time Stamp
	UDH  []udh     // User Data Header
	UD   []byte    // User Data
}

// WriteTo output byte data of this TPDU
func (d *Deliver) WriteTo(w io.Writer) (n int64, e error) {
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
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	var nn int64
	nn, e = d.OA.WriteTo(w)
	n += nn
	if e != nil {
		return
	}

	udh := encodeUDH(d.UDH)
	u := d.DCS.unitSize()
	l := len(udh) + len(d.UD)
	l = ((l * 8) - (l * 8 % u)) / u
	b = []byte{d.PID, d.DCS.encode()}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}
	b = encodeSCTimeStamp(d.SCTS)
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}
	b = []byte{byte(l)}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}
	if n, e = writeBytes(w, n, udh); e != nil {
		return
	}
	n, e = writeBytes(w, n, d.UD)
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
	d.SCTS = decodeSCTimeStamp([7]byte{b[2], b[3], b[4], b[5], b[6], b[7], b[8]})

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
	for _, h := range d.UDH {
		fmt.Fprintf(w, " %s\n", h)
	}
	if d.UD != nil && len(d.UD) != 0 {
		fmt.Fprintf(w, "%s\n", d.DCS.decodeData(d.UD))
	}
}

// DeliverReport is TPDU message from MS to SC
type DeliverReport struct {
	FCS byte            // Failure Cause
	PI  byte            // Parameter Indicator
	PID byte            // Protocol Identifier
	DCS dcs             // Data Coding Scheme
	UDH map[byte][]byte // User Data Header
	UD  []byte          // User Data
}
