package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// StatusReport is TPDU message from SC to MS
type StatusReport struct {
	MMS bool `json:"mms"` // M / More Messages to Send (true=more messages)
	LP  bool `json:"lp"`  // O / Loop Prevention
	SRQ bool `json:"srq"` // M / Status Report Qualifier (true=status report shall be returned)

	MR   byte      `json:"mr"`            // M / Message Reference
	RA   Address   `json:"ra"`            // M / Destination Address
	SCTS time.Time `json:"scts"`          // M / Service Centre Time Stamp
	DT   time.Time `json:"dt"`            // M / Discharge Time
	ST   byte      `json:"st"`            // M / Status
	PID  *byte     `json:"pid,omitempty"` // O / Protocol Identifier
	DCS  DCS       `json:"dcs,omitempty"` // O / Data Coding Scheme
	UD   UD        `json:"ud,omitempty"`  // O / User Data
}

// MarshalTP output byte data of this TPDU
func (d StatusReport) MarshalTP() []byte {
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
	if d.UD.UDH != nil && len(d.UD.UDH) != 0 {
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
	if len(d.UD.Text) != 0 || len(d.UD.UDH) != 0 {
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
		w.WriteByte(d.DCS.Encode())
	}
	if len(d.UD.Text) != 0 || len(d.UD.UDH) != 0 {
		d.UD.write(w, d.DCS)
	}
	return w.Bytes()
}

// UnmarshalStatusReport decode StatusReport from bytes
func UnmarshalStatusReport(b []byte) (d StatusReport, e error) {
	e = d.UnmarshalTP(b)
	return
}

// UnmarshalTP get data of this TPDU
func (d *StatusReport) UnmarshalTP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0]&0x03 != 0x02 {
		e = &InvalidDataError{
			Name: "invalid MTI"}
	}

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

// MarshalJSON provide custom marshaller
func (d StatusReport) MarshalJSON() ([]byte, error) {
	type alias StatusReport
	al := struct {
		*alias
		Pid *byte `json:"pid,omitempty"`
		Dcs *byte `json:"dcs,omitempty"`
		Ud  *UD   `json:"ud,omitempty"`
	}{
		alias: (*alias)(&d)}
	al.Pid = d.PID
	if d.DCS != nil {
		tmp := d.DCS.Encode()
		al.Dcs = &tmp
	}
	if !d.UD.isEmpty() {
		al.Ud = &d.UD
	}
	return json.Marshal(al)
}

// UnmarshalJSON provide custom marshaller
func (d *StatusReport) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type alias StatusReport
	al := struct {
		Pid *byte `json:"pid,omitempty"`
		Dcs *byte `json:"dcs,omitempty"`
		Ud  *UD   `json:"ud,omitempty"`
		*alias
	}{
		alias: (*alias)(d)}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	d.PID = al.Pid
	if al.Dcs != nil {
		d.DCS = DecodeDCS(*al.Dcs)
	}
	if al.Ud != nil {
		d.UD = *al.Ud
	}
	return nil
}

func (d StatusReport) String() string {
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
	if !d.UD.isEmpty() {
		fmt.Fprintf(w, "%s", d.UD.String())
	}
	return w.String()
}
