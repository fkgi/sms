package sms

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// Submit is TPDU message from MS to SC
type Submit struct {
	RD  bool // Reject Duplicates
	SRR bool // Status Report Request
	RP  bool // Reply Path

	MR  byte    // Message Reference
	DA  Address // Destination Address
	PID byte    // Protocol Identifier
	DCS dcs     // Data Coding Scheme
	VP  vp      // Validity Period
	UDH []udh   // User Data Header
	UD  []byte  // User Data
}

// Encode output byte data of this TPDU
func (d *Submit) Encode() []byte {
	w := new(bytes.Buffer)

	b := byte(0x01)
	if d.RD {
		b = b | 0x04
	}
	var vp []byte
	switch v := d.VP.(type) {
	case VPRelative:
		b = b | 0x10
		vp = []byte{v[0]}
	case VPEnhanced:
		b = b | 0x08
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	case VPAbsolute:
		b = b | 0x18
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	}
	if d.SRR {
		b = b | 0x20
	}
	if len(d.UDH) != 0 {
		b = b | 0x40
	}
	if d.RP {
		b = b | 0x80
	}
	w.WriteByte(b)

	w.WriteByte(d.MR)

	d.DA.WriteTo(w)

	w.WriteByte(d.PID)
	w.WriteByte(d.DCS.encode())
	w.Write(vp)

	writeUD(w, d.UD, d.UDH, d.DCS)

	return w.Bytes()
}

// Decode get data of this TPDU
func (d *Submit) Decode(b []byte) (e error) {
	d.RD = b[0]&0x04 == 0x04
	d.SRR = b[0]&0x20 == 0x20
	d.RP = b[0]&0x80 == 0x80

	r := bytes.NewReader(b[1:])

	if d.MR, e = r.ReadByte(); e != nil {
		return
	}
	if _, e = d.DA.ReadFrom(r); e != nil {
		return
	}
	if d.PID, e = r.ReadByte(); e != nil {
		return
	}
	if d.DCS, e = readDCS(r); e != nil {
		return
	}
	switch b[0] & 0x18 {
	case 0x00:
		d.VP = nil
	case 0x10:
		var p byte
		if p, e = r.ReadByte(); e == nil {
			d.VP = VPRelative([1]byte{p})
		}
	case 0x08:
		var p [7]byte
		if p, e = read7Bytes(r); e == nil {
			d.VP = VPEnhanced(p)
		}
	case 0x18:
		var p [7]byte
		if p, e = read7Bytes(r); e == nil {
			d.VP = VPAbsolute(p)
		}
	}
	if e == nil {
		d.UD, d.UDH, e = readUD(r, d.DCS, b[0]&0x40 == 0x40)
	}
	if e == nil && r.Len() != 0 {
		e = fmt.Errorf("invalid data: extra data")
	}
	return
}

// PrintStack show PDU parameter
func (d *Submit) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "SMS message stack: Submit\n")
	fmt.Fprintf(w, " | TP-RD:   %s\n", rdStat(d.RD))
	fmt.Fprintf(w, " | TP-SRR:  %s\n", srrStat(d.SRR))
	fmt.Fprintf(w, " | TP-RP:   %s\n", rpStat(d.RP))
	fmt.Fprintf(w, " | TP-MR:   %d\n", d.MR)
	fmt.Fprintf(w, " | TP-DA:   %s\n", d.DA)
	fmt.Fprintf(w, " | TP-PID:  %s\n", pidStat(d.PID))
	fmt.Fprintf(w, " | TP-DCS:  %s\n", d.DCS)
	if d.VP != nil {
		fmt.Fprintf(w, " | TP-VP:   %s\n", d.VP)
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
}

// SubmitReport is TPDU message from SC to MS
type SubmitReport struct {
	FCS  *byte     // Failure Cause
	SCTS time.Time // Service Centre Time Stamp
	PID  *byte     // Protocol Identifier
	DCS  dcs       // Data Coding Scheme
	UDH  []udh     // User Data Header
	UD   []byte    // User Data
}

// Encode output byte data of this TPDU
func (d *SubmitReport) Encode() []byte {
	w := new(bytes.Buffer)

	b := byte(0x01)
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
	w.Write(encodeSCTimeStamp(d.SCTS))
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
func (d *SubmitReport) Decode(b []byte) (e error) {
	r := bytes.NewReader(b[1:])

	var pi byte
	if pi, e = r.ReadByte(); e == nil && pi&0x80 == 0x80 {
		d.FCS = &pi
		pi, e = r.ReadByte()
	}
	if e != nil {
		return
	}
	var p [7]byte
	if p, e = read7Bytes(r); e != nil {
		return
	}
	d.SCTS = decodeSCTimeStamp(p)
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

// PrintStack show PDU parameter
func (d *SubmitReport) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "SMS message stack: Submit Report")
	if d.FCS != nil {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, " | TP-FCS:  %s\n", fcsStat(*d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
	}

	fmt.Fprintf(w, " | TP-SCTS: %s\n", d.SCTS)
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
}
