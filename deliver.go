package sms

import (
	"bytes"
	"fmt"
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

// Encode output byte data of this TPDU
func (d *Deliver) Encode() []byte {
	w := new(bytes.Buffer)

	b := byte(0x01)
	if !d.MMS {
		b = b | 0x04
	}
	if d.LP {
		b = b | 0x08
	}
	if d.SRI {
		b = b | 0x20
	}
	if d.UDH != nil && len(d.UDH) != 0 {
		b = b | 0x40
	}
	if d.RP {
		b = b | 0x80
	}
	w.WriteByte(b)
	d.OA.WriteTo(w)
	w.WriteByte(d.PID)
	w.WriteByte(d.DCS.encode())
	w.Write(encodeSCTimeStamp(d.SCTS))
	writeUD(w, d.UD, d.UDH, d.DCS)

	return w.Bytes()
}

// Decode get data of this TPDU
func (d *Deliver) Decode(b []byte) (e error) {
	d.MMS = b[0]&0x04 != 0x04
	d.LP = b[0]&0x08 == 0x08
	d.SRI = b[0]&0x20 == 0x20
	d.RP = b[0]&0x80 == 0x80

	r := bytes.NewReader(b[1:])

	if _, e = d.OA.ReadFrom(r); e != nil {
		return
	}
	if d.PID, e = r.ReadByte(); e != nil {
		return
	}
	if d.DCS, e = readDCS(r); e != nil {
		return
	}
	var tmp [7]byte
	if tmp, e = read7Bytes(r); e == nil {
		d.SCTS = decodeSCTimeStamp(tmp)
	}
	if e == nil {
		d.UD, d.UDH, e = readUD(r, d.DCS, b[0]&0x40 == 0x40)
	}
	if e == nil && r.Len() != 0 {
		e = fmt.Errorf("invalid data: extra data")
	}

	return
}

func (d *Deliver) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Deliver\n")
	fmt.Fprintf(w, " | TP-MMS:  %s\n", mmsStat(d.MMS))
	fmt.Fprintf(w, " | TP-LP:   %s\n", lpStat(d.LP))
	fmt.Fprintf(w, " | TP-SRI:  %s\n", sriStat(d.SRI))
	fmt.Fprintf(w, " | TP-RP:   %s\n", rpStat(d.RP))
	fmt.Fprintf(w, " | TP-OA:   %s\n", d.OA)
	fmt.Fprintf(w, " | TP-PID:  %s\n", pidStat(d.PID))
	fmt.Fprintf(w, " | TP-DCS:  %s\n", d.DCS)
	fmt.Fprintf(w, " | TP-SCTS: %s\n", d.SCTS)

	if len(d.UDH)+len(d.UD) != 0 {
		fmt.Fprintf(w, " | TP-UD:\n")
		for _, h := range d.UDH {
			fmt.Fprintf(w, "   | %s\n", h)
		}
		if len(d.UD) != 0 {
			fmt.Fprintf(w, "   | %s\n", d.DCS.Decode(d.UD))
		}
	}
	return w.String()
}

// DeliverReport is TPDU message from MS to SC
type DeliverReport struct {
	FCS *byte  // Failure Cause
	PID *byte  // Protocol Identifier
	DCS dcs    // Data Coding Scheme
	UDH []udh  // User Data Header
	UD  []byte // User Data
}

// Encode output byte data of this TPDU
func (d *DeliverReport) Encode() []byte {
	w := new(bytes.Buffer)

	b := byte(0x00)
	if len(d.UDH) != 0 {
		b = b | 0x40
	}
	w.WriteByte(b)
	if d.FCS != nil {
		w.WriteByte(*d.FCS)
	}
	b = byte(0x00)
	if d.PID != nil {
		b = b | 0x01
	}
	if d.DCS != nil {
		b = b | 0x02
	}
	if len(d.UDH)+len(d.UD) != 0 {
		b = b | 0x04
	}
	w.WriteByte(b)
	if d.PID != nil {
		w.WriteByte(*d.PID)
	}
	if d.DCS != nil {
		w.WriteByte(d.DCS.encode())
	}
	if len(d.UDH)+len(d.UD) != 0 {
		writeUD(w, d.UD, d.UDH, d.DCS)
	}
	return w.Bytes()
}

// Decode get data of this TPDU
func (d *DeliverReport) Decode(b []byte) (e error) {
	r := bytes.NewReader(b[1:])

	pi, e := r.ReadByte()
	if e == nil && pi&0x80 == 0x80 {
		d.FCS = &pi
		pi, e = r.ReadByte()
	}
	if e != nil {
		return
	}
	if pi&0x01 == 0x01 {
		var p byte
		if p, e = r.ReadByte(); e != nil {
			return
		}
		d.PID = &p
	}
	if pi&0x02 == 0x02 {
		if d.DCS, e = readDCS(r); e != nil {
			return
		}
	}
	if pi&0x04 == 0x04 {
		if d.DCS == nil {
			d.DCS = &GeneralDataCoding{
				AutoDelete: false,
				Compressed: false,
				MsgClass:   NoMessageClass,
				Charset:    GSM7bitAlphabet}
		}
		d.UD, d.UDH, e = readUD(r, d.DCS, b[0]&0x40 == 0x40)
	}
	if e == nil && r.Len() != 0 {
		e = fmt.Errorf("invalid data: extra data")
	}
	return
}

func (d *DeliverReport) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Deliver Report")
	if d.FCS != nil {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, " | TP-FCS:  %s\n", fcsStat(*d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
	}

	if d.PID != nil {
		fmt.Fprintf(w, " | TP-PID:  %s\n", pidStat(*d.PID))
	}
	if d.DCS != nil {
		fmt.Fprintf(w, " | TP-DCS:  %s\n", d.DCS)
	}
	if len(d.UDH)+len(d.UD) != 0 {
		fmt.Fprintf(w, " | TP-UD:\n")
		for _, h := range d.UDH {
			fmt.Fprintf(w, "   | %s\n", h)
		}
		if len(d.UD) != 0 {
			fmt.Fprintf(w, "   | %s\n", d.DCS.Decode(d.UD))
		}
	}
	return w.String()
}
