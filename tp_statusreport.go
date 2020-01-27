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
	rpRequest

	MMS bool `json:"mms"` // M / More Messages to Send (true=more messages)
	LP  bool `json:"lp"`  // O / Loop Prevention
	SRQ bool `json:"srq"` // M / Status Report Qualifier (true=status report shall be returned)

	TMR  byte       `json:"tmr"`           // M / Message Reference
	RA   Address    `json:"ra"`            // M / Destination Address
	SCTS time.Time  `json:"scts"`          // M / Service Centre Time Stamp
	DT   time.Time  `json:"dt"`            // M / Discharge Time
	ST   byte       `json:"st"`            // M / Status
	PID  *byte      `json:"pid,omitempty"` // O / Protocol Identifier
	DCS  DataCoding `json:"dcs,omitempty"` // O / Data Coding Scheme
	UD   UserData   `json:"ud,omitempty"`  // O / User Data
}

// MarshalTP output byte data of this TPDU
func (d StatusReport) MarshalTP() []byte {
	w := new(bytes.Buffer)

	b := byte(0x02)
	if !d.MMS {
		b |= 0x04
	}
	if d.LP {
		b |= 0x08
	}
	if d.SRQ {
		b |= 0x20
	}
	if len(d.UD.UDH) != 0 {
		b |= 0x40
	}
	w.WriteByte(b)
	w.WriteByte(d.TMR)
	l, a := d.RA.Marshal()
	w.WriteByte(l)
	w.Write(a)
	w.Write(marshalSCTimeStamp(d.SCTS))
	w.Write(marshalSCTimeStamp(d.DT))
	w.WriteByte(d.ST)
	b = byte(0x00)
	if d.PID != nil {
		b |= 0x01
	}
	if d.DCS != nil {
		b |= 0x02
	}
	if !d.UD.isEmpty() {
		b |= 0x04
	}
	if b == 0x00 {
		return w.Bytes()
	}
	w.WriteByte(b)
	if d.PID != nil {
		w.WriteByte(*d.PID)
	}
	if d.DCS != nil {
		w.WriteByte(d.DCS.Marshal())
	}
	if !d.UD.isEmpty() {
		d.UD.write(w, d.DCS)
	}
	return w.Bytes()
}

// MarshalRP output byte data of this RPDU
func (d StatusReport) MarshalRP() []byte {
	return d.rpRequest.marshal(false, d.MarshalTP())
}

// MarshalCP output byte data of this CPDU
func (d StatusReport) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
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
		return UnexpectedMessageTypeError{
			Expected: 0x02, Actual: b[0] & 0x03}
	}

	d.MMS = b[0]&0x04 != 0x04
	d.LP = b[0]&0x08 == 0x08
	d.SRQ = b[0]&0x20 == 0x20

	r := bytes.NewReader(b[1:])

	if d.TMR, e = r.ReadByte(); e != nil {
		return
	}
	if d.RA, e = readTPAddr(r); e != nil {
		return
	}
	var p [7]byte
	if p, e = read7Bytes(r); e != nil {
		return
	}
	d.SCTS = unmarshalSCTimeStamp(p)
	if p, e = read7Bytes(r); e != nil {
		return
	}
	d.DT = unmarshalSCTimeStamp(p)
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
	if pi&0x01 == 0x01 {
		var p byte
		if p, e = r.ReadByte(); e != nil {
			return
		}
		d.PID = &p
	}
	if pi&0x02 == 0x02 {
		if d.DCS, e = readDataCoding(r); e != nil {
			return
		}
	}
	if pi&0x04 == 0x04 {
		if e = d.UD.read(r, d.DCS, b[0]&0x40 == 0x40); e != nil {
			return
		}
	}
	if r.Len() != 0 {
		e = ErrExtraData
	}
	return
}

// UnmarshalRP get data of this RPDU
func (d *StatusReport) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshal(false, b); e == nil {
		e = d.UnmarshalTP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *StatusReport) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// MarshalJSON provide custom marshaller
func (d StatusReport) MarshalJSON() ([]byte, error) {
	type alias StatusReport
	al := struct {
		*alias
		Pid *byte     `json:"pid,omitempty"`
		Dcs *byte     `json:"dcs,omitempty"`
		Ud  *UserData `json:"ud,omitempty"`
	}{
		alias: (*alias)(&d)}
	al.Pid = d.PID
	if d.DCS != nil {
		tmp := d.DCS.Marshal()
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
		Pid *byte     `json:"pid,omitempty"`
		Dcs *byte     `json:"dcs,omitempty"`
		Ud  *UserData `json:"ud,omitempty"`
		*alias
	}{
		alias: (*alias)(d)}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	d.PID = al.Pid
	if al.Dcs != nil {
		d.DCS = UnmarshalDataCoding(*al.Dcs)
	}
	if al.Ud != nil {
		d.UD = *al.Ud
	}
	return nil
}

func (d StatusReport) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "TP-StatusReport\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
	fmt.Fprintf(w, "%sRP-OA:   %s\n", Indent, d.SCA)
	fmt.Fprintf(w, "%sRP-DA:   <nil>\n", Indent)
	fmt.Fprintf(w, "%sTP-MMS:  %s\n", Indent, mmsStat(d.MMS))
	fmt.Fprintf(w, "%sTP-LP:   %s\n", Indent, lpStat(d.LP))
	fmt.Fprintf(w, "%sTP-SRQ:  %s\n", Indent, srqStat(d.SRQ))
	fmt.Fprintf(w, "%sTP-MR:   %d\n", Indent, d.TMR)
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
		fmt.Fprintf(w, "%sTP-UD:\n", Indent)
		fmt.Fprintf(w, "%s", d.UD.String())
	}

	return w.String()
}
