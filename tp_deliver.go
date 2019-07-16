package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Deliver is TPDU message from SC to MS
type Deliver struct {
	MMS bool `json:"mms"` // M / More Messages to Send (true=more messages)
	LP  bool `json:"lp"`  // O / Loop Prevention
	SRI bool `json:"sri"` // O / Status Report Indication (true=status report shall be returned)
	RP  bool `json:"rp"`  // M / Reply Path

	OA   Address   `json:"oa"`   // M / Originating Address
	PID  byte      `json:"pid"`  // M / Protocol Identifier
	DCS  DCS       `json:"dcs"`  // M / Data Coding Scheme
	SCTS time.Time `json:"scts"` // M / Service Centre Time Stamp
	UD   UD        `json:"ud"`   // O / User Data
}

// Encode output byte data of this TPDU
func (d *Deliver) Encode() []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)

	b := byte(0x00)
	if !d.MMS {
		b |= 0x04
	}
	if d.LP {
		b |= 0x08
	}
	if d.SRI {
		b |= 0x20
	}
	if len(d.UD.UDH) != 0 {
		b |= 0x40
	}
	if d.RP {
		b |= 0x80
	}
	w.WriteByte(b)
	b, a := d.OA.Encode()
	w.WriteByte(b)
	w.Write(a)
	w.WriteByte(d.PID)
	if d.DCS == nil {
		w.WriteByte(0x00)
	} else {
		w.WriteByte(d.DCS.Encode())
	}
	w.Write(encodeSCTimeStamp(d.SCTS))
	d.UD.write(w, d.DCS)

	return w.Bytes()
}

// Decode get data of this TPDU
func (d *Deliver) Decode(b []byte) (e error) {
	if d == nil {
		return fmt.Errorf("nil data")
	}

	d.MMS = b[0]&0x04 != 0x04
	d.LP = b[0]&0x08 == 0x08
	d.SRI = b[0]&0x20 == 0x20
	d.RP = b[0]&0x80 == 0x80

	r := bytes.NewReader(b[1:])
	if d.OA, e = readAddr(r); e != nil {
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

// UnmarshalJSON provide custom marshaller
func (d *Deliver) UnmarshalJSON(b []byte) error {
	type alias Deliver
	al := struct {
		Dcs byte `json:"dcs"`
		*alias
	}{
		alias: (*alias)(d),
	}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	d.DCS = DecodeDCS(al.Dcs)
	return nil
}

// MarshalJSON provide custom marshaller
func (d *Deliver) MarshalJSON() ([]byte, error) {
	type alias Deliver
	al := struct {
		*alias
		Dcs byte `json:"dcs"`
	}{
		Dcs:   d.DCS.Encode(),
		alias: (*alias)(d),
	}
	return json.Marshal(al)
}

func (d *Deliver) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Deliver\n")
	fmt.Fprintf(w, "%sTP-MMS:  %s\n", Indent, mmsStat(d.MMS))
	fmt.Fprintf(w, "%sTP-LP:   %s\n", Indent, lpStat(d.LP))
	fmt.Fprintf(w, "%sTP-SRI:  %s\n", Indent, sriStat(d.SRI))
	fmt.Fprintf(w, "%sTP-RP:   %s\n", Indent, rpStat(d.RP))
	fmt.Fprintf(w, "%sTP-OA:   %s\n", Indent, d.OA)
	fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(d.PID))
	fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	fmt.Fprintf(w, "%sTP-SCTS: %s\n", Indent, d.SCTS)
	fmt.Fprintf(w, "%s", d.UD.String())
	return w.String()
}

// DeliverReport is TPDU message from MS to SC
type DeliverReport struct {
	FCS byte  `json:"fcs"` // C / Failure Cause
	PID *byte `json:"pid"` // O / Protocol Identifier
	DCS DCS   `json:"dcs"` // O / Data Coding Scheme
	UD  UD    `json:"ud"`  // O / User Data
}

// Encode output byte data of this TPDU
func (d *DeliverReport) Encode() []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)

	b := byte(0x00)
	if len(d.UD.UDH) != 0 {
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
func (d *DeliverReport) Decode(b []byte) (e error) {
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

// UnmarshalJSON provide custom marshaller
func (d *DeliverReport) UnmarshalJSON(b []byte) error {
	type alias DeliverReport
	al := struct {
		Dcs byte `json:"dcs"`
		*alias
	}{
		alias: (*alias)(d),
	}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	d.DCS = DecodeDCS(al.Dcs)
	return nil
}

// MarshalJSON provide custom marshaller
func (d *DeliverReport) MarshalJSON() ([]byte, error) {
	type alias DeliverReport
	al := struct {
		*alias
		Dcs byte `json:"dcs"`
	}{
		Dcs:   d.DCS.Encode(),
		alias: (*alias)(d),
	}
	return json.Marshal(al)
}

func (d *DeliverReport) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Deliver Report")
	if d.FCS&0x80 == 0x80 {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, "%sTP-FCS:  %s\n", Indent, fcsStat(d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
	}

	if d.PID != nil {
		fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(*d.PID))
	}
	if d.DCS != nil {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	}
	fmt.Fprintf(w, "%s", d.UD.String())
	return w.String()
}
