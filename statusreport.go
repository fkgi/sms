package sms

import (
	"bytes"
	"fmt"
	"time"
)

// StatusReport is TPDU message from SC to MS
type StatusReport struct {
	MMS bool // More Messages to Send (true=more messages)
	LP  bool // Loop Prevention
	SRQ bool // Status Report Qualifier (true=status report shall be returned)

	MR   byte      // Message Reference
	RA   Address   // Destination Address
	SCTS time.Time // Service Centre Time Stamp
	DT   time.Time // Discharge Time
	ST   byte      // Status
	PID  *byte     // Protocol Identifier
	DCS  dcs       // Data Coding Scheme
	UDH  []udh     // User Data Header
	UD   []byte    // User Data
}

// Encode output byte data of this TPDU
func (d *StatusReport) Encode() []byte {
	w := new(bytes.Buffer)

	b := byte(0x02)
	if !d.MMS {
		b = b | 0x04
	}
	if d.LP {
		b = b | 0x08
	}
	if d.SRQ {
		b = b | 0x20
	}
	if d.UDH != nil && len(d.UDH) != 0 {
		b = b | 0x40
	}
	w.WriteByte(b)
	w.WriteByte(d.MR)
	b, a := d.RA.Encode()
	w.WriteByte(b)
	w.Write(a)
	w.Write(encodeSCTimeStamp(d.SCTS))
	w.Write(encodeSCTimeStamp(d.DT))
	w.WriteByte(d.ST)
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
	if b == 0x00 {
		return w.Bytes()
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
func (d *StatusReport) Decode(b []byte) (e error) {
	d.MMS = b[0]&0x04 != 0x04
	d.LP = b[0]&0x08 == 0x08
	d.SRQ = b[0]&0x20 == 0x20

	r := bytes.NewReader(b[1:])

	if d.MR, e = r.ReadByte(); e != nil {
		return
	}
	if d.RA, e = readAddr(r); e != nil {
		return
	}
	var p [7]byte
	if p, e = read7Bytes(r); e != nil {
		return
	}
	d.SCTS = decodeSCTimeStamp(p)
	if p, e = read7Bytes(r); e != nil {
		return
	}
	d.DT = decodeSCTimeStamp(p)
	if d.ST, e = r.ReadByte(); e != nil {
		return
	}
	if r.Len() == 0 {
		return
	}
	var pi byte
	if pi, e = r.ReadByte(); e != nil {
		e = nil
		return
	}
	b = make([]byte, 1)
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
		tmp := make([]byte, r.Len())
		r.Read(tmp)
		e = &InvalidDataError{
			Name:  "extra part",
			Bytes: tmp}
	}
	return
}

func (d *StatusReport) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Status Report\n")
	fmt.Fprintf(w, "%sTP-MMS:  %s\n", Indent, mmsStat(d.MMS))
	fmt.Fprintf(w, "%sTP-LP:   %s\n", Indent, lpStat(d.LP))
	fmt.Fprintf(w, "%sTP-SRQ:  %s\n", Indent, srqStat(d.SRQ))
	fmt.Fprintf(w, "%sTP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sTP-RA:   %s\n", Indent, d.RA)
	fmt.Fprintf(w, "%sTP-SCTS: %s\n", Indent, d.SCTS)
	fmt.Fprintf(w, "%sTP-DT:   %s\n", Indent, d.SCTS)
	fmt.Fprintf(w, "%sTP-ST:   %s\n", Indent, stStat(d.ST))

	if d.PID != nil {
		fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(*d.PID))
	}
	if d.DCS != nil {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	}
	if len(d.UDH)+len(d.UD) != 0 {
		fmt.Fprintf(w, "%sTP-UD:\n", Indent)
		for _, h := range d.UDH {
			fmt.Fprintf(w, "%s%s%s\n", Indent, Indent, h)
		}
		if len(d.UD) != 0 {
			fmt.Fprintf(w, "%s%s%s\n", Indent, Indent, d.DCS.Decode(d.UD))
		}
	}
	return w.String()
}
