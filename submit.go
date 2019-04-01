package sms

import (
	"bytes"
	"fmt"
	"time"
)

// Submit is TPDU message from MS to SC
type Submit struct {
	RD  bool // M / Reject Duplicates
	SRR bool // O / Status Report Request
	RP  bool // M / Reply Path

	MR  byte    // M / Message Reference
	DA  Address // M / Destination Address
	PID byte    // M / Protocol Identifier
	DCS         // M / Data Coding Scheme
	VP          // O / Validity Period
	UD          // O / User Data
}

// Encode output byte data of this TPDU
func (d *Submit) Encode() []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)

	b := byte(0x01)
	if d.RD {
		b |= 0x04
	}
	var vp []byte
	switch v := d.VP.(type) {
	case VPRelative:
		b |= 0x10
		vp = []byte{byte(v)}
	case VPEnhanced:
		b |= 0x08
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	case VPAbsolute:
		b |= 0x18
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	default:
		// nil VP value
	}
	if d.SRR {
		b |= 0x20
	}
	if len(d.UDH) != 0 {
		b |= 0x40
	}
	if d.RP {
		b |= 0x80
	}
	w.WriteByte(b)
	w.WriteByte(d.MR)
	b, a := d.DA.Encode()
	w.WriteByte(b)
	w.Write(a)
	w.WriteByte(d.PID)
	if d.DCS == nil {
		w.WriteByte(0x00)
	} else {
		w.WriteByte(d.DCS.Encode())
	}
	w.Write(vp)
	d.UD.write(w, d.DCS)

	return w.Bytes()
}

// Decode get data of this TPDU
func (d *Submit) Decode(b []byte) (e error) {
	if d == nil {
		return fmt.Errorf("nil data")
	}

	d.RD = b[0]&0x04 == 0x04
	d.SRR = b[0]&0x20 == 0x20
	d.RP = b[0]&0x80 == 0x80

	r := bytes.NewReader(b[1:])
	if d.MR, e = r.ReadByte(); e != nil {
		return
	}
	if d.DA, e = readAddr(r); e != nil {
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
			d.VP = VPRelative(p)
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
		e = d.UD.read(r, d.DCS, b[0]&0x40 == 0x40)
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

func (d *Submit) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Submit\n")
	fmt.Fprintf(w, "%sTP-RD:   %s\n", Indent, rdStat(d.RD))
	fmt.Fprintf(w, "%sTP-SRR:  %s\n", Indent, srrStat(d.SRR))
	fmt.Fprintf(w, "%sTP-RP:   %s\n", Indent, rpStat(d.RP))
	fmt.Fprintf(w, "%sTP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sTP-DA:   %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(d.PID))
	fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	if d.VP != nil {
		fmt.Fprintf(w, "%sTP-VP:   %s\n", Indent, d.VP)
	}
	fmt.Fprintf(w, "%s", d.UD.String())
	return w.String()
}

// SubmitReport is TPDU message from SC to MS
type SubmitReport struct {
	FCS  byte      // C / Failure Cause
	SCTS time.Time // M / Service Centre Time Stamp
	PID  *byte     // O / Protocol Identifier
	DCS            // O / Data Coding Scheme
	UD             // O / User Data
}

// Encode output byte data of this TPDU
func (d *SubmitReport) Encode() []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)

	b := byte(0x01)
	if len(d.UDH) != 0 {
		b |= 0x40
	}
	w.WriteByte(b)
	if d.FCS&0x80 == 0x80 {
		w.WriteByte(d.FCS)
	}
	b = byte(0x00)
	if d.PID != nil {
		b |= 0x01
	}
	if d.DCS != nil {
		b |= 0x02
	}
	if len(d.UD.Text) != 0 || len(d.UD.UDH) != 0 {
		b |= 0x04
	}
	w.WriteByte(b)
	w.Write(encodeSCTimeStamp(d.SCTS))
	if d.PID != nil {
		w.WriteByte(*d.PID)
	}
	if d.DCS != nil {
		w.WriteByte(d.DCS.Encode())
	}
	if len(d.UD.Text) != 0 || len(d.UD.UDH) != 0 {
		d.UD.write(w, d.DCS)
	}
	return w.Bytes()
}

// Decode get data of this TPDU
func (d *SubmitReport) Decode(b []byte) (e error) {
	if d == nil {
		return fmt.Errorf("nil data")
	}

	r := bytes.NewReader(b[1:])

	var pi byte
	if pi, e = r.ReadByte(); e == nil && pi&0x80 == 0x80 {
		d.FCS = pi
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
		d.UD.read(r, d.DCS, b[0]&0x40 == 0x40)
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

func (d *SubmitReport) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Submit Report")
	if d.FCS&0x80 == 0x80 {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, "%sTP-FCS:  %s\n", Indent, fcsStat(d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
	}

	fmt.Fprintf(w, "%sTP-SCTS: %s\n", Indent, d.SCTS)
	if d.PID != nil {
		fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(*d.PID))
	}
	if d.DCS != nil {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	}
	fmt.Fprintf(w, "%s", d.UD.String())
	return w.String()
}
